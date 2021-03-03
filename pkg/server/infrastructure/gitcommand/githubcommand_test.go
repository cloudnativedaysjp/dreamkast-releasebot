package gitcommand

import (
	"testing"
)

func Test_GitCommandDriver(t *testing.T) {
	/* 動作確認済み (実行の度に PR の作成とマージが走るためコメントアウト)
	driver := NewGitCommandDriver(logrus.New(), "test", "dummy@example.com", os.Getenv("GITHUB_TOKEN"))

	t.Run(`Clone -> SwitchNewBranch -> CommitAll -> Push`, func(t *testing.T) {
		// Clone
		dir, err := driver.Clone(context.Background(), "ShotaKitazawa", "dotfiles")
		if err != nil {
			t.Fatalf("error: %s", err)
		}
		fmt.Println(dir)

		// SwitchNewBranch
		err = driver.SwitchNewBranch(context.Background(), dir, "demo")
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		// create new empty file
		fp, _ := os.Create("/tmp/dotfiles/.test")
		fp.Close()

		// CommitAll
		err = driver.CommitAll(context.Background(), dir, "for test")
		if err != nil {
			t.Fatalf("error: %s", err)
		}

		// Push
		err = driver.Push(context.Background(), dir)
		if err != nil {
			t.Fatalf("error: %s", err)
		}
	})
	*/

}
