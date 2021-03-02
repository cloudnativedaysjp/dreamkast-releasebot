package router

import (
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/controller/slackrtm"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/slackrtm_receiver"
	"github.com/sirupsen/logrus"
)

type Router struct {
	Logger *logrus.Logger

	ErrStream chan error
}

func New(l *logrus.Logger) *Router {
	return &Router{
		Logger: l,
	}
}

func (r *Router) Run(token string, targetRepositories []string, gitapi infrastructure.GitAPIDao, cache infrastructure.CacheDao) {
	/* init SlackRTM */
	m, err := slackrtm_receiver.NewReceiver(r.Logger, token)
	if err != nil {
		r.ErrStream <- err
		return
	}

	controller := slackrtm.NewMessageController(r.Logger, gitapi, cache, m.BotUserID, targetRepositories)

	m.DefaultHandleFunc(controller.Default)
	m.HandleFunc("help", controller.CommandHelp)
	m.HandleFunc("ping", controller.CommandPing)
	m.HandleFunc("release", controller.CommandRelease)

	if err := m.Serve(); err != nil {
		r.ErrStream <- err
	}
}
