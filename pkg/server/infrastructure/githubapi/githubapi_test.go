package githubapi

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_GitHubApiDriver(t *testing.T) {
	driver := NewGitHubApiDriver(logrus.New(), os.Getenv("GITHUB_TOKEN"))

	t.Run(`HealthCheck`, func(t *testing.T) {
		err := driver.HealthCheck()
		if err != nil {
			t.Fatalf("error: %s", err)
		}
	})

	/* 動作確認済み (実行の度に PR の作成とマージが走るためコメントアウト)
	t.Run(`CreatePullRequest & MergePullRequest`, func(t *testing.T) {
		ctx := context.Background()
		org := "ShotaKitazawa"
		repo := "dotfiles"
		headBranch := "demo"
		baseBranch := "master"
		label := "bug"

		// CreatePullRequest
		prNum, err := driver.CreatePullRequest(ctx, org, repo, headBranch, baseBranch, "demo", "hoge\n`fuga`\n**piyo**")
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		// LabelPullRequest
		if err := driver.LabelPullRequest(ctx, org, repo, prNum, label); err != nil {
			t.Fatalf("error: %s", err)
		}

		// MergePullRequest
		if err := driver.MergePullRequest(ctx, org, repo, prNum); err != nil {
			t.Fatalf("error: %s", err)
		}
	})
	*/

}
