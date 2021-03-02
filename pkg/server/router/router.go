package router

import (
	"fmt"

	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure/gitcommand"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure/githubapi"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure/localcache"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure/slackmsg"
	http_router "github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/router/http"
	slackrtm_router "github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/router/slackrtm"
	"github.com/sirupsen/logrus"
)

// Run is entrypoint
func Run(config Config) error {
	/* Logger */
	logger := logrus.New()

	/* Drivers */
	githubApiDriver := githubapi.NewGitHubApiDriver(logger, config.GitHubToken)
	gitCommandDriver := gitcommand.NewGitCommandDriver(logger, config.AppName, "dummy@example.com", config.GitHubToken)
	cacheDriver := localcache.New()
	slackDriver, err := slackmsg.NewSlackMsg(logger, config.BotToken)
	if err != nil {
		return err
	}

	/* Route */
	{
		/* Routing request from WebSocket (RTM receiver) */
		sr := slackrtm_router.New(logger)
		go sr.Run(
			config.BotToken,
			config.TargetRepositories,
			githubApiDriver,
			cacheDriver,
		)

		/* Routing request from http (interactive message) */
		hr := http_router.New(logger)
		go hr.Run(
			config.SlackVerificationToken,
			config.HTTPBindAddr,
			config.SlackInteractionListenPath,
			githubApiDriver,
			gitCommandDriver,
			cacheDriver,
			slackDriver,
			config.BaseBranch,
			config.EnableAutoMerge,
		)

		var err error
		for {
			select {
			case err = <-sr.ErrStream:
				return fmt.Errorf(`Error: slackrtm Router: %v`, err)
			case err = <-hr.ErrStream:
				return fmt.Errorf(`Error: HTTP Router: %v`, err)
			}
		}
	}
}
