package notifier

import (
	"fmt"
	"os"
	"time"

	"github.com/kohbanye/scholar-inbox-daily/internal/domain"
	"github.com/kohbanye/scholar-inbox-daily/internal/logger"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	api       *slack.Client
	channelID string
	logger    *logger.Logger
}

func NewSlackNotifier() (*SlackNotifier, error) {
	token := os.Getenv("SLACK_API_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	if token == "" || channelID == "" {
		return nil, fmt.Errorf("SLACK_API_TOKEN and SLACK_CHANNEL_ID environment variables must be set")
	}

	api := slack.New(token)
	return &SlackNotifier{
		api:       api,
		channelID: channelID,
		logger:    logger.New(),
	}, nil
}

func (n *SlackNotifier) PostPapers(papers []domain.Paper) error {
	n.logger.Info("Starting to post papers to Slack")
	today := time.Now().Format("2006-01-02")
	header := fmt.Sprintf("Scholar Inbox Daily Papers (%s)", today)

	var blocks []slack.Block
	blocks = append(blocks, slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", header, false, false)))
	blocks = append(blocks, slack.NewDividerBlock())

	for i, paper := range papers {
		n.logger.Info(fmt.Sprintf("Posting paper %d: %+v", i, paper))
		titleText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>*", paper.URL, paper.Title), false, false)
		metaText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Authors: %s\nDate: %s", paper.Authors, paper.Date), false, false)

		abstractText := paper.Abstract
		if len(abstractText) > 200 {
			abstractText = abstractText[:200] + "..."
		}
		abstractBlock := slack.NewTextBlockObject("mrkdwn", abstractText, false, false)

		blocks = append(blocks, slack.NewSectionBlock(titleText, nil, nil))
		blocks = append(blocks, slack.NewSectionBlock(metaText, nil, nil))
		blocks = append(blocks, slack.NewSectionBlock(abstractBlock, nil, nil))

		if i < len(papers)-1 {
			blocks = append(blocks, slack.NewDividerBlock())
		}
	}

	n.logger.Info(fmt.Sprintf("Posting %d papers to Slack", len(papers)))
	_, _, err := n.api.PostMessage(
		n.channelID,
		slack.MsgOptionBlocks(blocks...),
	)
	if err != nil {
		n.logger.Error("Failed to post message to Slack", err)
		return fmt.Errorf("error posting to Slack: %w", err)
	}

	n.logger.Info("Successfully posted papers to Slack")
	return nil
}
