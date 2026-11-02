package profile

import "errors"

var ErrNotFound = errors.New("profile not found")

type Profile struct {
	ID           string  `json:"id"`
	FullName     *string `json:"full_name"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	LinkedInURL  *string `json:"linkedin_url"`
	PortfolioURL *string `json:"portfolio_url"`
	CVText       *string `json:"cv_text"`
	CVPath       *string `json:"cv_path"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FullName     *string `json:"full_name"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	LinkedInURL  *string `json:"linkedin_url"`
	PortfolioURL *string `json:"portfolio_url"`
	CVText       *string `json:"cv_text"`
	CVPath       *string `json:"cv_path"`
}
