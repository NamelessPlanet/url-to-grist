package utils

import (
	"net/url"
	"strings"
)

// StripAnalytics takes a raw URL string, removes all known analytics‑tracking
// query parameters, and returns the cleaned URL. If the input cannot be parsed,
// the original string and the parsing error are returned.
func StripAnalytics(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}

	q := u.Query()
	for key := range q {
		if isAnalyticsKey(key) {
			q.Del(key)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// isAnalyticsKey returns true if the query key is a typical analytics/tracking
// parameter.  The list can be extended as needed.
func isAnalyticsKey(key string) bool {
	// Keys that are matched exactly
	exact := map[string]struct{}{
		"gclid":  {},
		"fbclid": {},
		"aclid":  {},
		"dclid":  {},
		"mc_cid": {},
		"mc_eid": {},
	}

	if _, ok := exact[key]; ok {
		return true
	}

	// Keys that start with these prefixes
	prefixes := []string{"utm_", "utm-"}

	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
