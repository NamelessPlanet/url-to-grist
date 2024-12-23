package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	geminiToken string

	prompt = "I will provide you with an URL and I would like you to generate a short summary or description of that page that would be suitable for a newsletter"
)

func init() {
	geminiToken = os.Getenv("GEMINI_TOKEN")

	if geminiToken == "" {
		fmt.Println("No GEMINI_TOKEN provided, disabling AI integration")
	}
}

func GenerateSummary(url string) (string, error) {
	var err error
	result := ""

	if geminiToken != "" {

		reqBody := Request{
			Contents: []RequestContent{
				RequestContent{Role: "user", Parts: []Part{Part{Text: prompt}}},
				RequestContent{Role: "user", Parts: []Part{Part{Text: url}}},
			},
		}
		reqJSON, err := json.Marshal(reqBody)
		if err != nil {
			return result, err
		}

		req, err := http.NewRequest(http.MethodPost, "https://generativelanguage.googleapis.com/v1/models/gemini-pro:generateContent", bytes.NewReader(reqJSON))
		if err != nil {
			return result, err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-goog-api-key", geminiToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return result, err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return result, err
		}

		response := Response{}
		json.Unmarshal(resBody, &response)

		if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
			return result, fmt.Errorf("Response empty - %s", string(resBody))
		}

		result = strings.TrimSpace(response.Candidates[0].Content.Parts[0].Text)
		result = strings.ReplaceAll(response.Candidates[0].Content.Parts[0].Text, "\"", "")
		result = strings.ReplaceAll(response.Candidates[0].Content.Parts[0].Text, "\n", "<br>")
	}

	return result, err
}
