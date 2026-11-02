package profile

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Get(ctx context.Context) (Profile, error) {
	var p Profile
	var fn, em, ph, li, po, cv, cpath sql.NullString
	err := r.DB.QueryRowContext(ctx, `SELECT id, full_name, email, phone, linkedin_url, portfolio_url, cv_text, cv_path, created_at, updated_at FROM profiles LIMIT 1`).Scan(
		&p.ID, &fn, &em, &ph, &li, &po, &cv, &cpath, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Profile{}, ErrNotFound
	}
	if err != nil {
		return Profile{}, err
	}
	p.FullName = nullString(fn)
	p.Email = nullString(em)
	p.Phone = nullString(ph)
	p.LinkedInURL = nullString(li)
	p.PortfolioURL = nullString(po)
	p.CVText = nullString(cv)
	p.CVPath = nullString(cpath)
	return p, nil
}

func (r *Repository) Upsert(ctx context.Context, in UpdateProfileRequest) (Profile, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	existing, err := r.Get(ctx)
	if err == ErrNotFound {
		id := uuid.NewString()
		_, err = r.DB.ExecContext(ctx, `INSERT INTO profiles (id, full_name, email, phone, linkedin_url, portfolio_url, cv_text, cv_path, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`,
			id, in.FullName, in.Email, in.Phone, in.LinkedInURL, in.PortfolioURL, in.CVText, in.CVPath, now, now)
		if err != nil {
			return Profile{}, err
		}
		return r.Get(ctx)
	}
	if err != nil {
		return Profile{}, err
	}
	// merge
	if in.FullName != nil {
		existing.FullName = in.FullName
	}
	if in.Email != nil {
		existing.Email = in.Email
	}
	if in.Phone != nil {
		existing.Phone = in.Phone
	}
	if in.LinkedInURL != nil {
		existing.LinkedInURL = in.LinkedInURL
	}
	if in.PortfolioURL != nil {
		existing.PortfolioURL = in.PortfolioURL
	}
	if in.CVText != nil {
		existing.CVText = in.CVText
	}
	if in.CVPath != nil {
		existing.CVPath = in.CVPath
	}
	_, err = r.DB.ExecContext(ctx, `UPDATE profiles SET full_name=?, email=?, phone=?, linkedin_url=?, portfolio_url=?, cv_text=?, cv_path=?, updated_at=? WHERE id=?`,
		existing.FullName, existing.Email, existing.Phone, existing.LinkedInURL, existing.PortfolioURL, existing.CVText, existing.CVPath, now, existing.ID)
	if err != nil {
		return Profile{}, err
	}
	return r.Get(ctx)
}

func nullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
