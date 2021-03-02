package slackrtm_receiver

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Receiver struct {
	Logger *logrus.Logger

	Client    *slack.Client
	RTM       *slack.RTM
	BotUserID string // starts with 'U', used Slack mention message from users

	Commands       []Command
	DefaultCommand Command
}

type Command struct {
	name    string
	aliases []string
	h       Handler
}

type Handler func(*slack.MessageEvent) (slack.Message, error)

func NewReceiver(l *logrus.Logger, token string) (*Receiver, error) {
	c := slack.New(token)
	res, err := c.AuthTest()
	if err != nil {
		return nil, err
	}
	rtm := c.NewRTM()
	go rtm.ManageConnection() // Start listening slack events

	return &Receiver{
		Logger:    l,
		Client:    c,
		RTM:       rtm,
		BotUserID: res.UserID,
	}, nil
}

func (r *Receiver) HandleFunc(commandName string, handler Handler) {
	// return if already exists
	for _, command := range r.Commands {
		if commandName == command.name {
			return
		}
	}
	// register
	r.Commands = append(r.Commands, Command{name: commandName, h: handler})
}

func (r *Receiver) DefaultHandleFunc(handler Handler) {
	r.DefaultCommand = Command{h: handler}
}

func (r *Receiver) Aliases(commandName string, aliases ...string) {
	for _, command := range r.Commands {
		if commandName == command.name {
			command.aliases = append(command.aliases, aliases...)
		}
	}
}

func (r *Receiver) Serve() error {
Receive:
	for msg := range r.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			input := strings.Fields(ev.Msg.Text)
			if len(input) == 0 || input[0] != fmt.Sprintf("<@%s>", r.BotUserID) {
				continue Receive
			}
			for _, command := range r.Commands {
				if input[1] == command.name || contains(command.aliases, input[1]) {
					if err := r.callHandlerAndPostMessage(command, ev); err != nil {
						r.Logger.Error(err)
					}
					continue Receive
				}
			}
			// if not match command of handler, call default handler
			if err := r.callHandlerAndPostMessage(r.DefaultCommand, ev); err != nil {
				r.Logger.Error(err)
			}
		case *slack.UnmarshallingErrorEvent:
			r.Logger.Warn(ev)
		}
	}
	return nil
}

func (r *Receiver) callHandlerAndPostMessage(command Command, ev *slack.MessageEvent) error {
	result, err := command.h(ev)
	if err != nil {
		_, _, err := r.Client.PostMessage(ev.Channel,
			slack.MsgOptionText(err.Error(), false),
		)
		return err
	}

	var messages []slack.MsgOption
	if result.Text != "" {
		messages = append(messages, slack.MsgOptionText(result.Text, false))
	}
	if len(result.Attachments) != 0 {
		messages = append(messages, slack.MsgOptionAttachments(result.Attachments...))
	}

	if len(messages) != 0 {
		if _, _, err = r.Client.PostMessage(ev.Channel, messages...); err != nil {
			return err
		}
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
