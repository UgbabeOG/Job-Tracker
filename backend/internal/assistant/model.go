package assistant

import "time"

type ParsedJob struct {
	Title          string `json:"title"`
	Company        string `json:"company"`
	Location       string `json:"location"`
	Description    string `json:"description"`
	Requirements   string `json:"requirements"`
	SourceURL      string `json:"source_url"`
	SourcePlatform string `json:"source_platform"`
}

type Suggestion struct {
	Field       string `json:"field"`
	Suggestion  string `json:"suggestion"`
	Reason      string `json:"reason"`
	NeedsReview bool   `json:"needs_review"`
}

type SuggestionsResponse struct {
	MatchScore  int          `json:"match_score"`
	Suggestions []Suggestion `json:"suggestions"`
	Missing     []string     `json:"missing_keywords"`
	Matched     []string     `json:"matched_keywords"`
	Total       int          `json:"total_keywords"`
}

func Now() string { return time.Now().UTC().Format(time.RFC3339) }
