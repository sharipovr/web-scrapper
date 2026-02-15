# Project 4: Web Scraper API

A Go web service that scrapes web pages and returns structured data including title, links, meta tags, images, and headings.

## 🎯 Learning Objectives

- HTTP client operations with `net/http`
- HTML parsing with `golang.org/x/net/html`
- Context-based timeout handling
- URL validation and manipulation
- Structured data extraction from HTML
- Error handling for external HTTP requests

## 📋 Features

- **URL Scraping**: Accepts a URL parameter and fetches the page
- **HTML Parsing**: Extracts structured information from HTML
- **Data Extraction**:
  - Page title
  - Meta tags (description, keywords, Open Graph, etc.)
  - All hyperlinks with anchor text
  - All images
  - Headings (h1-h6) with hierarchy
- **Timeout Handling**: 15-second timeout for requests
- **Error Handling**: Comprehensive error responses
- **URL Resolution**: Converts relative URLs to absolute URLs

## 🏗️ Project Structure

```
project4-web-scraper/
├── main.go      # Main application with scraper logic
├── go.mod       # Go module definition
└── README.md    # This file
```

## 🚀 Running the Application

### Start the server:
```bash
go mod download
go run main.go
```

The server will start on `http://localhost:8080`

### Build and run:
```bash
go build -o web-scraper
./web-scraper
```

## 📡 API Endpoints

### 1. Scrape Web Page
**GET** `/scrape?url=<target_url>`

Scrapes the provided URL and returns structured data.

**Query Parameters:**
- `url` (required): The URL to scrape (must be http or https)

**Example Request:**
```bash
curl "http://localhost:8080/scrape?url=http://example.com"
```

**Success Response (200 OK):**
```json
{
  "url": "http://example.com",
  "title": "Example Domain",
  "description": "Example description from meta tag",
  "keywords": "example, demo, test",
  "links": [
    {
      "href": "http://example.com/about",
      "text": "About Us"
    }
  ],
  "meta_tags": {
    "description": "Example description",
    "viewport": "width=device-width, initial-scale=1",
    "og:title": "Example Domain"
  },
  "images": [
    "http://example.com/logo.png"
  ],
  "headings": [
    {
      "level": 1,
      "text": "Example Domain"
    },
    {
      "level": 2,
      "text": "Welcome"
    }
  ]
}
```

**Error Responses:**

Missing URL (400 Bad Request):
```json
{
  "error": "missing_url",
  "message": "URL parameter is required"
}
```

Invalid URL (400 Bad Request):
```json
{
  "error": "invalid_url",
  "message": "URL must be a valid http or https URL"
}
```

Scrape Failed (500 Internal Server Error):
```json
{
  "error": "scrape_failed",
  "message": "Failed to scrape URL: <error details>"
}
```

### 2. Health Check
**GET** `/health`

Returns the service health status.

**Example Request:**
```bash
curl http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "time": "2026-02-14T10:30:00Z"
}
```

## 🧪 Testing Examples

### Test with Example.com:
```bash
curl "http://localhost:8080/scrape?url=http://example.com"
```

### Test with a news site:
```bash
curl "http://localhost:8080/scrape?url=https://news.ycombinator.com"
```

### Test error handling (missing URL):
```bash
curl http://localhost:8080/scrape
```

### Test error handling (invalid URL):
```bash
curl "http://localhost:8080/scrape?url=not-a-valid-url"
```

### Test with jq for pretty output:
```bash
curl -s "http://localhost:8080/scrape?url=http://example.com" | jq '.'
```

### Count extracted elements:
```bash
# Count links
curl -s "http://localhost:8080/scrape?url=http://example.com" | jq '.links | length'

# Count images
curl -s "http://localhost:8080/scrape?url=http://example.com" | jq '.images | length'

# Count headings
curl -s "http://localhost:8080/scrape?url=http://example.com" | jq '.headings | length'
```

## 🔑 Key Concepts Demonstrated

### 1. HTTP Client with Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
client := &http.Client{Timeout: 15 * time.Second}
resp, err := client.Do(req)
```

### 2. HTML Parsing
```go
doc, err := html.Parse(strings.NewReader(string(body)))
// Traverse the HTML tree recursively
extractData(doc, data, targetURL)
```

### 3. URL Validation
```go
parsedURL, err := url.Parse(targetURL)
if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
    // Invalid URL
}
```

### 4. Relative to Absolute URL Resolution
```go
base, _ := url.Parse(baseURL)
ref, _ := url.Parse(href)
absoluteURL := base.ResolveReference(ref).String()
```

### 5. Error Handling
- URL validation errors
- HTTP request errors
- Timeout errors
- HTML parsing errors
- Non-HTML content type errors

## 📚 Standard Library Packages Used

- `net/http` - HTTP client and server
- `encoding/json` - JSON encoding/decoding
- `context` - Timeout and cancellation
- `net/url` - URL parsing and manipulation
- `io` - Reading response body
- `strings` - String manipulation
- `time` - Timeout configuration
- `golang.org/x/net/html` - HTML parsing

## 🎓 Interview Talking Points

1. **Timeout Handling**: Implemented both context timeout and client timeout for robust control
2. **URL Resolution**: Handles relative URLs properly by converting them to absolute URLs
3. **Content Type Validation**: Checks that the response is HTML before parsing
4. **User Agent**: Sets a proper User-Agent header to avoid being blocked
5. **Resource Cleanup**: Uses `defer` for proper response body cleanup
6. **Error Propagation**: Uses `fmt.Errorf` with `%w` for error wrapping
7. **Tree Traversal**: Recursive function to traverse HTML DOM tree
8. **Structured Data**: Returns well-organized JSON with different data types
9. **HTTP Status Codes**: Proper use of 400, 405, 500 status codes
10. **Security**: Only allows http/https protocols to prevent potential security issues

## 🔄 Possible Enhancements

1. Add rate limiting to prevent abuse
2. Implement caching to avoid re-scraping the same URLs
3. Add robots.txt checking for ethical scraping
4. Support for JavaScript-rendered pages (would need external tools)
5. Extract social media metadata (Open Graph, Twitter Cards)
6. Add pagination support for multi-page scraping
7. Extract structured data (JSON-LD, microdata)
8. Add request queuing for batch scraping
9. Support custom HTTP headers
10. Add metrics collection (response times, success rates)

## ⚠️ Production Considerations

1. **Rate Limiting**: Implement rate limiting per domain
2. **Robots.txt**: Respect robots.txt files
3. **Caching**: Cache scraped data to reduce load
4. **Error Retry**: Implement exponential backoff for retries
5. **Max Response Size**: Limit response body size to prevent memory issues
6. **Concurrency Limits**: Control concurrent scraping operations
7. **User Agent**: Use descriptive, identifiable user agent
8. **Legal/Ethical**: Ensure compliance with terms of service and laws
