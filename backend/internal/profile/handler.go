package profile

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"jobtrack/internal/httpapi"
)

const maxBodySize = 1 << 20

var ErrPayloadTooLarge = errors.New("request body exceeds the size limit")

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/profile")
	path = strings.Trim(path, "/")
	if path != "" {
		httpapi.WriteNotFound(w, "Not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPut:
		h.put(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT")
		httpapi.WriteMethodNotAllowed(w)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repo.Get(r.Context())
	if errors.Is(err, ErrNotFound) {
		httpapi.WriteJSON(w, http.StatusOK, map[string]any{"data": nil})
		return
	}
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, map[string]any{"data": p})
}

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	if !requireJSON(w, r) {
		return
	}
	var input UpdateProfileRequest
	if err := decodeJSON(w, r, &input); err != nil {
		if errors.Is(err, ErrPayloadTooLarge) {
			httpapi.WritePayloadTooLarge(w)
			return
		}
		httpapi.WriteBadRequest(w, "Invalid JSON request body")
		return
	}
	p, err := h.Repo.Upsert(r.Context(), input)
	if err != nil {
		httpapi.WriteInternalError(w)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, map[string]any{"data": p})
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
