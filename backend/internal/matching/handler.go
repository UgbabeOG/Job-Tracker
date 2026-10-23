package matching

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"jobtrack/internal/httpapi"
)

const maxBodySize = 1 << 20 // 1 MiB

type MatchHandler struct{}

type matchRequestBody struct {
	CVText           string   `json:"cv_text"`
	JobDescription   string   `json:"job_description"`
	RequiredKeywords []string `json:"required_keywords"`
}

type matchResponse struct {
	Score           int      `json:"score"`
	MatchedKeywords []string `json:"matched_keywords"`
	MissingKeywords []string `json:"missing_keywords"`
	TotalKeywords   int      `json:"total_keywords"`
}

func NewHandler() *MatchHandler {
	return &MatchHandler{}
}

func (h *MatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		httpapi.WriteMethodNotAllowed(w)
		return
	}

	if !requireJSON(w, r) {
		return
	}

	var input matchRequestBody

	if err := decodeJSON(w, r, &input); err != nil {
		if errors.Is(err, ErrPayloadTooLarge) {
			httpapi.WritePayloadTooLarge(w)
			return
		}

		writeValidationError(w, "Invalid JSON request body", nil)
		return
	}

	result, err := Match(MatchRequest{
		CVText:           input.CVText,
		JobDescription:   input.JobDescription,
		RequiredKeywords: input.RequiredKeywords,
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyCV):
			writeValidationError(w, "The request contains invalid fields.", map[string]string{
				"cv_text": "CV text is required.",
			})

		case errors.Is(err, ErrEmptyJob):
			writeValidationError(w, "The request contains invalid fields.", map[string]string{
				"job_description": "Job description is required.",
			})

		case errors.Is(err, ErrNoKeywordsFound):
			writeMatchError(w, "NO_KEYWORDS_FOUND", "No valid keywords were found to score.", nil)

		default:
			httpapi.WriteInternalError(w)
		}

		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]any{
		"data": matchResponse{
			Score:           result.Score,
			MatchedKeywords: result.MatchedKeywords,
			MissingKeywords: result.MissingKeywords,
			TotalKeywords:   result.TotalKeywords,
		},
	})
}

func requireJSON(w http.ResponseWriter, r *http.Request) bool {
	contentType := strings.TrimSpace(
		strings.Split(r.Header.Get("Content-Type"), ";")[0],
	)

	if contentType != "application/json" {
		httpapi.WriteUnsupportedMediaType(w)
		return false
	}

	return true
}

var ErrPayloadTooLarge = errors.New("request body exceeds the size limit")

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		var maxBytesError *http.MaxBytesError

		if errors.As(err, &maxBytesError) {
			return ErrPayloadTooLarge
		}

		return err
	}

	var extra any

	if err := decoder.Decode(&extra); err != io.EOF {
		var maxBytesError *http.MaxBytesError

		if errors.As(err, &maxBytesError) {
			return ErrPayloadTooLarge
		}

		return errors.New("request body must contain a single JSON value")
	}

	return nil
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func writeValidationError(w http.ResponseWriter, message string, details any) {
	writeMatchError(w, "VALIDATION_ERROR", message, details)
}

func writeMatchError(w http.ResponseWriter, code, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: errorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
