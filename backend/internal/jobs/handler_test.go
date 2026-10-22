package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"jobtrack/internal/database"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	t.Setenv("JOBTRACK_DB", dbPath)

	db, err := database.Open()
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close test database: %v", err)
		}
	})

	repo := NewRepository(db)
	return NewHandler(repo)
}
func request(
	t *testing.T,
	h http.Handler,
	method, path, body string,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func createTestJob(t *testing.T, h http.Handler) string {
	t.Helper()

	rec := request(t, h, http.MethodPost, "/api/jobs",
		`{"company":"Test Company","job_title":"Go Engineer","status":"Saved"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create job status = %d, want %d; body: %s",
			rec.Code, http.StatusCreated, rec.Body.String())
	}

	var response struct {
		Data Job `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode created job: %v", err)
	}
	if response.Data.ID == "" {
		t.Fatal("created job has an empty ID")
	}

	return response.Data.ID
}

func TestCreateAndListJobs(t *testing.T) {
	h := newTestHandler(t)

	createTestJob(t, h)

	rec := request(t, h, http.MethodGet, "/api/jobs", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d; body: %s",
			rec.Code, http.StatusOK, rec.Body.String())
	}

	var response struct {
		Data []Job `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode job list: %v", err)
	}
	if len(response.Data) != 1 {
		t.Fatalf("job count = %d, want 1", len(response.Data))
	}
	if response.Data[0].Company != "Test Company" {
		t.Errorf("company = %q, want Test Company", response.Data[0].Company)
	}
	if response.Data[0].Status != "Saved" {
		t.Errorf("status = %q, want Saved", response.Data[0].Status)
	}
}

func TestPatchPreservesUnchangedFieldsAndClearsOptionalField(t *testing.T) {
	h := newTestHandler(t)
	id := createTestJob(t, h)

	rec := request(t, h, http.MethodPatch, "/api/jobs/"+id,
		`{"status":"Applied","notes":"Follow up"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want %d; body: %s",
			rec.Code, http.StatusOK, rec.Body.String())
	}

	rec = request(t, h, http.MethodPatch, "/api/jobs/"+id, `{"notes":null}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("clear notes status = %d, want %d; body: %s",
			rec.Code, http.StatusOK, rec.Body.String())
	}

	rec = request(t, h, http.MethodGet, "/api/jobs/"+id, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response struct {
		Data Job `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode job: %v", err)
	}

	job := response.Data
	if job.Company != "Test Company" {
		t.Errorf("company changed unexpectedly: %q", job.Company)
	}
	if job.JobTitle != "Go Engineer" {
		t.Errorf("job title changed unexpectedly: %q", job.JobTitle)
	}
	if job.Status != "Applied" {
		t.Errorf("status = %q, want Applied", job.Status)
	}
	if job.Notes != nil {
		t.Errorf("notes = %v, want nil", *job.Notes)
	}
}

func TestInvalidRequests(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		name        string
		method      string
		path        string
		body        string
		contentType string
		wantStatus  int
	}{
		{
			name:       "malformed JSON",
			method:     http.MethodPost,
			path:       "/api/jobs",
			body:       `{"company":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing required company",
			method:     http.MethodPost,
			path:       "/api/jobs",
			body:       `{"job_title":"Go Engineer"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty PATCH body",
			method:     http.MethodPatch,
			path:       "/api/jobs/example-id",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "read-only field",
			method:     http.MethodPatch,
			path:       "/api/jobs/example-id",
			body:       `{"id":"changed"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "unsupported content type",
			method:      http.MethodPost,
			path:        "/api/jobs",
			body:        `{"company":"Test","job_title":"Engineer"}`,
			contentType: "text/plain",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			} else if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s",
					rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestMissingJobAndDelete(t *testing.T) {
	h := newTestHandler(t)

	rec := request(t, h, http.MethodGet, "/api/jobs/not-a-real-id", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("missing GET status = %d, want %d",
			rec.Code, http.StatusNotFound)
	}

	id := createTestJob(t, h)

	rec = request(t, h, http.MethodDelete, "/api/jobs/"+id, "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d",
			rec.Code, http.StatusNoContent)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("DELETE response should have an empty body, got %q", rec.Body.String())
	}

	rec = request(t, h, http.MethodGet, "/api/jobs/"+id, "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET after delete status = %d, want %d",
			rec.Code, http.StatusNotFound)
	}
}

func TestUnsupportedMethod(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPut, "/api/jobs", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if !strings.Contains(rec.Header().Get("Allow"), "GET") {
		t.Errorf("Allow header = %q, expected supported methods", rec.Header().Get("Allow"))
	}
}

func TestRepositoryContextIsAvailable(t *testing.T) {
	h := newTestHandler(t)
	if h.Repo == nil {
		t.Fatal("handler repository is nil")
	}
	if err := h.Repo.DB.PingContext(context.Background()); err != nil {
		t.Fatalf("test database is not available: %v", err)
	}
}

func TestDashboardEndpoint(t *testing.T) {
	h := newTestHandler(t)

	createTestJob(t, h)
	request(t, h, http.MethodPost, "/api/jobs",
		`{"company":"B","job_title":"Designer","status":"Applied","deadline":"2026-10-27"}`)

	rec := request(t, h, http.MethodGet, "/api/dashboard", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d, want %d; body: %s",
			rec.Code, http.StatusOK, rec.Body.String())
	}

	var response struct {
		Data struct {
			TotalApplications   int            `json:"total_applications"`
			StatusCounts        map[string]int `json:"status_counts"`
			UpcomingDeadlines   int            `json:"upcoming_deadlines"`
			OverdueApplications int            `json:"overdue_applications"`
			TotalInterviews     int            `json:"total_interviews"`
			TotalOffers         int            `json:"total_offers"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode dashboard: %v", err)
	}

	if response.Data.TotalApplications != 2 {
		t.Errorf("total_applications = %d, want 2", response.Data.TotalApplications)
	}
	if response.Data.StatusCounts["Saved"] != 1 {
		t.Errorf("Saved count = %d, want 1", response.Data.StatusCounts["Saved"])
	}
	if response.Data.StatusCounts["Applied"] != 1 {
		t.Errorf("Applied count = %d, want 1", response.Data.StatusCounts["Applied"])
	}
	if response.Data.UpcomingDeadlines != 1 {
		t.Errorf("upcoming_deadlines = %d, want 1", response.Data.UpcomingDeadlines)
	}
	if response.Data.TotalInterviews != 0 {
		t.Errorf("total_interviews = %d, want 0", response.Data.TotalInterviews)
	}
}

func TestDashboardDashboardMethodNotAllowed(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/api/dashboard", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("dashboard POST status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestValidateJobInputRejectsInvalidDateAndURL(t *testing.T) {
	tests := []struct {
		name        string
		url         *string
		dateApplied *string
		deadline    *string
	}{
		{
			name:        "impossible date",
			dateApplied: stringPointer("2026-02-30"),
		},
		{
			name: "invalid URL scheme",
			url:  stringPointer("ftp://example.com/job"),
		},
		{
			name: "URL missing host",
			url:  stringPointer("https:///job"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJobInput(
				"Example Company",
				"Software Engineer",
				"Saved",
				tt.url,
				tt.dateApplied,
				tt.deadline,
			)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestOversizedRequestReturns413(t *testing.T) {
	h := newTestHandler(t)

	body := `{"company":"` + strings.Repeat("a", maxBodySize) +
		`","job_title":"Engineer"}`

	rec := request(t, h, http.MethodPost, "/api/jobs", body)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d; body: %s",
			rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
}

func stringPointer(value string) *string {
	return &value
}
