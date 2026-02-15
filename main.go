package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// ScrapedData represents the structured data extracted from a web page
type ScrapedData struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Keywords    string            `json:"keywords"`
	Links       []Link            `json:"links"`
	MetaTags    map[string]string `json:"meta_tags"`
	Images      []string          `json:"images"`
	Headings    []Heading         `json:"headings"`
}

// Link represents an extracted hyperlink
type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

// Heading represents an HTML heading (h1, h2, h3, etc.)
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()

	// Scraper endpoint
	mux.HandleFunc("/scrape", scrapeHandler)
	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)

	addr := ":8080"
	log.Printf("Web Scraper API starting on %s", addr)
	log.Printf("Example usage: curl 'http://localhost:8080/scrape?url=http://example.com'")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed")
		return
	}

	// Get URL parameter
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		sendError(w, http.StatusBadRequest, "missing_url", "URL parameter is required")
		return
	}

	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		sendError(w, http.StatusBadRequest, "invalid_url", "URL must be a valid http or https URL")
		return
	}

	// Scrape the web page
	data, err := scrapeWebPage(targetURL)
	if err != nil {
		log.Printf("Error scraping %s: %v", targetURL, err)
		sendError(w, http.StatusInternalServerError, "scrape_failed", fmt.Sprintf("Failed to scrape URL: %v", err))
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func scrapeWebPage(targetURL string) (*ScrapedData, error) {
	// Create HTTP client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebScraperBot/1.0)")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return nil, fmt.Errorf("content type is not HTML: %s", contentType)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract data
	data := &ScrapedData{
		URL:      targetURL,
		Links:    []Link{},
		MetaTags: make(map[string]string),
		Images:   []string{},
		Headings: []Heading{},
	}

	// Traverse the HTML tree and extract information
	extractData(doc, data, targetURL)

	return data, nil
}

func extractData(n *html.Node, data *ScrapedData, baseURL string) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			data.Title = getTextContent(n)
		case "meta":
			extractMetaTag(n, data)
		case "a":
			extractLink(n, data, baseURL)
		case "img":
			extractImage(n, data, baseURL)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			extractHeading(n, data)
		}
	}

	// Recursively process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractData(c, data, baseURL)
	}
}

func extractMetaTag(n *html.Node, data *ScrapedData) {
	var name, property, content string
	for _, attr := range n.Attr {
		switch attr.Key {
		case "name":
			name = attr.Val
		case "property":
			property = attr.Val
		case "content":
			content = attr.Val
		}
	}

	// Store meta tags by name or property
	if name != "" {
		data.MetaTags[name] = content
		// Extract specific meta tags
		if name == "description" {
			data.Description = content
		} else if name == "keywords" {
			data.Keywords = content
		}
	} else if property != "" {
		data.MetaTags[property] = content
	}
}

func extractLink(n *html.Node, data *ScrapedData, baseURL string) {
	var href string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	if href != "" {
		// Resolve relative URLs
		absoluteURL := resolveURL(baseURL, href)
		link := Link{
			Href: absoluteURL,
			Text: strings.TrimSpace(getTextContent(n)),
		}
		data.Links = append(data.Links, link)
	}
}

func extractImage(n *html.Node, data *ScrapedData, baseURL string) {
	var src string
	for _, attr := range n.Attr {
		if attr.Key == "src" {
			src = attr.Val
			break
		}
	}

	if src != "" {
		// Resolve relative URLs
		absoluteURL := resolveURL(baseURL, src)
		data.Images = append(data.Images, absoluteURL)
	}
}

func extractHeading(n *html.Node, data *ScrapedData) {
	level := 0
	switch n.Data {
	case "h1":
		level = 1
	case "h2":
		level = 2
	case "h3":
		level = 3
	case "h4":
		level = 4
	case "h5":
		level = 5
	case "h6":
		level = 6
	}

	text := strings.TrimSpace(getTextContent(n))
	if text != "" {
		heading := Heading{
			Level: level,
			Text:  text,
		}
		data.Headings = append(data.Headings, heading)
	}
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getTextContent(c)
	}
	return text
}

func resolveURL(baseURL, href string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	ref, err := url.Parse(href)
	if err != nil {
		return href
	}

	return base.ResolveReference(ref).String()
}

func sendError(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}
