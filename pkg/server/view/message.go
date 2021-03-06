package view

import (
	"fmt"
	"path/filepath"

	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/global"
	"github.com/slack-go/slack"
)

// [input from chat] Commands & Options
const (
	RELEASE = "release"
	PING    = "ping"
	HELP    = "help"
)

// [output to chat] messages
var (
	DefaultMessage = fmt.Sprintf(`
コマンドが存在しません。 以下のコマンドよりヘルプメッセージを確認してください。
> %[1]s%[2]s%[1]s
`, "`", "<@%[1]s> "+HELP)

	InvalidCommandSlackMessage = fmt.Sprintf(`
コマンドの引数が誤っています。以下のコマンドよりヘルプメッセージを確認してください。
> %[1]s%[2]s%[1]s
`, "`", "<@%[1]s> "+HELP)

	HelpMessage = fmt.Sprintf(`
%[1]s
%[2]s - ヘルプ出力
%[3]s - 疎通確認
%[4]s - 指定リポジトリのバージョンアップ
%[1]s
`, "```",
		"<@%[1]s> "+HELP,
		"<@%[1]s> "+PING,
		"<@%[1]s> "+RELEASE,
	)
)

/* views of Slack Command */

func CommandHelp(botUserId string) (msg slack.Message) {
	msg.Text = fmt.Sprintf(HelpMessage, botUserId)
	return
}

func CommandDefault(botUserId string) (msg slack.Message) {
	msg.Text = fmt.Sprintf(DefaultMessage, botUserId)
	return
}

func CommandInvalid(botUserId string) (msg slack.Message) {
	msg.Text = fmt.Sprintf(InvalidCommandSlackMessage, botUserId)
	return
}

func CommandPing() (msg slack.Message) {
	msg.Text = "pong"
	return
}

func CommandRelease(callbackId string, repoUrls []string) (msg slack.Message) {
	var options []slack.AttachmentActionOption
	for _, repoUrl := range repoUrls {
		repo := filepath.Base(repoUrl)
		org := filepath.Base(filepath.Dir(repoUrl))
		options = append(options,
			slack.AttachmentActionOption{
				Text:  fmt.Sprintf("%s/%s", org, repo),
				Value: fmt.Sprintf("%s__%s", org, repo),
			},
		)
	}

	msg.Attachments = append(msg.Attachments, slack.Attachment{
		Text:       "リリース対象のリポジトリを選択",
		Color:      global.ColorcodeLightGray,
		CallbackID: callbackId,
		Actions: []slack.AttachmentAction{
			{
				Name:    global.ActionNameRelease,
				Text:    "選択",
				Type:    "select",
				Options: options,
			},
			{
				Name:  global.ActionNameCancel,
				Text:  "cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	})
	return
}

/* views of Slack Interactive Message */

func IntMsgCancel(originalMessage slack.Message) (msg slack.Message) {
	respMessage := originalMessage
	respMessage.ResponseType = `in_channel`
	respMessage.ReplaceOriginal = true
	respMessage.Attachments = []slack.Attachment{
		{
			Text:  `キャンセルされました`,
			Color: global.ColorcodeCrimson,
		},
	}
	return respMessage
}

func IntMsgSelectLevel(callbackId string, originalMessage slack.Message, selectedOrg, selectedRepo string) slack.Message {
	msg := originalMessage
	msg.ResponseType = `in_channel`
	msg.ReplaceOriginal = true

	msg.Attachments = []slack.Attachment{{
		Text:       "更新レベルを選択",
		Color:      global.ColorcodeLightGray,
		CallbackID: callbackId,
		Actions: []slack.AttachmentAction{
			{
				Name: global.ActionNameReleaseVersionMajor,
				Text: global.ActionNameReleaseVersionMajor,
				Type: "button",
			},
			{
				Name: global.ActionNameReleaseVersionMinor,
				Text: global.ActionNameReleaseVersionMinor,
				Type: "button",
			},
			{
				Name: global.ActionNameReleaseVersionPatch,
				Text: global.ActionNameReleaseVersionPatch,
				Type: "button",
			},
			{
				Name:  global.ActionNameCancel,
				Text:  "cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}}
	return msg
}

func IntMsgConfirmRelease(callbackId string, originalMessage slack.Message, selectedOrg, selectedRepo, selectedLevel string) slack.Message {
	msg := originalMessage
	msg.ResponseType = `in_channel`
	msg.ReplaceOriginal = true

	msg.Attachments = []slack.Attachment{{
		Text:       fmt.Sprintf("OK? > Target: *%s/%s*, Update Level: *%s*", selectedOrg, selectedRepo, selectedLevel),
		Color:      global.ColorcodeLightGray,
		CallbackID: callbackId,
		Actions: []slack.AttachmentAction{
			{
				Name: global.ActionNameReleaseConfirm,
				Text: "OK",
				Type: "button",
			},
			{
				Name:  global.ActionNameCancel,
				Text:  "cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}}
	return msg
}

func IntMsgReleaseProcessing(originalMessage slack.Message, org, repo, level string) slack.Message {
	msg := originalMessage
	msg.ResponseType = `in_channel`
	msg.ReplaceOriginal = true
	msg.Text = fmt.Sprintf("Target: *%s/%s*, Update Level: *%s*", org, repo, level)
	msg.Attachments = []slack.Attachment{
		{Color: global.ColorcodeLightGray, Text: "processing..."},
	}
	return msg
}

func IntMsgReleaseDone(originalMessage slack.Message, org, repo, level string) slack.Message {
	msg := originalMessage
	msg.ResponseType = `in_channel`
	msg.ReplaceOriginal = true
	msg.Text = fmt.Sprintf("Target: *%s/%s*, Update Level: *%s*", org, repo, level)
	msg.Attachments = []slack.Attachment{
		{Color: global.ColorcodeDeepSkyBlue, Text: "Done!!"},
		{Color: global.ColorcodeDeepSkyBlue, Text: "Tags: <https://github.com/ShotaKitazawa/dotfiles/tags>"},
	}
	return msg
}

func IntMsgReleasePRLink(originalMessage slack.Message, org, repo, level string, prNum int) slack.Message {
	msg := originalMessage
	msg.ResponseType = `in_channel`
	msg.ReplaceOriginal = true
	msg.Text = fmt.Sprintf("Target: *%s/%s*, Update Level: *%s*", org, repo, level)
	msg.Attachments = []slack.Attachment{
		{Color: global.ColorcodeDeepSkyBlue, Text: fmt.Sprintf("Pull Request: <https://github.com/%s/%s/pull/%v>", org, repo, prNum)},
	}
	return msg
}

func IntMsgReleaseFailed(originalMessage slack.Message, org, repo, level string) slack.Message {
	msg := originalMessage
	msg.ResponseType = `in_channel`
	msg.ReplaceOriginal = true
	msg.Text = fmt.Sprintf("Target: *%s/%s*, Update Level: *%s*", org, repo, level)
	msg.Attachments = []slack.Attachment{
		{Text: "internal server error", Color: global.ColorcodeCrimson},
	}
	return msg
}
