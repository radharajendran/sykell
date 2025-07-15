package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"sykell-backend/internal/models"
	"sykell-backend/internal/repository"
	"sykell-backend/pkg/logger"
)

type CrawlerService struct {
	repo   *repository.CrawlerRepository
	client *http.Client
	mu     sync.RWMutex
}

func NewCrawlerService() *CrawlerService {
	return &CrawlerService{
		repo: repository.NewCrawlerRepository(),
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Limit redirects to prevent infinite loops
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (s *CrawlerService) AddURL(urlStr string) (*models.CrawlURL, error) {
	// Validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("URL must use http or https scheme")
	}

	return s.repo.CreateCrawlURL(urlStr)
}

func (s *CrawlerService) CrawlURL(id int) error {
	crawlURL, err := s.repo.GetCrawlURLByID(id)
	if err != nil {
		return err
	}

	if crawlURL == nil {
		return fmt.Errorf("crawl URL not found")
	}

	// Update status to running
	crawlURL.Status = models.StatusRunning
	now := time.Now()
	crawlURL.LastCrawledAt = &now

	if err := s.repo.UpdateCrawlURL(crawlURL); err != nil {
		return err
	}

	// Perform the actual crawling
	go s.performCrawl(crawlURL)

	return nil
}

func (s *CrawlerService) performCrawl(crawlURL *models.CrawlURL) {
	logger.Sugar().Infof("Starting crawl for URL: %s", crawlURL.URL)

	defer func() {
		if r := recover(); r != nil {
			logger.Sugar().Errorf("Panic during crawl of %s: %v", crawlURL.URL, r)
			crawlURL.Status = models.StatusError
			crawlURL.ErrorMessage = fmt.Sprintf("Internal error: %v", r)
			s.repo.UpdateCrawlURL(crawlURL)
		}
	}()

	// Fetch the webpage
	resp, err := s.client.Get(crawlURL.URL)
	if err != nil {
		crawlURL.Status = models.StatusError
		crawlURL.ErrorMessage = fmt.Sprintf("Failed to fetch URL: %v", err)
		s.repo.UpdateCrawlURL(crawlURL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		crawlURL.Status = models.StatusError
		crawlURL.ErrorMessage = fmt.Sprintf("HTTP error: %d", resp.StatusCode)
		s.repo.UpdateCrawlURL(crawlURL)
		return
	}

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		crawlURL.Status = models.StatusError
		crawlURL.ErrorMessage = fmt.Sprintf("Failed to parse HTML: %v", err)
		s.repo.UpdateCrawlURL(crawlURL)
		return
	}

	// Extract information from HTML
	s.extractHTMLInfo(crawlURL, doc)

	// Extract and check links
	links := s.extractLinks(doc, crawlURL.URL)
	s.categorizeAndCheckLinks(crawlURL, links)

	crawlURL.Status = models.StatusCompleted
	crawlURL.ErrorMessage = ""

	if err := s.repo.UpdateCrawlURL(crawlURL); err != nil {
		logger.Sugar().Errorf("Failed to update crawl URL: %v", err)
	}

	logger.Sugar().Infof("Completed crawl for URL: %s", crawlURL.URL)
}

func (s *CrawlerService) extractHTMLInfo(crawlURL *models.CrawlURL, doc *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch strings.ToLower(n.Data) {
			case "html":
				// Check for HTML version in doctype or html attributes
				for _, attr := range n.Attr {
					if attr.Key == "version" {
						crawlURL.HTMLVersion = attr.Val
					}
				}
				if crawlURL.HTMLVersion == "" {
					crawlURL.HTMLVersion = "HTML5" // Default assumption
				}
			case "title":
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					crawlURL.Title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "h1":
				crawlURL.H1Count++
			case "h2":
				crawlURL.H2Count++
			case "h3":
				crawlURL.H3Count++
			case "h4":
				crawlURL.H4Count++
			case "h5":
				crawlURL.H5Count++
			case "h6":
				crawlURL.H6Count++
			case "form":
				// Check if it's a login form
				if s.isLoginForm(n) {
					crawlURL.HasLoginForm = true
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func (s *CrawlerService) isLoginForm(form *html.Node) bool {
	hasPasswordField := false
	hasUsernameField := false

	var checkInputs func(*html.Node)
	checkInputs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var inputType, inputName string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "type":
					inputType = strings.ToLower(attr.Val)
				case "name":
					inputName = strings.ToLower(attr.Val)
				}
			}

			if inputType == "password" {
				hasPasswordField = true
			}

			// Common username/email field patterns
			if inputType == "text" || inputType == "email" {
				if strings.Contains(inputName, "user") ||
					strings.Contains(inputName, "email") ||
					strings.Contains(inputName, "login") {
					hasUsernameField = true
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			checkInputs(c)
		}
	}

	checkInputs(form)
	return hasPasswordField && hasUsernameField
}

func (s *CrawlerService) extractLinks(doc *html.Node, baseURL string) []string {
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if attr.Val != "" && !strings.HasPrefix(attr.Val, "#") {
						links = append(links, attr.Val)
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}

func (s *CrawlerService) categorizeAndCheckLinks(crawlURL *models.CrawlURL, links []string) {
	baseURL, err := url.Parse(crawlURL.URL)
	if err != nil {
		logger.Sugar().Errorf("Failed to parse base URL: %v", err)
		return
	}

	internalCount := 0
	externalCount := 0
	inaccessibleCount := 0

	// Create a context with timeout for link checking
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Channel to limit concurrent requests
	semaphore := make(chan struct{}, 10)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, link := range links {
		// Resolve relative URLs
		linkURL, err := baseURL.Parse(link)
		if err != nil {
			continue
		}

		isInternal := linkURL.Host == baseURL.Host || linkURL.Host == ""

		if isInternal {
			internalCount++
		} else {
			externalCount++
		}

		// Check if link is accessible (only for a sample to avoid overwhelming)
		if len(links) <= 50 || (len(links) > 50 && (internalCount+externalCount) <= 50) {
			wg.Add(1)
			go func(linkStr string, linkURL *url.URL) {
				defer wg.Done()

				select {
				case semaphore <- struct{}{}:
					defer func() { <-semaphore }()

					if s.checkLinkAccessibility(ctx, linkURL.String()) != nil {
						mu.Lock()
						inaccessibleCount++

						// Store broken link in database
						brokenLink := &models.BrokenLink{
							CrawlURLID:   crawlURL.ID,
							URL:          linkURL.String(),
							StatusCode:   0, // Will be set by checkLinkAccessibility if available
							ErrorMessage: "Link check failed",
						}
						s.repo.CreateBrokenLink(brokenLink)
						mu.Unlock()
					}
				case <-ctx.Done():
					return
				}
			}(link, linkURL)
		}
	}

	wg.Wait()

	crawlURL.InternalLinksCount = internalCount
	crawlURL.ExternalLinksCount = externalCount
	crawlURL.InaccessibleLinksCount = inaccessibleCount
}

func (s *CrawlerService) checkLinkAccessibility(ctx context.Context, urlStr string) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", urlStr, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

func (s *CrawlerService) GetCrawlURLs(page, limit int, status, search string) ([]models.CrawlURL, int, error) {
	offset := (page - 1) * limit
	return s.repo.GetCrawlURLs(limit, offset, status, search)
}

func (s *CrawlerService) GetCrawlResult(id int) (*models.CrawlResult, error) {
	crawlURL, err := s.repo.GetCrawlURLByID(id)
	if err != nil {
		return nil, err
	}

	if crawlURL == nil {
		return nil, fmt.Errorf("crawl URL not found")
	}

	brokenLinks, err := s.repo.GetBrokenLinks(id)
	if err != nil {
		return nil, err
	}

	return &models.CrawlResult{
		CrawlURL:    *crawlURL,
		BrokenLinks: brokenLinks,
	}, nil
}

func (s *CrawlerService) DeleteCrawlURLs(ids []int) error {
	return s.repo.DeleteCrawlURLs(ids)
}

func (s *CrawlerService) GetStats() (*models.CrawlStats, error) {
	return s.repo.GetCrawlStats()
}

func (s *CrawlerService) ReCrawlURLs(ids []int) error {
	for _, id := range ids {
		crawlURL, err := s.repo.GetCrawlURLByID(id)
		if err != nil {
			logger.Sugar().Errorf("Failed to get crawl URL %d: %v", id, err)
			continue
		}

		if crawlURL != nil {
			crawlURL.Status = models.StatusQueued
			crawlURL.ErrorMessage = ""
			if err := s.repo.UpdateCrawlURL(crawlURL); err != nil {
				logger.Sugar().Errorf("Failed to update crawl URL %d: %v", id, err)
			}
		}
	}

	return nil
}
