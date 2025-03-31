package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kohbanye/scholar-inbox-daily/internal/logger"
	"github.com/kohbanye/scholar-inbox-daily/internal/notifier"
	"github.com/kohbanye/scholar-inbox-daily/internal/scraper"
)

const (
	MaxPapers = 10
)

func handleRequest(ctx context.Context) error {
	log := logger.New()
	log.Info("Starting Scholar Inbox daily job")

	scholarScraper := scraper.NewScholarScraper()

	papers, err := scholarScraper.FetchPapers()
	if err != nil {
		log.Error("error fetching papers", err)
		return fmt.Errorf("error fetching papers: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"count": len(papers),
	}).Info("Found papers")

	if len(papers) > MaxPapers {
		log.WithFields(map[string]any{
			"limit": MaxPapers,
		}).Warn("Limiting papers")
		papers = papers[:MaxPapers]
	}

	slackNotifier, err := notifier.NewSlackNotifier()
	if err != nil {
		log.Error("error creating Slack notifier", err)
		return fmt.Errorf("error creating Slack notifier: %w", err)
	}

	if err := slackNotifier.PostPapers(papers); err != nil {
		log.Error("error posting to Slack", err)
		return fmt.Errorf("error posting to Slack: %w", err)
	}

	log.Info("Successfully posted papers to Slack")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
