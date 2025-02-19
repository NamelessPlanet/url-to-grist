package scraper

import (
	"fmt"
	"strings"

	"url-to-grist/internal/types"

	"github.com/gocolly/colly/v2"
)

func FetchURLDetails(entry *types.Entry) (*types.Entry, error) {
	var err error

	c := colly.NewCollector()

	c.OnHTML("title", func(e *colly.HTMLElement) {
		if entry.Title == "" {
			entry.Title = e.Text
		}
	})

	c.OnHTML("meta[property]", func(e *colly.HTMLElement) {
		switch strings.ToLower(e.Attr("property")) {
		case "og:title":
			entry.Title = e.Attr("content")
		case "og:description":
			if entry.Summary == "" {
				entry.Summary = e.Attr("content")
			}
		case "og:article:author":
			entry.Byline = e.Attr("content")
		}
	})

	c.OnHTML("meta[name]", func(e *colly.HTMLElement) {
		switch strings.ToLower(e.Attr("name")) {
		case "og:title":
			entry.Title = e.Attr("content")
		case "og:description":
			if entry.Summary == "" {
				entry.Summary = e.Attr("content")
			}
		case "og:article:author":
			entry.Byline = e.Attr("content")
		case "article:author":
			entry.Byline = e.Attr("content")
		case "author":
			entry.Byline = e.Attr("content")
		case "octolytics-dimension-user_login":
			entry.Byline = e.Attr("content")
		}
	})

	c.OnHTML("span[itemprop='author']", func(e *colly.HTMLElement) {
		if entry.Byline == "" {
			linkChild := e.DOM.Find("link[itemprop='name']")
			if linkChild != nil {
				entry.Byline, _ = linkChild.Attr("content")
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.UserAgent = "url-to-grist - https://github.com/NamelessPlanet/url-to-grist"
	err = c.Visit(entry.URL)
	if err != nil {
		fmt.Printf("Failed to scrape details from '%s' - %s\n", entry.URL, err)
		return entry, err
	}

	return entry, err
}
