package main

import (
	"log"
	"os"

	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/router"
	"github.com/urfave/cli/v2"
)

var (
	appName    = "dreamkast-chatbot"
	appVersion = ""
)

func main() {
	app := cli.NewApp()
	app = &cli.App{
		Name:                 appName,
		Version:              appVersion,
		Usage:                "",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   `Log Level (only support "debug", "info", "warn", "error", or "critical")`,
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:     "slack-bot-token",
				Usage:    "Slack Bot Token",
				Required: true,
				EnvVars:  []string{"SLACK_BOT_TOKEN"},
			},
			&cli.StringFlag{
				Name:     "github-token",
				Usage:    "GitHub Token for Bot",
				Required: true,
				EnvVars:  []string{"GITHUB_TOKEN"},
			},
			&cli.StringFlag{
				Name:     "slack-verification-token",
				Usage:    "Verification token for Slack interactive message",
				Required: true,
				EnvVars:  []string{"SLACK_VERIFICATION_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "http-bind-addr",
				Usage:   "http bind address (format is ADDR:PORT)",
				Value:   "0.0.0.0:8080",
				EnvVars: []string{"HTTP_BIND_ADDR"},
			},
			&cli.StringFlag{
				Name:    "slack-interaction-listen-path",
				Usage:   "path of listen for Slack interactive message",
				Value:   "/interaction",
				EnvVars: []string{"SLACK_INTERACTION_LISTEN_PATH"},
			},
			&cli.StringSliceFlag{
				Name:     "target-repositories",
				Usage:    `specify target repositories`,
				Required: true,
				EnvVars:  []string{"TARGET_REPOSITORIES"},
			},
			&cli.StringFlag{
				Name:     "base-branch",
				Usage:    `specify base branch for each repositories`,
				Required: true,
				EnvVars:  []string{"BASE_BRANCH"},
			},
			&cli.BoolFlag{
				Name:    "enable-auto-merge",
				Usage:   `whether auto merge in release PR created by Bot`,
				Value:   false,
				EnvVars: []string{"ENABLE_AUTO_MERGE"},
			},
		},
		Action: func(c *cli.Context) error {
			config := router.Config{
				AppName:                    appName,
				AppVersion:                 appVersion,
				LogLevel:                   c.String("log-level"),
				BotToken:                   c.String("slack-bot-token"),
				GitHubToken:                c.String("github-token"),
				HTTPBindAddr:               c.String("http-bind-addr"),
				SlackInteractionListenPath: c.String("slack-interaction-listen-path"),
				SlackVerificationToken:     c.String("slack-verification-token"),
				TargetRepositories:         c.StringSlice("target-repositories"),
				BaseBranch:                 c.String("base-branch"),
				EnableAutoMerge:            c.Bool("enable-auto-merge"),
			}
			return router.Run(config)
		},
	}
	// Execute
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
