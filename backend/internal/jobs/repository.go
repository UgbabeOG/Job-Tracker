package jobs

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) List(ctx context.Context) ([]Job, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, company, job_title, job_url, date_applied,
		       deadline, status, notes, job_description, created_at, updated_at
		FROM jobs ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Job, 0)
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, job)
	}
	return result, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id string) (Job, error) {
	job, err := scanJob(r.DB.QueryRowContext(ctx, `
		SELECT id, company, job_title, job_url, date_applied,
		       deadline, status, notes, job_description, created_at, updated_at
		FROM jobs WHERE id = ?
	`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Job{}, ErrNotFound
	}
	return job, err
}

func (r *Repository) Create(ctx context.Context, input CreateJobRequest) (Job, error) {
	status := input.Status
	if status == "" {
		status = "Saved"
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.NewString()

	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO jobs (
			id, company, job_title, job_url, date_applied,
			deadline, status, notes, job_description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, input.Company, input.JobTitle, input.JobURL, input.DateApplied,
		input.Deadline, status, input.Notes, input.JobDescription, now, now)
	if err != nil {
		return Job{}, err
	}

	return r.Get(ctx, id)
}

func (r *Repository) Update(ctx context.Context, job Job) (Job, error) {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE jobs SET
			company = ?, job_title = ?, job_url = ?, date_applied = ?,
			deadline = ?, status = ?, notes = ?, job_description = ?, updated_at = ?
		WHERE id = ?
	`, job.Company, job.JobTitle, job.JobURL, job.DateApplied,
		job.Deadline, job.Status, job.Notes, job.JobDescription,
		time.Now().UTC().Format(time.RFC3339), job.ID)
	if err != nil {
		return Job{}, err
	}
	return r.Get(ctx, job.ID)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM jobs WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanJob(row rowScanner) (Job, error) {
	var job Job
	var jobURL, dateApplied, deadline, notes, jobDescription sql.NullString

	err := row.Scan(
		&job.ID, &job.Company, &job.JobTitle,
		&jobURL, &dateApplied, &deadline,
		&job.Status, &notes, &jobDescription, &job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, ErrNotFound
		}
		return Job{}, err
	}

	job.JobURL = nullStringPointer(jobURL)
	job.DateApplied = nullStringPointer(dateApplied)
	job.Deadline = nullStringPointer(deadline)
	job.Notes = nullStringPointer(notes)
	job.JobDescription = nullStringPointer(jobDescription)

	return job, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
