# url-to-grist

Takes an URL and imports it into Grist for later processing

## Features

* CLI import
  ```
  GRIST_TABLE_URL="EXAMPLE URL" \
  GRIST_API_KEY="EXAMPLE KEY" \
  GEMINI_TOKEN="EXAMPLE TOKEN" \
    go run . "https://example.com"
  ```
* Web Server - `go run .`
  * Can be triggered with by making a request to `http://localhost:8000/?url=https://example.com`
