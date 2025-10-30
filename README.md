# Grokipedia-API

A RESTful API service that provides programmatic access to Grokipedia content. This service allows you to fetch articles and search for content on Grokipedia using real-time web scraping with headless browser technology.

## Features

- 🔍 **Real-time Search** - Search articles using headless Chrome automation
- 📄 **Article Fetching** - Retrieve full article content by path
- 🚀 **Fast & Lightweight** - Go-based server with efficient scraping
- 🌐 **CORS-enabled** - Ready for web applications
- 📊 **JSON API** - Clean RESTful JSON responses
- 🐳 **Docker Support** - Easy deployment with Docker/Kubernetes

## Prerequisites

### For Docker (Recommended)
- Docker and Docker Compose

### For Local Development
- Go 1.21 or higher
- Chrome or Chromium browser (for headless search functionality)
- Internet connection (to access Grokipedia)

## Quick Start

### Option 1: Using Docker (Recommended)

Pull and run the pre-built image from GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/OWNER/REPO:latest

# Run the container
docker run -d -p 8080:8080 --name grokipedia-api ghcr.io/OWNER/REPO:latest

# Test the API
curl http://localhost:8080/health
```

Or use Docker Compose:

```bash
docker-compose up -d
```

### Option 2: Build from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/OWNER/REPO.git
   cd REPO
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

4. **Test the API:**
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Search for articles
   curl "http://localhost:8080/api/search?q=artificial+intelligence"

   # Get an article
   curl "http://localhost:8080/api/article/page/Machine_learning"
   ```

The server will start on `http://localhost:8080`

## Deployment

### Docker Deployment

The project includes automated Docker image building via GitHub Actions. Every push to the main branch automatically builds and publishes a Docker image to GitHub Container Registry.

**Available image tags:**
- `latest` - Latest build from main branch
- `v*.*.*` - Semantic version tags
- `main-<sha>` - Specific commit builds

**Deploy with Docker:**
```bash
docker run -d \
  -p 8080:8080 \
  --name grokipedia-api \
  --restart unless-stopped \
  ghcr.io/OWNER/REPO:latest
```

**Deploy with Docker Compose:**
```yaml
version: '3.8'
services:
  grokipedia-api:
    image: ghcr.io/OWNER/REPO:latest
    ports:
      - "8080:8080"
    restart: unless-stopped
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grokipedia-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: grokipedia-api
  template:
    metadata:
      labels:
        app: grokipedia-api
    spec:
      containers:
      - name: grokipedia-api
        image: ghcr.io/OWNER/REPO:latest
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: grokipedia-api
spec:
  selector:
    app: grokipedia-api
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## API Endpoints

### 1. Health Check

Check if the API is running.

**Endpoint:** `GET /health`

**Example:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "version": "1.0.0",
  "time": "2025-10-29T10:30:00Z"
}
```

### 2. Get Article

Fetch a specific article by its path.

**Endpoint:** `GET /api/article/{path}`

**Parameters:**
- `path` - The article path (e.g., "page/Artificial_intelligence", "page/Machine_learning")

**Example:**
```bash
curl http://localhost:8080/api/article/page/Artificial_intelligence
```

**Response:**
```json
{
  "title": "Artificial intelligence",
  "url": "https://grokipedia.com/page/Artificial_intelligence",
  "content": "Fundamentals\n\nDefining Artificial Intelligence...",
  "summary": "First paragraph summary...",
  "categories": []
}
```

### 3. Search Articles

Search for articles matching a query using real-time headless browser automation.

**Endpoint:** `GET /api/search?q={query}`

**Parameters:**
- `q` - Search query string (required)

**Example:**
```bash
curl "http://localhost:8080/api/search?q=machine+learning"
```

**Response:**
```json
{
  "query": "machine learning",
  "count": 12,
  "results": [
    {
      "title": "Machine learning",
      "url": "https://grokipedia.com/page/Machine_learning",
      "snippet": "No description available"
    },
    {
      "title": "Machine Learning Control",
      "url": "https://grokipedia.com/page/Machine_Learning_Control",
      "snippet": "No description available"
    },
    {
      "title": "Quantum machine learning",
      "url": "https://grokipedia.com/page/Quantum_machine_learning",
      "snippet": "No description available"
    }
  ]
}
```

