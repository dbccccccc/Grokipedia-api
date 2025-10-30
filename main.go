package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
)

const (
	defaultBaseURL = "https://grokipedia.com"
	defaultPort    = "8080"
)

var (
	baseURL string
	port    string
)

// Article represents a Grokipedia article
type Article struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Content     string   `json:"content"`
	Summary     string   `json:"summary"`
	Categories  []string `json:"categories,omitempty"`
	LastUpdated string   `json:"last_updated,omitempty"`
}

// SearchResult represents a search result
type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

// fetchHTML fetches HTML content from a URL
func fetchHTML(urlStr string) (*goquery.Document, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Grokipedia-API-Client/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// getArticle fetches and parses a Grokipedia article
func getArticle(articlePath string) (*Article, error) {
	// Ensure the path starts with /
	if !strings.HasPrefix(articlePath, "/") {
		articlePath = "/" + articlePath
	}

	fullURL := baseURL + articlePath
	log.Printf("Fetching article from URL: %s", fullURL)

	doc, err := fetchHTML(fullURL)
	if err != nil {
		return nil, err
	}

	article := &Article{
		URL: fullURL,
	}

	// Extract title
	article.Title = doc.Find("h1").First().Text()
	if article.Title == "" {
		article.Title = doc.Find("title").First().Text()
	}

	// Extract main content, walking the rendered article structure
	var contentParts []string
	lastLine := ""

	addContent := func(sel *goquery.Selection, candidateForSummary bool) {
		clean := sel.Clone()
		clean.Find("button, svg, style, script").Remove()

		text := strings.TrimSpace(clean.Text())
		if text == "" {
			return
		}

		text = strings.Join(strings.Fields(text), " ")
		if utf8.RuneCountInString(text) < 3 {
			return
		}

		if text == lastLine {
			return
		}

		contentParts = append(contentParts, text)
		lastLine = text

		if candidateForSummary && article.Summary == "" && utf8.RuneCountInString(text) > 50 {
			article.Summary = text
		}
	}

	processContent := func(root *goquery.Selection) {
		root.Find("*").Each(func(i int, s *goquery.Selection) {
			nodeName := goquery.NodeName(s)

			switch nodeName {
			case "h2", "h3", "h4", "h5", "h6", "blockquote", "pre", "p", "li":
				addContent(s, nodeName == "p" || nodeName == "blockquote")
			case "span":
				classAttr, _ := s.Attr("class")
				if strings.Contains(classAttr, "katex") || strings.Contains(classAttr, "sr-only") {
					return
				}

				if strings.Contains(classAttr, "break-words") || strings.Contains(classAttr, "leading-7") {
					addContent(s, true)
				}
			}
		})
	}

	articleRoot := doc.Find("article")
	if articleRoot.Length() == 0 {
		articleRoot = doc.Find("main")
	}
	if articleRoot.Length() > 0 {
		processContent(articleRoot)
	}

	if len(contentParts) == 0 {
		processContent(doc.Selection)
	}

	article.Content = strings.Join(contentParts, "\n\n")

	// Fall back to meta description for summary if needed
	if article.Summary == "" {
		if desc, ok := doc.Find(`meta[name="description"]`).Attr("content"); ok {
			desc = strings.TrimSpace(desc)
			if desc != "" {
				article.Summary = desc
			}
		}
	}

	if article.Summary == "" {
		if ogDesc, ok := doc.Find(`meta[property="og:description"]`).Attr("content"); ok {
			ogDesc = strings.TrimSpace(ogDesc)
			if ogDesc != "" {
				article.Summary = ogDesc
			}
		}
	}

	// Capture last updated timestamp if provided
	if modified, ok := doc.Find(`meta[property="article:modified_time"]`).Attr("content"); ok {
		modified = strings.TrimSpace(modified)
		if modified != "" {
			article.LastUpdated = modified
		}
	}

	// Extract categories if available
	doc.Find(".categories a, .category a").Each(func(i int, s *goquery.Selection) {
		category := strings.TrimSpace(s.Text())
		if category != "" {
			article.Categories = append(article.Categories, category)
		}
	})

	return article, nil
}

// searchArticles searches for articles on Grokipedia using headless Chrome
// This function uses chromedp to execute JavaScript and get real-time search results
func searchArticles(query string) ([]SearchResult, error) {
	log.Printf("Starting headless browser search for: %s", query)

	// Create allocator options with headless mode
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create Chrome context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set timeout (30 seconds should be enough)
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	searchURL := fmt.Sprintf("%s/search?q=%s", baseURL, query)
	log.Printf("Navigating to: %s", searchURL)

	var results []SearchResult
	var htmlContent string

	// Run chromedp tasks
	err := chromedp.Run(ctx,
		// Navigate to search page
		chromedp.Navigate(searchURL),

		// Wait for search results container to appear
		chromedp.WaitVisible(`main`, chromedp.ByQuery),

		// Wait a bit for JavaScript to render results
		chromedp.Sleep(3*time.Second),

		// Get HTML for debugging
		chromedp.OuterHTML(`main`, &htmlContent, chromedp.ByQuery),

		// Extract search results using JavaScript
		chromedp.Evaluate(`
			(function() {
				const results = [];
				const seen = new Set(); // Track unique titles to avoid duplicates

				// Find all search result items (they're in divs with cursor-pointer class)
				const items = document.querySelectorAll('main div.cursor-pointer');

				console.log('Found ' + items.length + ' search result items');

				items.forEach(item => {
					// Find the title span
					const titleSpan = item.querySelector('span.line-clamp-1 span');
					if (!titleSpan) return;

					const title = titleSpan.textContent.trim();
					if (!title || seen.has(title)) return; // Skip duplicates

					seen.add(title);

					// Construct the page URL from the title
					// Convert title to URL slug (replace spaces with underscores)
					const slug = title.replace(/ /g, '_');
					const url = 'https://grokipedia.com/page/' + encodeURIComponent(slug);

					// Try to find snippet/description
					let snippet = '';
					const paragraphs = item.querySelectorAll('p');
					for (let p of paragraphs) {
						const text = p.textContent.trim();
						if (text && text.length > 20) {
							snippet = text.substring(0, 200);
							break;
						}
					}

					results.push({
						title: title,
						url: url,
						snippet: snippet || 'No description available'
					});
				});

				return results.slice(0, 20); // Return top 20 unique results
			})();
		`, &results),
	)

	if err != nil {
		log.Printf("Headless browser error: %v", err)
		return nil, fmt.Errorf("headless browser search failed: %w", err)
	}

	log.Printf("HTML content length: %d bytes", len(htmlContent))
	log.Printf("Found %d search results for query: %s", len(results), query)

	return results, nil
}

// Handlers

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
		Time:    time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articlePath := vars["path"]

	if articlePath == "" {
		sendError(w, http.StatusBadRequest, "Article path is required")
		return
	}

	article, err := getArticle(articlePath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch article: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		sendError(w, http.StatusBadRequest, "Search query parameter 'q' is required")
		return
	}

	results, err := searchArticles(query)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Search failed: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"query":   query,
		"count":   len(results),
		"results": results,
	})
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func init() {
	// Load configuration from environment variables
	baseURL = os.Getenv("GROKIPEDIA_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
}

func main() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/article/{path:.*}", getArticleHandler).Methods("GET")
	r.HandleFunc("/api/search", searchHandler).Methods("GET")

	// Apply middleware
	handler := corsMiddleware(loggingMiddleware(r))

	log.Printf("Starting Grokipedia API server")
	log.Printf("Base URL: %s", baseURL)
	log.Printf("Port: %s", port)
	log.Printf("Endpoints:")
	log.Printf("  GET /health - Health check")
	log.Printf("  GET /api/article/{path} - Get article by path")
	log.Printf("  GET /api/search?q={query} - Search articles")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
