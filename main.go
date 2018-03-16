package main

import (
	"context"
	"flag"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// TOKEN Auth token for GitHub Authentication
	TOKEN = "Your GitHub Auth Token"
)

func main() {

	// Authentication
	ctx := context.Background()

	accessToken := flag.String("token", TOKEN, "github access token")

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tokenizedClient := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(tokenizedClient)

	organization := flag.String("org", "", "github organization")
	owner := flag.String("owner", "", "repository owner")
	path := flag.String("filepath", "pom.xml", "path of the file to be updated")
	dependency := flag.String("dependency", "", "dependency to be updated")

	flag.Parse()

	// 1. Fetch all Repositories of the Organization
	repos := FetchOrganizationRepositories(ctx, client, *organization)

	repo := *repos[4].Name
	// 2. Fork the Repository
	ForkRepository(ctx, client, *organization, repo)

	// 3. Get Repository Contents for Computation
	_, baseFileSHA := GetRepositoryContents(ctx, client, *owner, repo, *path)

	// 4. Edit pom.xml, commit and push for Version Bump
	EditAndCommitRepositoryContent(ctx, client, *owner, repo, *path, baseFileSHA, *dependency)

	// 5. Open a PR to the repo.
	OpenPullRequest(ctx, client, *organization, *owner, repo)
}
