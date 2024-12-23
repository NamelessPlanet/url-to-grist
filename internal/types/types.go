package types

type Entry struct {
	URL             string `json:"url"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	Byline          string `json:"byline"`
	AISummary       string `json:"aiSummary"`
	Category        string `json:"category"`
	Year            string `json:"year"`
	Month           string `json:"month"`
	Featured        bool   `json:"featured"`
	Sponsored       bool   `json:"sponsored"`
	PostedToSocials bool   `json:"postedToSocials"`
}
