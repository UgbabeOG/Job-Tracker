package assistant

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"jobtrack/internal/httpapi"
	"jobtrack/internal/jobs"
	"jobtrack/internal/matching"
	"jobtrack/internal/profile"
)

const maxBodySize = 1 << 20

var ErrPayloadTooLarge = errors.New("request body exceeds the size limit")

type Handler struct {
	Repo    *Repository
	JobRepo *jobs.Repository
	Prof    *profile.Repository
}

func NewHandler(repo *Repository, jobRepo *jobs.Repository, prof *profile.Repository) *Handler {
	return &Handler{Repo: repo, JobRepo: jobRepo, Prof: prof}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/assistant")
	path = strings.Trim(path, "/")
	if path == "" {
		httpapi.WriteNotFound(w, "Not found")
		return
	}
	parts := strings.Split(path, "/")
	if parts[0] == "parse-url" {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			httpapi.WriteMethodNotAllowed(w)
			return
		}
		h.parseURL(w, r)
		return
	}
	if parts[0] == "suggestions" && len(parts) == 2 {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			httpapi.WriteMethodNotAllowed(w)
			return
		}
		h.suggestions(w, r, parts[1])
		return
	}
	if parts[0] == "audit" {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			httpapi.WriteMethodNotAllowed(w)
			return
		}
		h.audit(w, r)
		return
	}
	httpapi.WriteNotFound(w, "Not found")
}

type parseReq struct {
	URL           string `json:"url"`
	SaveToTracker bool   `json:"save_to_tracker"`
	Company       string `json:"company"`
	JobTitle      string `json:"job_title"`
}

func (h *Handler) parseURL(w http.ResponseWriter, r *http.Request) {
	if !requireJSON(w, r) {
		return
	}
	var input parseReq
	if err := decodeJSON(w, r, &input); err != nil {
		if errors.Is(err, ErrPayloadTooLarge) {
			httpapi.WritePayloadTooLarge(w)
			return
		}
		httpapi.WriteBadRequest(w, "Invalid JSON request body")
		return
	}
	if strings.TrimSpace(input.URL) == "" {
		httpapi.WriteBadRequest(w, "url is required")
		return
	}
	parsed, err := ParseFromURL(input.URL)
	if err != nil {
		httpapi.WriteBadRequest(w, "Unable to parse URL")
		return
	}
	if input.Company != "" {
		parsed.Company = input.Company
	}
	if input.JobTitle != "" {
		parsed.Title = input.JobTitle
	}
	_ = h.Repo.Audit(r.Context(), "parse_url", map[string]any{"url": input.URL, "platform": parsed.SourcePlatform})

	var jobID string
	if input.SaveToTracker {
		job, err := h.JobRepo.Create(r.Context(), jobs.CreateJobRequest{
			Company:        parsed.Company,
			JobTitle:       parsed.Title,
			JobURL:         &parsed.SourceURL,
			JobDescription: &parsed.Description,
			Status:         "Saved",
		})
		if err == nil {
			jobID = job.ID
			_ = h.Repo.SaveSource(r.Context(), job.ID, parsed.SourceURL, parsed.SourcePlatform, parsed)
		}
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"parsed": parsed,
			"job_id": jobID,
		},
	})
}

func (h *Handler) suggestions(w http.ResponseWriter, r *http.Request, jobID string) {
	job, err := h.JobRepo.Get(r.Context(), jobID)
	if errors.Is(err, jobs.ErrNotFound) {
		httpapi.WriteNotFound(w, "Job not found")
		return
	}
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}
	prof, _ := h.Prof.Get(r.Context())
	cv := ""
	if prof.CVText != nil {
		cv = *prof.CVText
	}
	desc := ""
	if job.JobDescription != nil {
		desc = *job.JobDescription
	}
	res, err := matching.Match(matching.MatchRequest{CVText: cv, JobDescription: desc})
	if err != nil {
		// build minimal suggestions
		httpapi.WriteJSON(w, http.StatusOK, map[string]any{"data": SuggestionsResponse{MatchScore: 0, Suggestions: []Suggestion{{Field: "cv", Suggestion: "Add your CV in Profile to get suggestions.", NeedsReview: true}}, Missing: []string{}, Matched: []string{}, Total: 0}})
		return
	}
	sugs := SuggestionsResponse{
		MatchScore: res.Score,
		Missing:    res.MissingKeywords,
		Matched:    res.MatchedKeywords,
		Total:      res.TotalKeywords,
	}
	// build actionable suggestions
	if prof.FullName != nil {
		sugs.Suggestions = append(sugs.Suggestions, Suggestion{Field: "full_name", Suggestion: "Use your full name as in your CV/resume.", Reason: "Keep consistent across applications.", NeedsReview: false})
	}
	if prof.Email != nil {
		sugs.Suggestions = append(sugs.Suggestions, Suggestion{Field: "email", Suggestion: "Use your professional email.", Reason: "Avoid personal/unprofessional emails.", NeedsReview: false})
	}
	if prof.LinkedInURL != nil {
		sugs.Suggestions = append(sugs.Suggestions, Suggestion{Field: "linkedin", Suggestion: "Have your LinkedIn URL ready to paste.", Reason: "Many forms ask for it.", NeedsReview: false})
	}
	if len(res.MissingKeywords) > 0 {
		sugs.Suggestions = append(sugs.Suggestions, Suggestion{Field: "cover_letter", Suggestion: "Mention relevant experience for missing skills: " + strings.Join(res.MissingKeywords, ", ") + ".", Reason: "Address gaps where possible.", NeedsReview: true})
	}
	_ = h.Repo.SaveSuggestions(r.Context(), jobID, sugs)
	_ = h.Repo.Audit(r.Context(), "get_suggestions", map[string]any{"job_id": jobID, "score": res.Score})
	httpapi.WriteJSON(w, http.StatusOK, map[string]any{"data": sugs})
}

func (h *Handler) audit(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.ListAudit(r.Context())
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, map[string]any{"data": items})
}

func requireJSON(w http.ResponseWriter, r *http.Request) bool {
	ct := strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0])
	if ct != "application/json" {
		httpapi.WriteUnsupportedMediaType(w)
		return false
	}
	return true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			return ErrPayloadTooLarge
		}
		return err
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			return ErrPayloadTooLarge
		}
		return errors.New("request body must contain a single JSON value")
	}
	return nil
}
