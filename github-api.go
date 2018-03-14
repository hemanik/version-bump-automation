package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// ORGANIZATION base repository owner
	ORGANIZATION = "snowdrop"
	// FORK_OWNER owner for the foked repo
	FORK_OWNER = "hemanik"
)

func fetchOrganisationRepositories(client *github.Client) []*github.Repository {
	repos, _, err := client.Repositories.ListByOrg(context.Background(), ORGANIZATION, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for i, repo := range repos {
		fmt.Printf("\n %v. %s\n", i+1, *repo.HTMLURL)
	}

	return repos
}

func forkRepository(client *github.Client, repo string) {
	forkedRepo, _, forkErr := client.Repositories.CreateFork(context.Background(), ORGANIZATION, repo, nil)

	if _, ok := forkErr.(*github.AcceptedError); ok {
		log.Println("\n scheduled on GitHub side")
	}

	fmt.Printf("\n Forked Repo: %v", forkedRepo)
}

func fetchBaseTree(client *github.Client, repo string) string {
	// Get reference to HEAD
	ref, _, err := client.Git.GetRef(context.Background(), FORK_OWNER, repo, "refs/heads/master")
	if err != nil {
		fmt.Printf("Git.GetRef returned error: %v", err)
	}
	latestCommitSHA := ref.Object.GetSHA()
	fmt.Printf("\n SHA_LATEST_COMMIT: %v", latestCommitSHA)

	// Grab the commit HEAD point to.
	commit, _, err := client.Git.GetCommit(context.Background(), ORGANIZATION, repo, latestCommitSHA)
	if err != nil {
		fmt.Printf("Git.GetCommit returned error: %v", err)
	}
	baseTreeSHA := commit.GetSHA()
	fmt.Printf("\n SHA_BASE_TREE: %v", baseTreeSHA)

	// Get hold of the tree commit points to.
	treeInput := []github.TreeEntry{
		{
			Path:    github.String("tests/pom.xml"),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String("file content"),
		},
	}

	tree, _, err := client.Git.CreateTree(context.Background(), FORK_OWNER, repo, baseTreeSHA, treeInput)
	if err != nil {
		fmt.Printf("Git.CreateTree returned error: %v", err)
	}
	newTreeSHA := tree.GetSHA()
	fmt.Printf("\n SHA_NEW_TREE: %v", newTreeSHA)

	return baseTreeSHA
}

func editAndCommitRepositoryContent(client *github.Client, repo string, baseTreeSHA string) {
	message := "version bump"
	content := []byte("file content")
	sha := baseTreeSHA
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message:   &message,
		Content:   content,
		SHA:       &sha,
		Committer: &github.CommitAuthor{Name: github.String("arquillian"), Email: github.String("arquillian-ike@redhat.com")},
	}
	updateResponse, _, err := client.Repositories.UpdateFile(context.Background(), FORK_OWNER, repo, "tests/pom.xml", repositoryContentsOptions)
	if err != nil {
		fmt.Printf("Repositories.UpdateFile returned error: %v", err)
	}

	responseSHA := updateResponse.GetSHA()
	fmt.Printf("%s\n", responseSHA)
}

func openPullRequest(client *github.Client, repo string) {
	prInput := &github.NewPullRequest{
		Title: github.String("Bumps new version"),
		Body:  github.String("Bumps Latest Version"),
		Head:  github.String("hemanik:development"),
		Base:  github.String("master"),
	}

	pull, _, err := client.PullRequests.Create(context.Background(), ORGANIZATION, "test1", prInput)
	if err != nil {
		fmt.Printf("PullRequests.Create returned error: %v", err)
	}
	fmt.Printf("%s\n", pull.GetHTMLURL())
}

func main() {

	// Authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "47b3edcff1ef72fbe48cf892417bc4bc65cbae6e"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// 1. Fetch All Repositories of the Organisation
	fmt.Print("\n Fetching Repositories ........")
	repos := fetchOrganisationRepositories(client)

	// 2. Fork all the Repositories
	repo := *repos[4].Name
	fmt.Printf("\n Repo to be forked: %v", repo)
	forkRepository(client, repo)

	// 3. Get Base Tree For Computation
	fmt.Println("\n Calculating SHA's ............. ")
	baseTreeSHA := fetchBaseTree(client, repo)

	// 4. Edit pom.xml, commit and push for Version Bump
	fmt.Print("\n Editing POM .................")
	editAndCommitRepositoryContent(client, repo, baseTreeSHA)

	// 5. Open a PR to the repo.
	fmt.Print("\n Opening a new PR ................. ")
	openPullRequest(client, repo)
}
