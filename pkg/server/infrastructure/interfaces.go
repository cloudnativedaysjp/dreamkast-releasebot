package infrastructure

import (
	"context"

	"github.com/slack-go/slack"
)

type GitAPIDao interface {
	HealthCheck() error
	CreatePullRequest(ctx context.Context, org, repo, headBranch, baseBranch, title, body string) (prNum int, err error)
	LabelPullRequest(ctx context.Context, org, repo string, prNum int, label string) (err error)
	MergePullRequest(ctx context.Context, org, repo string, prNum int) error
	DeleteBranch(ctx context.Context, org, repo, headBranch string) error
}

type GitCommandDao interface {
	HealthCheck() error
	Clone(ctx context.Context, org, repo string) (dirPath string, err error)
	SwitchNewBranch(ctx context.Context, dirPath, branch string) error
	CommitAll(ctx context.Context, dirPath, commitMsg string) error
	Push(ctx context.Context, dirPath string) error
	Remove(ctx context.Context, dir string) (err error)
}

type SlackMsgDao interface {
	HealthCheck() error
	SendMessage(ctx context.Context, msg slack.Message, channel string) error
	SendThreadMessage(ctx context.Context, msg slack.Message, channel, ts string) error
	UpdateMessage(ctx context.Context, msg slack.Message, channel, ts string) error
}

type CacheDao interface {
	HealthCheck() error
	Read(key string) (string, error)
	Write(key, value string) error
	Remove(key string) error
}
