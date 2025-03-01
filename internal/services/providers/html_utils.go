package providers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shashank-sharma/backend/internal/logger"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	MaxBodySize = 5 * 1024 * 1024 
	MaxContentLength = 100 * 1024
)

type HTMLContentFetcher struct {
	client *http.Client
}

func NewHTMLContentFetcher(timeout time.Duration) *HTMLContentFetcher {
	return &HTMLContentFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchContent fetches and extracts the main content from a URL
func (f *HTMLContentFetcher) FetchContent(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", nil
	}

	// Create a request with the given context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Anonymous/1.0)")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, MaxBodySize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Extract content
	content, err := ExtractMainContent(body)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error extracting content from %s: %v", url, err))
		// Return partial content if possible
		if len(content) > 0 {
			return content, nil
		}
		return "", err
	}

	return content, nil
}

// ExtractMainContent extracts the main content from HTML
func ExtractMainContent(body []byte) (string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %v", err)
	}

	// Get the title
	title := extractTitle(doc)

	// Extract main content using various heuristics
	content := extractContent(doc)

	// Combine the title and content 
	if title != "" && content != "" {
		return fmt.Sprintf("%s\n\n%s", title, content), nil
	} else if content != "" {
		return content, nil
	} else if title != "" {
		return title, nil
	}

	return "", fmt.Errorf("no content extracted")
}

func extractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.DataAtom == atom.Title {
		return extractTextFromNode(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := extractTitle(c); title != "" {
			return title
		}
	}

	return ""
}

func extractContent(n *html.Node) string {
	var content string
	
	// Look for article, main, or content containers 
	articleContent := extractTextFromTag(n, atom.Article)
	mainContent := extractTextFromTag(n, atom.Main)
	
	// Also look for common content divs by ID or class 
	contentByID := extractTextFromElementWithAttr(n, "id", []string{"content", "main-content", "article-content", "post-content"})
	contentByClass := extractTextFromElementWithAttr(n, "class", []string{"content", "main-content", "article-content", "post-content", "entry-content"})
	
	// Use the longest content found
	candidates := []string{articleContent, mainContent, contentByID, contentByClass}
	for _, c := range candidates {
		if len(c) > len(content) {
			content = c
		}
	}
	
	// If we couldn't find structured content, just extract paragraphs
	if content == "" {
		content = extractTextFromTag(n, atom.P)
	}
	
	// Limit size of extracted content
	if len(content) > MaxContentLength {
		content = content[:MaxContentLength] + "..."
	}
	
	return strings.TrimSpace(content)
}

func extractTextFromTag(n *html.Node, tag atom.Atom) string {
	if n.Type == html.ElementNode && n.DataAtom == tag {
		return extractTextFromNode(n)
	}
	
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractTextFromTag(c, tag)
	}
	
	return result
}

func extractTextFromElementWithAttr(n *html.Node, attrName string, attrValues []string) string {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrName {
				for _, value := range attrValues {
					if strings.Contains(attr.Val, value) {
						return extractTextFromNode(n)
					}
				}
			}
		}
	}
	
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractTextFromElementWithAttr(c, attrName, attrValues)
	}
	
	return result
}

func extractTextFromNode(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractTextFromNode(c)
	}
	
	if n.Type == html.ElementNode {
		switch n.DataAtom {
		case atom.P, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Br, atom.Li:
			result += "\n"
		}
	}
	
	return result
} 