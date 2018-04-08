package main

import (
	"github.com/google/go-github/github"
	"github.com/hemanik/automation-script/server"
	"github.com/hemanik/automation-script/updater"
	"golang.org/x/net/context"
)

func main() {
	updater.StartVersionUpdate(eventHandler, eventServer)
}

func eventHandler(ctx context.Context, githubClient *github.Client, organization string, owner string, path string, dependency string) server.GitHubEventHandler {
	return &updater.GitHubEventsHandler{Ctx: ctx, Client: githubClient, Organization: organization, Owner: owner, Path: path, Dependency: dependency}
}

func eventServer(webhookSecret []byte, eventHandler server.GitHubEventHandler) *server.Server {
	return &server.Server{
		GitHubEventHandler: eventHandler,
		HmacSecret:         webhookSecret,
	}
}
