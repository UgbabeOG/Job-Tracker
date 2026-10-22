package jobs

import (
	"strings"
	"time"
)

var lagosLocation = mustLoadLocation("Africa/Lagos")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

// nowInLagos returns the current moment expressed in the Africa/Lagos time zone.
// The handler uses this so "today" is the Lagos calendar date.
func nowInLagos() time.Time {
	return time.Now().In(lagosLocation)
}

// today returns today's calendar date in Africa/Lagos as YYYY-MM-DD.
func today() string {
	return nowInLagos().Format("2006-01-02")
}

// ClassifyDeadline buckets a single job's deadline relative to `now`.
// Returns "overdue", "upcoming", or "none".
func ClassifyDeadline(now time.Time, deadline *string, status string) string {
	if status == "Offer" || status == "Rejected" {
		return "none"
	}
	if deadline == nil || strings.TrimSpace(*deadline) == "" {
		return "none"
	}

	deadlineDate, err := time.Parse("2006-01-02", *deadline)
	if err != nil {
		return "none"
	}

	todayDate := time.Date(
		now.In(lagosLocation).Year(),
		now.In(lagosLocation).Month(),
		now.In(lagosLocation).Day(),
		0, 0, 0, 0, time.UTC,
	)

	if deadlineDate.Before(todayDate) {
		return "overdue"
	}

	sevenDaysLater := todayDate.AddDate(0, 0, 7)
	if !deadlineDate.After(sevenDaysLater) {
		return "upcoming"
	}

	return "none"
}

// DashboardStats is the shape returned by GET /api/dashboard.
type DashboardStats struct {
	TotalApplications   int            `json:"total_applications"`
	StatusCounts        map[string]int `json:"status_counts"`
	UpcomingDeadlines   int            `json:"upcoming_deadlines"`
	OverdueApplications int            `json:"overdue_applications"`
	TotalInterviews     int            `json:"total_interviews"`
	TotalOffers         int            `json:"total_offers"`
}

// ComputeDashboardStats aggregates stats from a list of jobs using `now` as "today".
func ComputeDashboardStats(now time.Time, jobs []Job) DashboardStats {
	stats := DashboardStats{
		StatusCounts: map[string]int{
			"Saved":     0,
			"Applied":   0,
			"Interview": 0,
			"Offer":     0,
			"Rejected":  0,
		},
	}

	for _, job := range jobs {
		stats.TotalApplications++
		stats.StatusCounts[job.Status]++

		if job.Status == "Interview" {
			stats.TotalInterviews++
		}
		if job.Status == "Offer" {
			stats.TotalOffers++
		}

		switch ClassifyDeadline(now, job.Deadline, job.Status) {
		case "overdue":
			stats.OverdueApplications++
		case "upcoming":
			stats.UpcomingDeadlines++
		}
	}

	return stats
}
