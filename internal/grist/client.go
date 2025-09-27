package grist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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

	records := Records{
		Record{
			Fields: Fields{
				URL:       entry.URL,
				Title:     strings.ReplaceAll(entry.Title, "\n", "<br>"),
				Summary:   strings.ReplaceAll(entry.Summary, "\n", "<br>"),
				Byline:    entry.Byline,
				Category:  entry.Category,
				Year:      entry.Year,
				Month:     entry.Month,
				AISummary: strings.ReplaceAll(entry.AISummary, "\n", "<br>"),
				Featured:  entry.Featured,
				Sponsored: entry.Sponsored,
			},
		},
	}

	jsonData, err := json.Marshal(records)
	if err != nil {
		return entry, err
	}

	req, _ := http.NewRequest("POST", gristTableURL, bytes.NewBuffer(jsonData))
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
