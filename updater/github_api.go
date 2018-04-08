package updater

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
)

const (
	tempDirectory = "/tmp/temp/"
)

// FetchOrganizationRepositories fetches repos by organization
func FetchOrganizationRepositories(ctx context.Context, client *github.Client, organization string) []*github.Repository {
	fmt.Print("\n Fetching Repositories ........\n")
	repos, _, err := client.Repositories.ListByOrg(ctx, organization, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for i, repo := range repos {
		fmt.Printf(" %v. %s\n", i+1, repo.GetHTMLURL())
	}
	return repos
}

// ForkRepository forks a repo
func ForkRepository(ctx context.Context, client *github.Client, organization string, repo string) {
	fmt.Printf("\n Repo to be forked: %v \n\n", repo)

	forkedRepo, _, forkErr := client.Repositories.CreateFork(ctx, organization, repo, nil)

	if _, ok := forkErr.(*github.AcceptedError); ok {
		log.Println("scheduled on GitHub side")
	}

	fmt.Printf("\n Repository Fork Created %v", forkedRepo.GetHTMLURL())
}

// GetRepositoryContents fetches repository contents
func GetRepositoryContents(ctx context.Context, client *github.Client, forkOwner string, repo string, path string) (string, string) {
	fmt.Println("\n\n Caching Repository Contents...........")

	contents, _, _, err := client.Repositories.GetContents(ctx, forkOwner, repo, path, nil)
	if err != nil {
		fmt.Printf("Repositories.GetContents returned error: %v", err)
	}
	fileContent, _ := contents.GetContent()
	fileSha := contents.GetSHA()
	storeRepositoryContentsInTempDirectory(fileContent, path)
	return fileContent, fileSha
}

// BumpProjectVersion bumps new version
func BumpProjectVersion(dependency string) {
	if err := exec.Command("sh", "version-bump.sh", dependency).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("\n Successfully bumped new version.")
}

// EditAndCommitRepositoryContent updates repo content
func EditAndCommitRepositoryContent(ctx context.Context, client *github.Client, forkOwner string, repo string, path string, baseTreeSHA string, dependency string) {
	fmt.Print("\n Editing POM .................\n")

	BumpProjectVersion(dependency)

	message := "chore: bumps new project version."
	content := []byte(getNewRepositoryContentsFromTempDirectory(path))
	sha := baseTreeSHA
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: content,
		SHA:     &sha,
		//Committer: &github.CommitAuthor{Name: github.String("arquillian"), Email: github.String("arquillian-ike@redhat.com")},
	}
	updateResponse, _, err := client.Repositories.UpdateFile(ctx, forkOwner, repo, path, repositoryContentsOptions)
	if err != nil {
		fmt.Printf("Repositories.UpdateFile returned error: %v", err)
	}

	fmt.Printf("\n Added commit %s\n", updateResponse.GetMessage())
}

// OpenPullRequest opens a PR
func OpenPullRequest(ctx context.Context, client *github.Client, organization string, forkOwner string, repo string) {
	fmt.Print("\n Opening a new PR ................. ")

	input := &github.NewPullRequest{
		Title: github.String("Bumps New Project Version"),
		Body:  github.String("Bumps latest version"),
		Head:  github.String(forkOwner + ":master"),
		Base:  github.String("master"),
	}

	pull, _, err := client.PullRequests.Create(ctx, organization, repo, input)
	if err != nil {
		fmt.Printf("PullRequests.Create returned error: %v", err)
	}
	fmt.Printf("%s\n", pull.GetHTMLURL())
}

func storeRepositoryContentsInTempDirectory(fileContent string, path string) {
	if _, err := os.Stat(tempDirectory); os.IsNotExist(err) {
		os.Mkdir(tempDirectory, 0777)
	}

	if err := ioutil.WriteFile(tempDirectory+path, []byte(fileContent), 0644); err != nil {
		log.Fatal(err)
	}
}

func getNewRepositoryContentsFromTempDirectory(path string) string {
	content, err := ioutil.ReadFile(tempDirectory + path)
	if err != nil {
		log.Fatal(err)
	}
	os.RemoveAll(tempDirectory) // clean up
	return string(content)
}
