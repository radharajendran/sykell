package service

import (
	"sync"
	"time"

	"sykell-backend/internal/models"
	"sykell-backend/internal/repository"
	"sykell-backend/pkg/logger"
)

type CrawlerJobProcessor struct {
	crawlerService *CrawlerService
	repo           *repository.CrawlerRepository
	stopChan       chan bool
	wg             sync.WaitGroup
	isRunning      bool
	mu             sync.RWMutex
}

func NewCrawlerJobProcessor() *CrawlerJobProcessor {
	return &CrawlerJobProcessor{
		crawlerService: NewCrawlerService(),
		repo:           repository.NewCrawlerRepository(),
		stopChan:       make(chan bool),
	}
}

func (p *CrawlerJobProcessor) Start() {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return
	}
	p.isRunning = true
	p.mu.Unlock()

	logger.Sugar().Info("Starting crawler job processor")

	p.wg.Add(1)
	go p.processJobs()
}

func (p *CrawlerJobProcessor) Stop() {
	p.mu.Lock()
	if !p.isRunning {
		p.mu.Unlock()
		return
	}
	p.isRunning = false
	p.mu.Unlock()

	logger.Sugar().Info("Stopping crawler job processor")
	close(p.stopChan)
	p.wg.Wait()
}

func (p *CrawlerJobProcessor) processJobs() {
	defer p.wg.Done()

	ticker := time.NewTicker(10 * time.Second) // Check for queued jobs every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-p.stopChan:
			logger.Sugar().Info("Crawler job processor stopped")
			return
		case <-ticker.C:
			p.processQueuedJobs()
		}
	}
}

func (p *CrawlerJobProcessor) processQueuedJobs() {
	// Get queued crawl URLs
	queuedURLs, _, err := p.repo.GetCrawlURLs(10, 0, models.StatusQueued, "")
	if err != nil {
		logger.Sugar().Errorf("Failed to get queued URLs: %v", err)
		return
	}

	if len(queuedURLs) == 0 {
		return
	}

	logger.Sugar().Infof("Processing %d queued crawl jobs", len(queuedURLs))

	for _, crawlURL := range queuedURLs {
		// Start crawling in background
		go func(url models.CrawlURL) {
			// Update status to running
			url.Status = models.StatusRunning
			now := time.Now()
			url.LastCrawledAt = &now

			if err := p.repo.UpdateCrawlURL(&url); err != nil {
				logger.Sugar().Errorf("Failed to update crawl URL status: %v", err)
				return
			}

			// Perform crawl
			p.crawlerService.performCrawl(&url)
		}(crawlURL)
	}
}
