package slackrtm

import (
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/view"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// GitHubWebhookController is Controller
type MessageController struct {
	logger *logrus.Logger

	gitApi             infrastructure.GitAPIDao
	cache              infrastructure.CacheDao
	botUserID          string
	targetRepositories []string
}

func NewMessageController(logger *logrus.Logger, gitapi infrastructure.GitAPIDao, cache infrastructure.CacheDao, botUserID string, targetRepositories []string) *MessageController {
	return &MessageController{logger, gitapi, cache, botUserID, targetRepositories}
}

func (mc *MessageController) Default(ev *slack.MessageEvent) (slack.Message, error) {
	return view.CommandDefault(mc.botUserID), nil
}

func (mc *MessageController) CommandHelp(ev *slack.MessageEvent) (slack.Message, error) {
	return view.CommandHelp(mc.botUserID), nil
}

func (mc *MessageController) CommandPing(ev *slack.MessageEvent) (slack.Message, error) {
	if err := mc.gitApi.HealthCheck(); err != nil {
		return slack.Message{}, err
	}
	return view.CommandPing(), nil
}

func (mc *MessageController) CommandRelease(ev *slack.MessageEvent) (slack.Message, error) {
	// generate unique id
	callbackId := fmt.Sprintf("%v", time.Now().UnixNano())

	return view.CommandRelease(callbackId, mc.targetRepositories), nil
}
