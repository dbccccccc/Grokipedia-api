# Grokipedia REST API Documentation

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, no authentication is required.

## Response Format

All responses are in JSON format.

### Success Response

```json
{
  "field1": "value1",
  "field2": "value2"
}
```

### Error Response

```json
{
  "error": "Error Type",
  "message": "Detailed error message"
}
```

## HTTP Status Codes

- `200 OK` - Request successful
- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Endpoints

### 1. Health Check

Check the API server status.

**Endpoint:** `GET /health`

**Parameters:** None

**Response:**

```json
{
  "status": "ok",
  "version": "1.0.0",
  "time": "2025-10-29T10:30:00Z"
}
```

**Example:**

```bash
curl http://localhost:8080/health
```

**PowerShell:**

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
```

---

### 2. Get Article

Retrieve a specific article from Grokipedia.

**Endpoint:** `GET /api/article/{path}`

**Path Parameters:**

| Parameter | Type   | Required | Description                                    |
|-----------|--------|----------|------------------------------------------------|
| path      | string | Yes      | The article path from Grokipedia URL           |

**Response:**

```json
{
  "title": "Article Title",
  "url": "https://grokipedia.com/article-path",
  "content": "Full article content with paragraphs separated by newlines...",
  "summary": "First paragraph or summary of the article",
  "categories": ["Category1", "Category2"],
  "last_updated": "2025-10-29"
}
```

**Response Fields:**

| Field        | Type     | Description                                      |
|--------------|----------|--------------------------------------------------|
| title        | string   | Article title                                    |
| url          | string   | Full URL to the article on Grokipedia           |
| content      | string   | Full article content                             |
| summary      | string   | Article summary (usually first paragraph)        |
| categories   | string[] | List of categories (if available)                |
| last_updated | string   | Last update date (if available)                  |

**Example:**

```bash
# Replace 'your-article-path' with actual path from Grokipedia
curl http://localhost:8080/api/article/your-article-path
```

**PowerShell:**

```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/article/your-article-path" -Method Get
$response | ConvertTo-Json
```

**Python:**

```python
import requests

response = requests.get('http://localhost:8080/api/article/your-article-path')
article = response.json()
print(f"Title: {article['title']}")
print(f"Summary: {article['summary']}")
```

**JavaScript:**

```javascript
fetch('http://localhost:8080/api/article/your-article-path')
  .then(response => response.json())
  .then(article => {
    console.log('Title:', article.title);
    console.log('Summary:', article.summary);
  });
```

---

### 3. Search Articles

Search for articles on Grokipedia.

**Endpoint:** `GET /api/search`

**Query Parameters:**

| Parameter | Type   | Required | Description                    |
|-----------|--------|----------|--------------------------------|
| q         | string | Yes      | Search query                   |

**Response:**

```json
{
  "query": "search terms",
  "count": 5,
  "results": [
    {
      "title": "Article Title",
      "url": "https://grokipedia.com/article-path",
      "snippet": "Preview text from the article..."
    }
  ]
}
```

**Response Fields:**

| Field   | Type     | Description                           |
|---------|----------|---------------------------------------|
| query   | string   | The search query that was executed    |
| count   | integer  | Number of results found               |
| results | array    | Array of search result objects        |

**Search Result Object:**

| Field   | Type   | Description                              |
|---------|--------|------------------------------------------|
| title   | string | Article title                            |
| url     | string | Full URL to the article                  |
| snippet | string | Preview/excerpt from the article         |

**Example:**

```bash
curl "http://localhost:8080/api/search?q=machine+learning"
```

**PowerShell:**

```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/search?q=machine+learning" -Method Get
Write-Host "Found $($response.count) results"
$response.results | ForEach-Object {
    Write-Host "- $($_.title)"
}
```

**Python:**

```python
import requests

response = requests.get('http://localhost:8080/api/search', params={'q': 'machine learning'})
data = response.json()
print(f"Found {data['count']} results")
for result in data['results']:
    print(f"- {result['title']}: {result['snippet'][:100]}...")
```

**JavaScript:**

```javascript
const query = 'machine learning';
fetch(`http://localhost:8080/api/search?q=${encodeURIComponent(query)}`)
  .then(response => response.json())
  .then(data => {
    console.log(`Found ${data.count} results`);
    data.results.forEach(result => {
      console.log(`- ${result.title}`);
    });
  });
```

---

## Error Handling

### Common Errors

#### 400 Bad Request

Missing required parameters.

```json
{
  "error": "Bad Request",
  "message": "Article path is required"
}
```

#### 404 Not Found

Article or resource not found.

```json
{
  "error": "Not Found",
  "message": "Article not found"
}
```

#### 500 Internal Server Error

Server-side error (e.g., network issues, parsing errors).

```json
{
  "error": "Internal Server Error",
  "message": "Failed to fetch article: connection timeout"
}
```

---

## Rate Limiting

Currently, there is no rate limiting implemented. However, please be respectful of Grokipedia's servers and avoid making excessive requests.

**Recommendations:**
- Implement caching on your client side
- Add delays between requests
- Use batch operations when possible

---

## CORS

CORS is enabled for all origins (`*`). This allows the API to be called from web browsers.

---

## Best Practices

1. **Cache responses** - Store frequently accessed articles locally
2. **Handle errors gracefully** - Always check response status codes
3. **Use appropriate timeouts** - Network requests may take time
4. **Respect the source** - Don't overload Grokipedia's servers
5. **Validate paths** - Ensure article paths are valid before making requests

---

## Examples

### Complete Workflow Example (Python)

```python
import requests
import json

class GrokipediaAPI:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
    
    def health_check(self):
        response = requests.get(f"{self.base_url}/health")
        return response.json()
    
    def get_article(self, path):
        response = requests.get(f"{self.base_url}/api/article/{path}")
        response.raise_for_status()
        return response.json()
    
    def search(self, query):
        response = requests.get(f"{self.base_url}/api/search", params={"q": query})
        response.raise_for_status()
        return response.json()

# Usage
api = GrokipediaAPI()

# Check health
print(api.health_check())

# Search for articles
results = api.search("artificial intelligence")
print(f"Found {results['count']} results")

# Get first result if available
if results['count'] > 0:
    first_url = results['results'][0]['url']
    # Extract path from URL
    path = first_url.replace('https://grokipedia.com/', '')
    article = api.get_article(path)
    print(f"Title: {article['title']}")
```

---

## Troubleshooting

### Server not responding

1. Check if the server is running: `curl http://localhost:8080/health`
2. Verify the port is not blocked by firewall
3. Check server logs for errors

### Articles not loading

1. Verify the article path is correct
2. Check if the article exists on Grokipedia.com
3. The website structure may have changed - check the scraping logic

### Search returns no results

1. Try different search terms
2. The search functionality depends on Grokipedia's search page structure
3. Check if Grokipedia's search page is accessible

---

## Support

For issues, questions, or contributions, please refer to the main README.md file.

