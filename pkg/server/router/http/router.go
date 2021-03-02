package http

import (
	"net/http"

	http_controller "github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/controller/http"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/global"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/slack_interactivemsg"
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

// Run is entrypoint
func (r *Router) Run(verificationToken, bindAddr, listenPath string, gitapi infrastructure.GitAPIDao, gitcommand infrastructure.GitCommandDao, cache infrastructure.CacheDao, slackmsg infrastructure.SlackMsgDao, baseBranch string, enableAutoMerge bool) {
	/* Controllers */
	controller := http_controller.NewInteractionMsgController(r.Logger, gitapi, gitcommand, cache, slackmsg, baseBranch, enableAutoMerge)

	/* register Slack Interactive Message handlers */
	interactiveMsgHandlers := slack_interactivemsg.New(r.Logger, verificationToken)
	interactiveMsgHandlers.HandleFunc(global.ActionNameCancel, controller.SelectedCancel)
	interactiveMsgHandlers.HandleFunc(global.ActionNameRelease, controller.SelectedRepository)
	interactiveMsgHandlers.HandleFunc(global.ActionNameReleaseVersionMajor, controller.SelectedReleaseLevel)
	interactiveMsgHandlers.HandleFunc(global.ActionNameReleaseVersionMinor, controller.SelectedReleaseLevel)
	interactiveMsgHandlers.HandleFunc(global.ActionNameReleaseVersionPatch, controller.SelectedReleaseLevel)
	interactiveMsgHandlers.HandleFunc(global.ActionNameReleaseConfirm, controller.SelectedConfirmRelease)

	/* run HTTP server */
	m := http.NewServeMux()
	m.HandleFunc(listenPath, interactiveMsgHandlers.Handler)
	m.HandleFunc(`/health`, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	// listen
	if err := http.ListenAndServe(bindAddr, m); err != nil {
		r.ErrStream <- err
	}
}
