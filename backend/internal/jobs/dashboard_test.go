package jobs

import (
	"testing"
	"time"
)

func loadLagos(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("Africa/Lagos")
	if err != nil {
		t.Fatalf("load Africa/Lagos location: %v", err)
	}
	return loc
}

func TestClassifyDeadlineOverdue(t *testing.T) {
	loc := loadLagos(t)
	today := time.Date(2026, 10, 20, 9, 0, 0, 0, loc)

	tests := []struct {
		name      string
		deadline  *string
		status    string
		want      string
	}{
		{"deadline yesterday", stringPointer("2026-10-19"), "Applied", "overdue"},
		{"deadline long ago", stringPointer("2026-01-01"), "Saved", "overdue"},
		{"deadline yesterday Offer", stringPointer("2026-10-19"), "Offer", "none"},
		{"deadline yesterday Rejected", stringPointer("2026-10-19"), "Rejected", "none"},
		{"no deadline", nil, "Applied", "none"},
		{"empty deadline", stringPointer(""), "Applied", "none"},
		{"invalid deadline", stringPointer("not-a-date"), "Applied", "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyDeadline(today, tt.deadline, tt.status); got != tt.want {
				t.Errorf("ClassifyDeadline = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassifyDeadlineUpcomingBoundaries(t *testing.T) {
	loc := loadLagos(t)
	today := time.Date(2026, 10, 20, 9, 0, 0, 0, loc)

	tests := []struct {
		name     string
		deadline *string
		status   string
		want     string
	}{
		{"deadline today", stringPointer("2026-10-20"), "Applied", "upcoming"},
		{"deadline tomorrow", stringPointer("2026-10-21"), "Saved", "upcoming"},
		{"deadline 7 days away", stringPointer("2026-10-27"), "Interview", "upcoming"},
		{"deadline 8 days away", stringPointer("2026-10-28"), "Applied", "none"},
		{"deadline today Offer", stringPointer("2026-10-20"), "Offer", "none"},
		{"deadline tomorrow Rejected", stringPointer("2026-10-21"), "Rejected", "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyDeadline(today, tt.deadline, tt.status); got != tt.want {
				t.Errorf("ClassifyDeadline = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestComputeDashboardStats(t *testing.T) {
	loc := loadLagos(t)
	now := time.Date(2026, 10, 20, 9, 0, 0, 0, loc)

	jobs := []Job{
		{ID: "1", Company: "A", JobTitle: "Engineer", Status: "Applied", Deadline: stringPointer("2026-10-19")},
		{ID: "2", Company: "B", JobTitle: "Designer", Status: "Interview", Deadline: stringPointer("2026-10-20")},
		{ID: "3", Company: "C", JobTitle: "PM", Status: "Saved", Deadline: stringPointer("2026-10-27")},
		{ID: "4", Company: "D", JobTitle: "Lead", Status: "Offer", Deadline: stringPointer("2026-10-15")},
		{ID: "5", Company: "E", JobTitle: "QA", Status: "Rejected", Deadline: stringPointer("2026-10-15")},
		{ID: "6", Company: "F", JobTitle: "Intern", Status: "Applied"}, // no deadline
	}

	stats := ComputeDashboardStats(now, jobs)

	if stats.TotalApplications != 6 {
		t.Errorf("total_applications = %d, want 6", stats.TotalApplications)
	}
	if stats.StatusCounts["Applied"] != 2 {
		t.Errorf("Applied count = %d, want 2", stats.StatusCounts["Applied"])
	}
	if stats.StatusCounts["Interview"] != 1 {
		t.Errorf("Interview count = %d, want 1", stats.StatusCounts["Interview"])
	}
	if stats.StatusCounts["Saved"] != 1 {
		t.Errorf("Saved count = %d, want 1", stats.StatusCounts["Saved"])
	}
	if stats.StatusCounts["Offer"] != 1 {
		t.Errorf("Offer count = %d, want 1", stats.StatusCounts["Offer"])
	}
	if stats.StatusCounts["Rejected"] != 1 {
		t.Errorf("Rejected count = %d, want 1", stats.StatusCounts["Rejected"])
	}
	if stats.TotalInterviews != 1 {
		t.Errorf("total_interviews = %d, want 1", stats.TotalInterviews)
	}
	if stats.TotalOffers != 1 {
		t.Errorf("total_offers = %d, want 1", stats.TotalOffers)
	}
	if stats.OverdueApplications != 1 {
		t.Errorf("overdue_applications = %d, want 1", stats.OverdueApplications)
	}
	if stats.UpcomingDeadlines != 2 {
		t.Errorf("upcoming_deadlines = %d, want 2", stats.UpcomingDeadlines)
	}
}

func TestComputeDashboardStatsIncludesZeroCounts(t *testing.T) {
	loc := loadLagos(t)
	now := time.Date(2026, 10, 20, 9, 0, 0, 0, loc)

	stats := ComputeDashboardStats(now, nil)

	for _, status := range []string{"Saved", "Applied", "Interview", "Offer", "Rejected"} {
		if _, ok := stats.StatusCounts[status]; !ok {
			t.Errorf("status_counts missing %q", status)
		}
		if stats.StatusCounts[status] != 0 {
			t.Errorf("status_counts[%q] = %d, want 0", status, stats.StatusCounts[status])
		}
	}
	if stats.TotalApplications != 0 {
		t.Errorf("total_applications = %d, want 0", stats.TotalApplications)
	}
}
