package updater

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/google/go-github/github"
)

// GitHubEventsHandler is the event handler for the plugin.
// Implements server.GitHubEventHandler interface which contains the logic for incoming GitHub events
type GitHubEventsHandler struct {
	Ctx          context.Context
	Client       *github.Client
	Owner        string
	Path         string
	Organization string
	Dependency   string
}

// HandleEvent is an entry point for the plugin logic. This method is invoked by the Server when
// events are dispatched from the /hook service
func (gh *GitHubEventsHandler) HandleEvent(event interface{}, payload []byte) error {

	switch event.(type) {
	case *github.PushEvent:
		// this is a commit push, do something with it
		fmt.Println("Push event received.")
		// 1. Fetch all Repositories of the Organization
		// repos := FetchOrganizationRepositories(ctx, client, *organization)

		repo := "openshift-deployment-testing"
		// 2. Fork the Repository
		//ForkRepository(ctx, client, *organization, repo)

		// 3. Get Repository Contents for Computation
		_, baseFileSHA := GetRepositoryContents(gh.Ctx, gh.Client, gh.Owner, repo, gh.Path)

		// 4. Edit pom.xml, commit and push for Version Bump
		EditAndCommitRepositoryContent(gh.Ctx, gh.Client, gh.Owner, repo, gh.Path, baseFileSHA, gh.Dependency)

		// 5. Open a PR to the repo.
		OpenPullRequest(gh.Ctx, gh.Client, gh.Organization, gh.Owner, repo)

	case *github.PullRequestEvent:
		// this is a pull request, do something with it
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
	default:
		log.Printf("unknown event type")
	}
	return nil
}
