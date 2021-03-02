package slack_interactivemsg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackInteractiveMsg struct {
	Logger            *logrus.Logger
	VerificationToken string
	Actions           []Action
}

type Action struct {
	name string
	h    Handler
}
type Handler func(w http.ResponseWriter, r *http.Request)

func New(l *logrus.Logger, verificationToken string) SlackInteractiveMsg {
	return SlackInteractiveMsg{Logger: l, VerificationToken: verificationToken}
}

func (s *SlackInteractiveMsg) HandleFunc(actionName string, handler Handler) {
	s.Actions = append(s.Actions, Action{actionName, handler})
}

func (s SlackInteractiveMsg) Handler(w http.ResponseWriter, r *http.Request) {
	/* Check Method */
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	/* Decode Body */
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to read request body: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonStr, err := url.QueryUnescape(strings.TrimPrefix(string(buf), "payload="))
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to unespace request body: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to decode json message from slack: %s", jsonStr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/* Verify Token */
	if message.Token != s.VerificationToken {
		s.Logger.Debug("invalid varification token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	originalMessage := message.OriginalMessage
	callbackAction := *(message.ActionCallback.AttachmentActions[0])

	for _, action := range s.Actions {
		if callbackAction.Name == action.name {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "originalMessage", originalMessage)
			ctx = context.WithValue(ctx, "channelID", message.Channel.Conversation.ID)
			ctx = context.WithValue(ctx, "callbackAction", callbackAction)
			ctx = context.WithValue(ctx, "callbackID", message.CallbackID)
			action.h(w, r.WithContext(ctx))
			break
		}
	}
}
