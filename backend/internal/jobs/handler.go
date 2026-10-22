package jobs

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"jobtrack/internal/httpapi"
)

const maxBodySize = 1 << 20 // 1 MiB

var ErrPayloadTooLarge = errors.New("request body exceeds the size limit")

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/api/dashboard" {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			httpapi.WriteMethodNotAllowed(w)
			return
		}
		h.dashboard(w, r)
		return
	}

	path = strings.Trim(strings.TrimPrefix(path, "/api/jobs"), "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			httpapi.WriteMethodNotAllowed(w)
		}
		return
	}

	if strings.Contains(path, "/") {
		httpapi.WriteNotFound(w, "Job not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, r, path)
	case http.MethodPatch:
		h.update(w, r, path)
	case http.MethodDelete:
		h.delete(w, r, path)
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		httpapi.WriteMethodNotAllowed(w)
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	data, err := h.Repo.List(r.Context())
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	if !requireJSON(w, r) {
		return
	}

	var input CreateJobRequest

	if err := decodeJSON(w, r, &input); err != nil {
		if errors.Is(err, ErrPayloadTooLarge) {
			httpapi.WritePayloadTooLarge(w)
			return
		}

		httpapi.WriteBadRequest(w, "Invalid JSON request body")
		return
	}

	if err := ValidateJobInput(
		input.Company,
		input.JobTitle,
		input.Status,
		input.JobURL,
		input.DateApplied,
		input.Deadline,
	); err != nil {
		httpapi.WriteBadRequest(w, err.Error())
		return
	}

	job, err := h.Repo.Create(r.Context(), input)
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, map[string]any{
		"data": job,
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request, id string) {
	job, err := h.Repo.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpapi.WriteNotFound(w, "Job not found")
		return
	}
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]any{
		"data": job,
	})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, id string) {
	if !requireJSON(w, r) {
		return
	}

	var fields map[string]json.RawMessage

	if err := decodeJSON(w, r, &fields); err != nil {
		if errors.Is(err, ErrPayloadTooLarge) {
			httpapi.WritePayloadTooLarge(w)
			return
		}

		httpapi.WriteBadRequest(w, "Invalid JSON request body")
		return
	}

	if len(fields) == 0 {
		httpapi.WriteBadRequest(w, "PATCH body must contain at least one field")
		return
	}

	allowed := map[string]bool{
		"company":      true,
		"job_title":    true,
		"job_url":      true,
		"date_applied": true,
		"deadline":     true,
		"status":       true,
		"notes":        true,
	}

	for key := range fields {
		if !allowed[key] {
			httpapi.WriteBadRequest(
				w,
				"Unknown or read-only field: "+key,
			)
			return
		}
	}

	job, err := h.Repo.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpapi.WriteNotFound(w, "Job not found")
		return
	}
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	for key, raw := range fields {
		switch key {
		case "company":
			if err := json.Unmarshal(raw, &job.Company); err != nil {
				httpapi.WriteBadRequest(w, "company must be a string")
				return
			}

		case "job_title":
			if err := json.Unmarshal(raw, &job.JobTitle); err != nil {
				httpapi.WriteBadRequest(w, "job_title must be a string")
				return
			}

		case "job_url":
			if err := json.Unmarshal(raw, &job.JobURL); err != nil {
				httpapi.WriteBadRequest(w, "job_url must be a string or null")
				return
			}

		case "date_applied":
			if err := json.Unmarshal(raw, &job.DateApplied); err != nil {
				httpapi.WriteBadRequest(w, "date_applied must be a string or null")
				return
			}

		case "deadline":
			if err := json.Unmarshal(raw, &job.Deadline); err != nil {
				httpapi.WriteBadRequest(w, "deadline must be a string or null")
				return
			}

		case "status":
			if err := json.Unmarshal(raw, &job.Status); err != nil {
				httpapi.WriteBadRequest(w, "status must be a string")
				return
			}

		case "notes":
			if err := json.Unmarshal(raw, &job.Notes); err != nil {
				httpapi.WriteBadRequest(w, "notes must be a string or null")
				return
			}
		}
	}

	if err := ValidateJobInput(
		job.Company,
		job.JobTitle,
		job.Status,
		job.JobURL,
		job.DateApplied,
		job.Deadline,
	); err != nil {
		httpapi.WriteBadRequest(w, err.Error())
		return
	}

	updated, err := h.Repo.Update(r.Context(), job)
	if errors.Is(err, ErrNotFound) {
		httpapi.WriteNotFound(w, "Job not found")
		return
	}
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]any{
		"data": updated,
	})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id string) {
	err := h.Repo.Delete(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpapi.WriteNotFound(w, "Job not found")
		return
	}
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.Repo.List(r.Context())
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}

	stats := ComputeDashboardStats(nowInLagos(), jobs)

	httpapi.WriteJSON(w, http.StatusOK, map[string]any{
		"data": stats,
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
