package grist

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"url-to-grist/internal/types"
)

var (
	gristTableURL string
	gristAPIKey   string
)

func init() {
	gristTableURL = os.Getenv("GRIST_TABLE_URL")
	gristAPIKey = os.Getenv("GRIST_API_KEY")

	if gristTableURL == "" || gristAPIKey == "" {
		panic("GRIST_TABLE_ID and GRIST_API_KEY must be provided")
	}
}

func Import(entry *types.Entry) (*types.Entry, error) {
	var err error

	c := http.Client{}

	var jsonData = fmt.Sprintf(`{
		"records": [
			{
				"fields": {
					"URL": "%s",
					"Title": "%s",
					"Summary": "%s",
					"Byline": "%s",
					"Category": "%s",
					"Year": "%s",
					"Month": "%s",
					"AI_Summary": "%s",
					"Featured": %t,
					"Sponsored": %t
				}
			}
		]
	}`,
		entry.URL,
		strings.ReplaceAll(entry.Title, "\n", "<br>"),
		strings.ReplaceAll(entry.Summary, "\n", "<br>"),
		entry.Byline,
		entry.Category,
		entry.Year,
		entry.Month,
		strings.ReplaceAll(entry.AISummary, "\n", "<br>"),
		entry.Featured,
		entry.Sponsored,
	)

	req, _ := http.NewRequest("POST", gristTableURL, bytes.NewBuffer([]byte(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gristAPIKey))

	resp, err := c.Do(req)
	if err != nil {
		return entry, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return entry, err
	}

	fmt.Printf("Grist response: %s\n", resp.Status)
	if resp.StatusCode >= 400 {
		fmt.Println(string(respBody))
		fmt.Println(jsonData)
		return entry, fmt.Errorf("Unexpected status response from Grist")
	}

	return entry, err
}
