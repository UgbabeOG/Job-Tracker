package jobs

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("job not found")
	ErrInvalid  = errors.New("invalid job data")
)

var validStatuses = map[string]bool{
	"Saved":     true,
	"Applied":   true,
	"Interview": true,
	"Offer":     true,
	"Rejected":  true,
}

type Job struct {
	ID          string  `json:"id"`
	Company     string  `json:"company"`
	JobTitle    string  `json:"job_title"`
	JobURL      *string `json:"job_url"`
	DateApplied *string `json:"date_applied"`
	Deadline    *string `json:"deadline"`
	Status      string  `json:"status"`
	Notes       *string `json:"notes"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateJobRequest struct {
	Company     string  `json:"company"`
	JobTitle    string  `json:"job_title"`
	JobURL      *string `json:"job_url"`
	DateApplied *string `json:"date_applied"`
	Deadline    *string `json:"deadline"`
	Status      string  `json:"status"`
	Notes       *string `json:"notes"`
}

func ValidateJobInput(company, title, status string, jobURL, dateApplied, deadline *string) error {
	if strings.TrimSpace(company) == "" {
		return errors.New("company is required")
	}
	if strings.TrimSpace(title) == "" {
		return errors.New("job_title is required")
	}
	if status != "" && !validStatuses[status] {
		return errors.New("status must be Saved, Applied, Interview, Offer, or Rejected")
	}
	if err := validateURL(jobURL); err != nil {
		return err
	}
	if err := validateDate("date_applied", dateApplied); err != nil {
		return err
	}
	if err := validateDate("deadline", deadline); err != nil {
		return err
	}
	return nil
}

func validateURL(value *string) error {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}

	parsed, err := url.ParseRequestURI(*value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("job_url must be a valid HTTP or HTTPS URL")
	}
	return nil
}

func validateDate(field string, value *string) error {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil || parsed.Format("2006-01-02") != *value {
		return errors.New(field + " must use YYYY-MM-DD format")
	}
	return nil
}
