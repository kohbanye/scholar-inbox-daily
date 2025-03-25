package notifier

import (
	"fmt"
	"os"
	"time"

	"github.com/kohbanye/scholar-inbox-daily/internal/domain"
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	api       *slack.Client
	channelID string
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
	}, nil
}

func (n *SlackNotifier) PostPapers(papers []domain.Paper) error {
	today := time.Now().Format("2006-01-02")
	header := fmt.Sprintf("Scholar Inbox Daily Papers (%s)", today)

	var blocks []slack.Block
	blocks = append(blocks, slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", header, false, false)))
	blocks = append(blocks, slack.NewDividerBlock())

	for i, paper := range papers {
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

	_, _, err := n.api.PostMessage(
		n.channelID,
		slack.MsgOptionBlocks(blocks...),
	)
	return err
}
