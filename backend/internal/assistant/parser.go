package assistant

import (
	"net/url"
	"strings"
)

func DetectPlatform(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return "unknown"
	}
	host := strings.ToLower(parsed.Host)
	if strings.Contains(host, "greenhouse") {
		return "greenhouse"
	}
	if strings.Contains(host, "lever") {
		return "lever"
	}
	if strings.Contains(host, "linkedin") {
		return "linkedin"
	}
	if strings.Contains(host, "indeed") {
		return "indeed"
	}
	if strings.Contains(host, "workday") {
		return "workday"
	}
	if strings.Contains(host, "ashbyhq") || strings.Contains(host, "ashby") {
		return "ashby"
	}
	return "unknown"
}

// Very simple read-only parser - extracts basic metadata from common patterns.
// For safety we don't execute JS; we rely on what's in HTML.
func ParseFromURL(raw string) (ParsedJob, error) {
	p := ParsedJob{SourceURL: raw, SourcePlatform: DetectPlatform(raw)}
	// Best effort - client can enhance; we keep it simple and safe
	// Try to extract from URL path for some platforms
	if p.SourcePlatform == "greenhouse" {
		// e.g. boards.greenhouse.io/company/jobs/123
		if strings.Contains(raw, "/jobs/") {
			p.Title = "Job (parse from page for details)"
		}
	}
	p.Description = ""
	p.Requirements = ""
	return p, nil
}
