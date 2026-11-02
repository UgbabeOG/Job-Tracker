package assistant

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) SaveSource(ctx context.Context, jobID, sourceURL, platform string, parsed any) error {
	id := uuid.NewString()
	data, _ := json.Marshal(parsed)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO job_sources (id, job_id, source_url, source_platform, parsed_json, created_at) VALUES (?,?,?,?,?,?)`,
		id, jobID, sourceURL, platform, string(data), Now())
	return err
}

func (r *Repository) SaveSuggestions(ctx context.Context, jobID string, suggestions SuggestionsResponse) error {
	id := uuid.NewString()
	data, _ := json.Marshal(suggestions)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO apply_suggestions (id, job_id, suggestions_json, match_score, created_at) VALUES (?,?,?,?,?)`,
		id, jobID, string(data), suggestions.MatchScore, Now())
	return err
}

func (r *Repository) Audit(ctx context.Context, action string, details any) error {
	id := uuid.NewString()
	data, _ := json.Marshal(details)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO assistant_audit (id, action, details, created_at) VALUES (?,?,?,?)`,
		id, action, string(data), Now())
	return err
}

func (r *Repository) ListAudit(ctx context.Context) ([]map[string]any, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, action, details, created_at FROM assistant_audit ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []map[string]any
	for rows.Next() {
		var id, action, created string
		var details sql.NullString
		if err := rows.Scan(&id, &action, &details, &created); err != nil {
			return nil, err
		}
		m := map[string]any{"id": id, "action": action, "created_at": created}
		if details.Valid {
			var d any
			_ = json.Unmarshal([]byte(details.String), &d)
			m["details"] = d
		}
		res = append(res, m)
	}
	return res, rows.Err()
}
