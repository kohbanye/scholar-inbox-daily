package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kohbanye/scholar-inbox-daily/internal/notifier"
	"github.com/kohbanye/scholar-inbox-daily/internal/scraper"
)

func handleRequest(ctx context.Context) error {
	log.Println("Starting Scholar Inbox daily job")

	scholarScraper := scraper.NewScholarScraper()

	papers, err := scholarScraper.FetchPapers()
	if err != nil {
		return fmt.Errorf("error fetching papers: %w", err)
	}

	log.Printf("Found %d papers", len(papers))

	slackNotifier, err := notifier.NewSlackNotifier()
	if err != nil {
		return fmt.Errorf("error creating Slack notifier: %w", err)
	}

	if err := slackNotifier.PostPapers(papers); err != nil {
		return fmt.Errorf("error posting to Slack: %w", err)
	}

	log.Println("Successfully posted papers to Slack")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
