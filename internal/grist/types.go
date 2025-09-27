package grist

type Fields struct {
	URL       string `json:"URL"`
	Title     string `json:"Title"`
	Summary   string `json:"Summary"`
	Byline    string `json:"Byline"`
	Category  string `json:"Category"`
	Year      string `json:"Year"`
	Month     string `json:"Month"`
	AISummary string `json:"AI_Summary"`
	Featured  bool   `json:"Featured"`
	Sponsored bool   `json:"Sponsored"`
}

type Record struct {
	Fields Fields `json:"fields"`
}

type Records []Record