**Note:** Search uses headless Chrome to execute JavaScript and retrieve real-time results. The first search may take 5-10 seconds as the browser initializes.

## Usage Examples

### Important Note

Grokipedia's URL structure may vary. To find the correct article path:
1. Visit https://grokipedia.com in your browser
2. Navigate to the article you want
3. Copy the path from the URL (everything after `grokipedia.com/`)
4. Use that path in the API

### Using cURL

```bash
# Health check
curl http://localhost:8080/health

# Get an article (replace with actual article path from Grokipedia)
curl http://localhost:8080/api/article/your-article-path

# Search for articles
curl "http://localhost:8080/api/search?q=quantum+computing"

# Example with a specific article
# First, find an article on https://grokipedia.com
# Then use its path, for example:
curl http://localhost:8080/api/article/example-topic
```

### Using JavaScript (fetch)

```javascript
// Get an article
fetch('http://localhost:8080/api/article/page/Artificial_intelligence')
  .then(response => response.json())
  .then(data => console.log(data));

// Search
fetch('http://localhost:8080/api/search?q=machine+learning')
  .then(response => response.json())
  .then(data => console.log(data));
```

### Using Python (requests)

```python
import requests

# Get an article
response = requests.get('http://localhost:8080/api/article/page/Artificial_intelligence')
article = response.json()
print(article['title'])

# Search
response = requests.get('http://localhost:8080/api/search', params={'q': 'machine learning'})
results = response.json()
print(f"Found {results['count']} results")
```

## Building for Production

Build a standalone binary:

```bash
# For current platform
go build -o grokipedia-api

# For Linux
GOOS=linux GOARCH=amd64 go build -o grokipedia-api-linux

# For Windows
GOOS=windows GOARCH=amd64 go build -o grokipedia-api.exe

# For macOS
GOOS=darwin GOARCH=amd64 go build -o grokipedia-api-mac
```

Run the binary:

```bash
./grokipedia-api
```

## Configuration

You can customize the server by modifying these constants in `main.go`:

- `baseURL` - The Grokipedia base URL (default: "https://grokipedia.com")
- `port` - Server port (default: "8080")

To change the port, modify the `port` variable in the `main()` function or set it via environment variable:

```bash
PORT=3000 go run main.go
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful request
- `400 Bad Request` - Missing or invalid parameters
- `404 Not Found` - Article not found
- `500 Internal Server Error` - Server error

Error response format:
```json
{
  "error": "Bad Request",
  "message": "Article path is required"
}
```

## Technical Details

### Architecture
- **Language:** Go 1.21+
- **Web Scraping:** goquery for HTML parsing
- **Headless Browser:** chromedp for JavaScript-rendered content
- **Router:** Gorilla Mux
- **Containerization:** Multi-stage Docker build with Alpine Linux

### How It Works
1. **Article Fetching:** Uses simple HTTP requests to fetch server-side rendered article pages
2. **Search:** Uses headless Chrome to execute JavaScript and retrieve real-time search results from Grokipedia
3. **Deduplication:** JavaScript-based deduplication ensures unique search results
4. **URL Construction:** Automatically constructs proper `/page/{slug}` URLs from article titles

## Limitations

- This API scrapes content from Grokipedia's public website
- Search functionality requires Chrome/Chromium (included in Docker image)
- First search request may take 5-10 seconds as the headless browser initializes
- Rate limiting may apply based on Grokipedia's server policies
- Content structure may change if Grokipedia updates their website
- No authentication is currently implemented

## Development

### Project Structure

```
.
├── main.go       # Main application code
├── go.mod        # Go module dependencies
└── README.md     # This file
```

### Adding New Features

To add new endpoints, follow this pattern:

1. Create a handler function
2. Register the route in `main()`
3. Update this README with documentation

## License

This project is provided as-is for educational and personal use. Please respect Grokipedia's terms of service when using this API.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Troubleshooting

### Server won't start
- Check if port 8080 is already in use
- Try a different port by modifying the code

### Articles not loading
- Verify internet connection
- Check if Grokipedia.com is accessible
- The website structure may have changed

### Search returns no results
- Ensure Chrome/Chromium is installed (or use Docker image which includes it)
- Try different search terms
- The search functionality depends on Grokipedia's search page structure
- Check server logs for headless browser errors

## Support

For issues or questions, please open an issue on the project repository.

