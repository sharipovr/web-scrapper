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

## 🎓 Technical Discussion Points

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

## 🔧 Troubleshooting & Common Issues (Q&A)

### Q1: Why am I getting "403 Forbidden" errors?

**Answer:** A 403 error means the website is blocking your scraper. This happens because:

1. **Anti-bot services** like Cloudflare, Akamai, or AWS WAF detect automated traffic
2. The website actively blocks scrapers to protect their content
3. Your requests look suspicious (missing headers, bot-like behavior)

**What the scraper does:**
- Detects which anti-bot service is blocking you (Cloudflare, Akamai, etc.)
- Logs the service name so you know what you're dealing with
- Shows response headers to help diagnose the issue

**Example output:**
```
❌ Request blocked - Status: 403
🛡️  Anti-bot service detected: Cloudflare
```

**Solutions:**
- Some sites will always block scrapers - this is expected
- For sites with basic blocking: the scraper uses realistic browser headers
- For advanced blocking (JavaScript challenges): you'd need a headless browser
- Check if the site offers an API as an alternative

---

### Q2: Why am I getting "context deadline exceeded" errors?

**Answer:** This means the request **timed out** - it took longer than 15 seconds (our configured timeout).

**Common causes:**
1. **Website is very slow** or overloaded
2. **Anti-bot delays** - some sites deliberately slow down suspected bots
3. **Network issues** - poor connection or packet loss
4. **Geo-blocking** - slower response for certain regions
5. **Challenge pages** - Cloudflare presenting a challenge that never completes

**What happens:**
```
Error: failed to fetch URL: Get "https://example.com": context deadline exceeded
```

**Solutions:**
- If it's a legitimate slow site, increase the timeout
- If it's anti-bot behavior, the site might be presenting a challenge page
- Try the site in a regular browser first to see if it loads

---

### Q3: How can I tell what's blocking my scraper?

**Answer:** The scraper has built-in **anti-bot detection** that identifies common blocking services:

**Detected services:**
- **Cloudflare** - Most common, very sophisticated
- **Akamai** - Enterprise-level protection
- **AWS WAF** - Amazon's Web Application Firewall
- **Imperva/Incapsula** - Security service
- **DataDome** - Bot detection platform
- **Sucuri** - Security scanner
- **PerimeterX** - Bot detection

**How it works:**
The scraper checks response headers for telltale signs:
```go
// Looks for headers like:
"CF-Ray"           // Cloudflare
"X-Akamai-*"       // Akamai
"X-Amzn-*"         // AWS WAF
```

**Running with debug mode** shows all headers:
```bash
./project4-web-scraper -debug
```

---

### Q4: Why am I getting empty results but status 200 OK?

**Answer:** This was a tricky bug! You're getting **compressed (gzipped) data** but the scraper wasn't decompressing it.

**What was happening:**
1. We set `Accept-Encoding: gzip, deflate, br` header
2. Server sent compressed HTML (binary data)
3. Go's automatic decompression was disabled (we overrode it)
4. Parser tried to parse binary gibberish → extracted nothing

**The garbled output looked like:**
```
�j5���C3��HY8� �c���*�$...
```

**The fix:**
- **Removed** the `Accept-Encoding` header
- Let Go's `http.Client` handle compression automatically
- Now it requests compression AND decompresses transparently

**Key lesson:** In Go, don't manually set `Accept-Encoding` unless you'll handle decompression yourself!

---

### Q5: Did removing headers make the scraper slower?

**Answer:** No! It actually made it **faster and more reliable**.

**What's happening now:**
1. ✅ Go automatically adds `Accept-Encoding: gzip`
2. ✅ Server sends compressed data (smaller, faster transfer)
3. ✅ Go automatically decompresses it (transparent)
4. ✅ Result: Fast transfer + working extraction

**Before (broken):**
- Compressed transfer ✅
- Manual decompression ❌
- Result: Fast but broken

**Now (working):**
- Compressed transfer ✅
- Auto decompression ✅
- Result: Fast AND working!

---

### Q6: How do I use debug mode?

**Answer:** Debug mode gives you detailed information about what's happening.

**Without debug (production/clean mode):**
```bash
./project4-web-scraper
```
Shows only:
- Errors when they occur
- Anti-bot detection warnings
- Minimal clean output

**With debug (troubleshooting mode):**
```bash
./project4-web-scraper -debug
```
Shows everything:
- ✅ Success confirmations
- 📊 Response sizes
- 📄 Content types
- 🔍 HTML preview (first 500 characters)
- 📈 Extraction statistics (how many links, images, etc.)
- 📋 Full HTTP response headers

**Use debug mode when:**
- A site isn't working and you don't know why
- You want to see what HTML is actually being received
- You're learning and want to understand the process
- You need to diagnose anti-bot blocks

**Example debug output:**
```
✅ Successfully fetched https://example.com
📊 Response size: 1256 bytes
📄 Content-Type: text/html; charset=UTF-8
🔍 HTML Preview:
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
...
📈 Extraction results: Title="Example Domain", Links=1, Images=0, Headings=2, MetaTags=4
```

---

### Q7: Why do some sites work and others don't?

**Answer:** Web scraping success depends on the site's protection level:

**✅ Works well with:**
- Simple static HTML sites
- Sites without anti-bot protection
- Sites with basic security (our browser headers bypass these)
- Public content sites that don't mind scrapers

**⚠️ Might work with:**
- Sites with moderate protection (hit or miss)
- Sites that rate-limit but don't fully block
- Sites with simple bot detection

**❌ Won't work with:**
- Sites protected by Cloudflare's JavaScript challenge
- Sites requiring login/authentication
- Sites with heavy JavaScript rendering (single-page apps)
- Sites actively blocking all automated access
- Sites checking for browser fingerprints

**For advanced cases, you'd need:**
- Headless browser (like `chromedp` in Go)
- Proxy rotation
- Cookie/session management
- JavaScript execution

**For technical discussions:** Understanding these limitations and trade-offs is important!

---

### Q8: Is my ISP blocking the requests?

**Answer:** Very unlikely! Here's why:

**What ISPs typically do:**
- Block entire domains via DNS (you'd get connection errors, not empty results)
- Throttle speed (you'd see slowness/timeouts, not empty data)
- They rarely modify HTTP content

**What you're more likely experiencing:**
1. **JavaScript-rendered content** - Site loads content after page load
2. **Bot detection** - Site serves different HTML to scrapers
3. **Compression issues** - Like we found with the gzip problem
4. **API-based content** - Content loaded separately from main page

**How to test:**
1. Try the same site in your browser
2. Check the HTML preview in debug mode
3. Compare with "View Source" in your browser
4. Use developer tools to see network requests

---

## 🎯 Key Technical Takeaways

When discussing this project, highlight these learning points:

1. **Problem Diagnosis**: Used systematic debugging (headers, logs, preview) to find gzip issue
2. **Understanding Go's HTTP Client**: Learned about automatic compression handling
3. **Real-world Constraints**: Experienced bot detection, timeouts, and blocking
4. **Feature vs Debug**: Implemented debug flag for troubleshooting without cluttering production logs
5. **Anti-bot Detection**: Built intelligent detection of protection services
6. **Error Context**: Provided meaningful error messages that help users understand failures
7. **Trade-offs**: Understood limitations (JavaScript rendering, Cloudflare) and alternatives
