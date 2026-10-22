package matching

import (
	"errors"
	"reflect"
	"testing"
)

func TestMatchCaseInsensitive(t *testing.T) {
	request := MatchRequest{
		CVText:           "I use go, sql and git.",
		JobDescription:   "Backend developer",
		RequiredKeywords: []string{"Go", "SQL", "Git"},
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 100 {
		t.Fatalf("expected score 100, got %d", result.Score)
	}
}

func TestMatchDeduplicatesKeywords(t *testing.T) {
	request := MatchRequest{
		CVText:         "Go and SQL",
		JobDescription: "Backend developer",
		RequiredKeywords: []string{
			"Go",
			"go",
			" GO ",
			"SQL",
		},
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalKeywords != 2 {
		t.Fatalf("expected 2 unique keywords, got %d", result.TotalKeywords)
	}

	if result.Score != 100 {
		t.Fatalf("expected score 100, got %d", result.Score)
	}
}

func TestMatchPhrase(t *testing.T) {
	request := MatchRequest{
		CVText:           "Built services using REST API design and Go.",
		JobDescription:   "Backend developer",
		RequiredKeywords: []string{"REST API", "Go"},
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 100 {
		t.Fatalf("expected score 100, got %d", result.Score)
	}
}

func TestMatchHandlesPunctuation(t *testing.T) {
	request := MatchRequest{
		CVText:           "Skills: Go, SQL; Git.",
		JobDescription:   "Backend developer",
		RequiredKeywords: []string{"Go", "SQL", "Git"},
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 100 {
		t.Fatalf("expected score 100, got %d", result.Score)
	}
}

func TestMatchDoesNotUseAccidentalSubstring(t *testing.T) {
	request := MatchRequest{
		CVText:           "I enjoy golfing and building software.",
		JobDescription:   "Backend developer",
		RequiredKeywords: []string{"Go"},
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 0 {
		t.Fatalf("expected score 0, got %d", result.Score)
	}

	if !reflect.DeepEqual(result.MissingKeywords, []string{"Go"}) {
		t.Fatalf("unexpected missing keywords: %v", result.MissingKeywords)
	}
}

func TestMatchReturnsMatchedAndMissingKeywords(t *testing.T) {
	request := MatchRequest{
		CVText:           "Go SQL Git",
		JobDescription:   "Backend developer",
		RequiredKeywords: []string{"Go", "SQL", "Git", "Docker"},
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Score != 75 {
		t.Fatalf("expected score 75, got %d", result.Score)
	}

	if result.TotalKeywords != 4 {
		t.Fatalf("expected total 4, got %d", result.TotalKeywords)
	}
}

func TestMatchExtractsFromVocabulary(t *testing.T) {
	request := MatchRequest{
		CVText:         "Experienced with Go, Git and SQL.",
		JobDescription: "We need Go, Git, SQL and Docker experience.",
	}

	result, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalKeywords != 4 {
		t.Fatalf("expected 4 extracted keywords, got %d", result.TotalKeywords)
	}

	if result.Score != 75 {
		t.Fatalf("expected score 75, got %d", result.Score)
	}
}

func TestMatchRejectsEmptyCV(t *testing.T) {
	_, err := Match(MatchRequest{
		CVText:           "   ",
		JobDescription:   "Go developer",
		RequiredKeywords: []string{"Go"},
	})

	if !errors.Is(err, ErrEmptyCV) {
		t.Fatalf("expected ErrEmptyCV, got %v", err)
	}
}

func TestMatchRejectsEmptyJobDescription(t *testing.T) {
	_, err := Match(MatchRequest{
		CVText:           "Go developer",
		JobDescription:   "",
		RequiredKeywords: []string{"Go"},
	})

	if !errors.Is(err, ErrEmptyJob) {
		t.Fatalf("expected ErrEmptyJob, got %v", err)
	}
}

func TestMatchRejectsZeroExtractedKeywords(t *testing.T) {
	_, err := Match(MatchRequest{
		CVText:         "Experienced software developer.",
		JobDescription: "Looking for an enthusiastic team member.",
	})

	if !errors.Is(err, ErrNoKeywordsFound) {
		t.Fatalf("expected ErrNoKeywordsFound, got %v", err)
	}
}

func TestMatchIsDeterministic(t *testing.T) {
	request := MatchRequest{
		CVText:           "Go SQL",
		JobDescription:   "Backend developer",
		RequiredKeywords: []string{"SQL", "Docker", "Go", "SQL"},
	}

	first, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := Match(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic results: %#v != %#v", first, second)
	}
}
