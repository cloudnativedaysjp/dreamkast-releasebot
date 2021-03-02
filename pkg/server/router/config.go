package router

type Config struct {
	AppName                    string
	AppVersion                 string
	LogLevel                   string
	BotToken                   string
	GitHubToken                string
	HTTPBindAddr               string
	SlackInteractionListenPath string
	SlackVerificationToken     string
	TargetRepositories         []string
	BaseBranch                 string
	EnableAutoMerge            bool
}
