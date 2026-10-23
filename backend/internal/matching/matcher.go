package matching

import (
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
)

var (
	ErrEmptyCV         = errors.New("cv text is required")
	ErrEmptyJob        = errors.New("job description is required")
	ErrNoKeywordsFound = errors.New("no keywords found")
	whitespacePattern  = regexp.MustCompile(`\s+`)
)

// SkillVocabulary is the fixed vocabulary used to extract candidate
// keywords from a job description when RequiredKeywords is omitted.
//
// Keep this list deterministic and documented so matching results do not
// depend on undocumented keyword-guessing behaviour.
var SkillVocabulary = []string{
	"AWS",
	"Azure",
	"C",
	"C++",
	"C#",
	"CSS",
	"Django",
	"Docker",
	"Express",
	"Flask",
	"Git",
	"GitHub",
	"Go",
	"GraphQL",
	"HTML",
	"Java",
	"JavaScript",
	"Kubernetes",
	"MongoDB",
	"MySQL",
	"Next.js",
	"Node.js",
	"PostgreSQL",
	"Python",
	"React",
	"REST API",
	"SQLite",
	"SQL",
	"TypeScript",
	"Vue",
}

type MatchRequest struct {
	CVText           string
	JobDescription   string
	RequiredKeywords []string
}

type MatchResult struct {
	Score           int
	MatchedKeywords []string
	MissingKeywords []string
	TotalKeywords   int
}

// Match compares CV text against required job keywords.
//
// When RequiredKeywords is empty, candidate keywords are extracted from
// JobDescription using SkillVocabulary.
func Match(request MatchRequest) (MatchResult, error) {
	if strings.TrimSpace(request.CVText) == "" {
		return MatchResult{}, ErrEmptyCV
	}

	if strings.TrimSpace(request.JobDescription) == "" {
		return MatchResult{}, ErrEmptyJob
	}

	keywords := deduplicateKeywords(request.RequiredKeywords)

	if len(keywords) == 0 {
		keywords = extractKeywords(request.JobDescription)
	}

	if len(keywords) == 0 {
		return MatchResult{}, ErrNoKeywordsFound
	}

	matched := make([]string, 0, len(keywords))
	missing := make([]string, 0, len(keywords))

	for _, keyword := range keywords {
		if containsKeyword(request.CVText, keyword) {
			matched = append(matched, keyword)
		} else {
			missing = append(missing, keyword)
		}
	}

	score := int(math.Round(
		float64(len(matched)) / float64(len(keywords)) * 100,
	))

	return MatchResult{
		Score:           score,
		MatchedKeywords: matched,
		MissingKeywords: missing,
		TotalKeywords:   len(keywords),
	}, nil
}

func extractKeywords(jobDescription string) []string {
	found := make([]string, 0)

	for _, keyword := range SkillVocabulary {
		if containsKeyword(jobDescription, keyword) {
			found = append(found, keyword)
		}
	}

	return found
}

func deduplicateKeywords(keywords []string) []string {
	seen := make(map[string]string)

	for _, keyword := range keywords {
		cleaned := normalizeKeyword(keyword)
		if cleaned == "" {
			continue
		}

		key := strings.ToLower(cleaned)

		if _, exists := seen[key]; !exists {
			seen[key] = cleaned
		}
	}

	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, seen[key])
	}

	return result
}

func normalizeKeyword(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	return whitespacePattern.ReplaceAllString(keyword, " ")
}

func containsKeyword(text, keyword string) bool {
	text = strings.ToLower(text)
	keyword = strings.ToLower(normalizeKeyword(keyword))

	if keyword == "" {
		return false
	}

	pattern := `(^|[^\pL\pN])` +
		regexp.QuoteMeta(keyword) +
		`($|[^\pL\pN])`

	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false
	}

	return matched
}
