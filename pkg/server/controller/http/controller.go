package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/infrastructure"
	"github.com/cloudnativedaysjp/dreamkast-releasebot/pkg/server/view"
)

const (
	temporaryBranchForRelease = "bot/release"
)

// InteractionMsgController is Controller
type InteractionMsgController struct {
	logger *logrus.Logger

	gitApi          infrastructure.GitAPIDao
	gitCommand      infrastructure.GitCommandDao
	cache           infrastructure.CacheDao
	slackMsg        infrastructure.SlackMsgDao
	baseBranch      string
	enableAutoMerge bool
}

// NewInteractionMsgController is initialize function
func NewInteractionMsgController(logger *logrus.Logger, gitapi infrastructure.GitAPIDao, gitcommand infrastructure.GitCommandDao, cache infrastructure.CacheDao, slackmsg infrastructure.SlackMsgDao, baseBranch string, enableAutoMerge bool) *InteractionMsgController {
	return &InteractionMsgController{logger, gitapi, gitcommand, cache, slackmsg, baseBranch, enableAutoMerge}
}

func (ic *InteractionMsgController) SelectedCancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	originalMessage := ctx.Value("originalMessage").(slack.Message)

	/* View */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(view.IntMsgCancel(originalMessage)); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ic *InteractionMsgController) SelectedRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	originalMessage := ctx.Value("originalMessage").(slack.Message)
	callbackAction := ctx.Value("callbackAction").(slack.AttachmentAction)
	callbackId := ctx.Value("callbackID").(string)

	/* Validation */
	if callbackAction.Type != "select" || callbackAction.SelectedOptions[0].Value == "" {
		err := fmt.Errorf(`invalid request`)
		ic.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value := strings.Split(callbackAction.SelectedOptions[0].Value, "__")
	if len(value) < 2 {
		err := fmt.Errorf(`invalid request`)
		ic.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	org := value[0]
	repo := value[1]

	// set to statestore
	if err := ic.cache.Write(fmt.Sprintf("%s_org", callbackId), org); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := ic.cache.Write(fmt.Sprintf("%s_repo", callbackId), repo); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/* View */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(view.IntMsgSelectLevel(callbackId, originalMessage, org, repo)); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ic *InteractionMsgController) SelectedReleaseLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	originalMessage := ctx.Value("originalMessage").(slack.Message)
	callbackAction := ctx.Value("callbackAction").(slack.AttachmentAction)
	callbackId := ctx.Value("callbackID").(string)

	/* Validation */
	if callbackAction.Type != "button" || callbackAction.Name == "" {
		err := fmt.Errorf(`invalid request`)
		ic.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	level := callbackAction.Name

	// set to statestore
	if err := ic.cache.Write(fmt.Sprintf("%s_level", callbackId), level); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get from statestore
	org, err := ic.cache.Read(fmt.Sprintf("%s_org", callbackId))
	if err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	repo, err := ic.cache.Read(fmt.Sprintf("%s_repo", callbackId))
	if err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/* View */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(view.IntMsgConfirmRelease(callbackId, originalMessage, org, repo, level)); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ic *InteractionMsgController) SelectedConfirmRelease(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	originalMessage := ctx.Value("originalMessage").(slack.Message)
	callbackAction := ctx.Value("callbackAction").(slack.AttachmentAction)
	channelId := ctx.Value("channelID").(string)
	callbackId := ctx.Value("callbackID").(string)

	// for debug
	ctx = context.Background()

	/* Validation */
	if callbackAction.Type != "button" {
		err := fmt.Errorf(`invalid request`)
		ic.logger.Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get from statestore
	org, err := ic.cache.Read(fmt.Sprintf("%s_org", callbackId))
	if err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	repo, err := ic.cache.Read(fmt.Sprintf("%s_repo", callbackId))
	if err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	level, err := ic.cache.Read(fmt.Sprintf("%s_level", callbackId))
	if err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/* Logic */
	// TODO: Lock
	go func() {
		defer func() {
			if err := ic.cache.Remove(fmt.Sprintf("%s_org", callbackId)); err != nil {
				ic.logger.Warn(fmt.Sprintf("failed to remove cache: %v", fmt.Sprintf("%s_org", callbackId)))
			}
			if err := ic.cache.Remove(fmt.Sprintf("%s_repo", callbackId)); err != nil {
				ic.logger.Warn(fmt.Sprintf("failed to remove cache: %v", fmt.Sprintf("%s_repo", callbackId)))
			}
			if err := ic.cache.Remove(fmt.Sprintf("%s_level", callbackId)); err != nil {
				ic.logger.Warn(fmt.Sprintf("failed to remove cache: %v", fmt.Sprintf("%s_level", callbackId)))
			}
		}()

		// clone repo to working dir
		repoDir, err := ic.gitCommand.Clone(ctx, org, repo)
		if err != nil {
			ic.logger.Error(err)
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
			}
			return
		}
		// remove working dir finally
		defer func() {
			if err := ic.gitCommand.Remove(ctx, repoDir); err != nil {
				ic.logger.Error(err)
				return
			}
		}()
		// switch -> empty commit -> push
		if err := ic.gitCommand.SwitchNewBranch(ctx, repoDir, temporaryBranchForRelease); err != nil {
			ic.logger.Error("SwitchNewBranch() failed. ", err)
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
			}
			return
		}
		if err := ic.gitCommand.CommitAll(ctx, repoDir, "[Bot] for release!!"); err != nil {
			ic.logger.Error("CommitAll() failed. ", err)
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
			}
			return
		}
		if err := ic.gitCommand.Push(ctx, repoDir); err != nil {
			ic.logger.Error("Push() failed. ", err)
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
			}
			return
		}
		defer func() {
			if err := ic.gitApi.DeleteBranch(ctx, org, repo, temporaryBranchForRelease); err != nil {
				ic.logger.Warn(fmt.Sprintf("failed to remove remote branch: %v", temporaryBranchForRelease))
			}
		}()
		// create -> label -> merge PullRequest
		prNum, err := ic.gitApi.CreatePullRequest(ctx, org, repo, temporaryBranchForRelease, ic.baseBranch, "[dreamkast-releasebot] Automatic Release", "Automatic Release")
		if err != nil {
			ic.logger.Error("CreatePullRequest() failed. ", err)
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
			}
			return
		}
		if err := ic.gitApi.LabelPullRequest(ctx, org, repo, prNum, level); err != nil {
			ic.logger.Error("LabelPullRequest() failed. ", err)
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
			}
			return
		}
		if ic.enableAutoMerge {
			if err := ic.gitApi.MergePullRequest(ctx, org, repo, prNum); err != nil {
				ic.logger.Error("MergePullRequest() failed. ", err)
				if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseFailed(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
					ic.logger.Error("UpdateMessage() failed. ", err)
				}
				return
			}
			// update Slack Message
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleaseDone(originalMessage, org, repo, level), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
				return
			}
		} else {
			// update Slack Message
			if err := ic.slackMsg.UpdateMessage(ctx, view.IntMsgReleasePRLink(originalMessage, org, repo, level, prNum), channelId, originalMessage.Timestamp); err != nil {
				ic.logger.Error("UpdateMessage() failed. ", err)
				return
			}
		}
	}()

	/* View */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(view.IntMsgReleaseProcessing(originalMessage, org, repo, level)); err != nil {
		ic.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
