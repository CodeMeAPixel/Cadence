package web

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PageContent struct {
	URL         string
	Title       string
	Description string
	Body        string
	AllText     string
	Headers     map[string][]string
	StatusCode  int
	FetchedAt   time.Time
}

type Fetcher struct {
	client  *http.Client
	timeout time.Duration
}

func NewFetcher(timeout time.Duration) *Fetcher {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (f *Fetcher) Fetch(url string) (*PageContent, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	content := &PageContent{
		URL:        url,
		StatusCode: resp.StatusCode,
		FetchedAt:  time.Now(),
		Headers:    resp.Header,
	}

	content.Title = strings.TrimSpace(doc.Find("title").First().Text())

	metaDesc, _ := doc.Find("meta[name='description']").Attr("content")
	content.Description = metaDesc

	doc.Find("script, style").Remove()

	allText := doc.Find("body").Text()
	content.Body = strings.TrimSpace(allText)
	content.AllText = extractStructuredText(doc)

	return content, nil
}

func extractStructuredText(doc *goquery.Document) string {
	var texts []string

	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			texts = append(texts, text)
		}
	})

	doc.Find("p").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if len(text) > 20 {
			texts = append(texts, text)
		}
	})

	doc.Find("li").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			texts = append(texts, text)
		}
	})

	return strings.Join(texts, "\n")
}

func (pc *PageContent) GetMainContent() string {
	if pc.AllText != "" {
		return pc.AllText
	}
	return pc.Body
}
