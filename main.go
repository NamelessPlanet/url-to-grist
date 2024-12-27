package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"url-to-grist/internal/ai"
	"url-to-grist/internal/grist"
	"url-to-grist/internal/scraper"
	"url-to-grist/internal/types"
)

var (
	port              string
	webserverPassword string
)

func init() {
	var ok bool
	port, ok = os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	webserverPassword = os.Getenv("WEBSERVER_PASSWORD")
}

func main() {
	if len(os.Args) == 1 {
		// HTTP Server
		if err := startServer(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// CLI
		for _, url := range os.Args[1:] {
			if entry, err := processURL(url, nil); err != nil {
				fmt.Printf("Failed to import '%s' - %s\n", url, err)
			} else {
				entryJSON, _ := json.MarshalIndent(entry, "", "  ")
				fmt.Printf("Successfully imported '%s' - '%s'\n", url, entryJSON)
			}
		}
		os.Exit(0)
	}
}

func startServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Enforce password access
		if webserverPassword != "" {
			userPass := r.URL.Query().Get("password")
			if userPass != webserverPassword {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		url := ""
		category := ""
		summary := ""
		featured := false
		sponsored := false

		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Has("url") {
				url = r.URL.Query().Get("url")
			}
		case http.MethodPost:
			data, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			url = string(data)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Populate optional extra metadata
		if r.URL.Query().Has("category") {
			category = r.URL.Query().Get("category")
		}
		if r.URL.Query().Has("summary") {
			summary = r.URL.Query().Get("summary")
		}
		if r.URL.Query().Has("featured") {
			featured = true
		}
		if r.URL.Query().Has("sponsored") {
			sponsored = true
		}

		if url != "" {
			entry := &types.Entry{
				URL:       url,
				Category:  category,
				Sponsored: sponsored,
				Featured:  featured,
				Summary:   summary,
			}

			entry, err := processURL(url, entry)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			entryJSON, _ := json.MarshalIndent(entry, "", "  ")
			_, _ = w.Write([]byte(fmt.Sprintf("Imported successfully - %s\n\n%s", url, string(entryJSON))))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func processURL(url string, baseEntry *types.Entry) (*types.Entry, error) {
	entry := &types.Entry{}
	if baseEntry != nil {
		entry = baseEntry
	}

	entry.URL = url
	entry.Year = time.Now().Format("2006")
	entry.Month = time.Now().Format("January")

	entry, err := scraper.FetchURLDetails(entry)
	if err != nil {
		fmt.Printf("Failed to fetch URL details - %s\n", err)
		return entry, err
	}

	aiSummary, err := ai.GenerateSummary(entry.URL)
	if err != nil {
		fmt.Printf("Failed to generate AI summary - %s\n", err)
	} else {
		entry.AISummary = aiSummary
	}

	entry, err = grist.Import(entry)
	if err != nil {
		fmt.Printf("Failed to fetch URL details - %s\n", err)
		return entry, err
	}

	return entry, err
}
