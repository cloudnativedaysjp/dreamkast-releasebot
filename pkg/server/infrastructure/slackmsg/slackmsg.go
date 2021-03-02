package slackmsg

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackMsg struct {
	logger    *logrus.Logger
	client    *slack.Client
	botUserId string
}

func NewSlackMsg(l *logrus.Logger, token string) (*SlackMsg, error) {
	c := slack.New(token)
	res, err := c.AuthTest()
	if err != nil {
		return nil, err
	}

	return &SlackMsg{l, c, res.UserID}, nil
}

func (s *SlackMsg) HealthCheck() error {
	_, err := s.client.AuthTest()
	return err
}

func (s *SlackMsg) SendMessage(ctx context.Context, msg slack.Message, channel string) error {
	_, _, err := s.client.PostMessageContext(ctx, channel,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SlackMsg) SendThreadMessage(ctx context.Context, msg slack.Message, channel, ts string) error {
	_, _, err := s.client.PostMessageContext(
		ctx, channel,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
		slack.MsgOptionTS(ts),
		slack.MsgOptionBroadcast(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SlackMsg) UpdateMessage(ctx context.Context, msg slack.Message, channel, ts string) error {
	_, _, _, err := s.client.UpdateMessageContext(
		ctx, channel, ts,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
	)
	if err != nil {
		return err
	}
	return nil
}
