package updater

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/hemanik/automation-script/server"
	"github.com/hemanik/automation-script/utils"
	"golang.org/x/oauth2"
)

var (
	port              = flag.Int("port", 8888, "Port to listen on.")
	githubTokenFile   = flag.String("github-token-file", "/etc/github/oauth", "Path to the file containing the GitHub OAuth secret.")
	webhookSecretFile = flag.String("hmac-secret-file", "/etc/webhook/hmac", "Path to the file containing the GitHub HMAC secret.")
	organization      = flag.String("org", "", "github organization")
	owner             = flag.String("owner", "", "repository owner")
	path              = flag.String("filepath", "pom.xml", "path of the file to be updated")
	dependency        = flag.String("dependency", "", "dependency to be updated")
)

// EventHandlerCreator is a func type that creates server.GitHubEventHandler instance which is the central point for
// the plugin logic
type EventHandlerCreator func(context context.Context, client *github.Client, organization string, owner string, path string, dependency string) server.GitHubEventHandler

// ServerCreator is a func type that wires Server and server.GitHubEventHandler together
type ServerCreator func(hmacSecret []byte, evenHandler server.GitHubEventHandler) *server.Server

// StartVersionUpdate hjguhihk
func StartVersionUpdate(newEventHandler EventHandlerCreator, newServer ServerCreator) {

	flag.Parse()

	webhookSecret, err := utils.LoadSecret(*webhookSecretFile)
	if err != nil {
		fmt.Printf("unable to load webhook secret from %q", *webhookSecretFile)
	}

	oauthSecret, err := utils.LoadSecret(*githubTokenFile)
	if err != nil {
		fmt.Printf("unable to load oauth token from %q", *githubTokenFile)
	}

	// Authentication
	ctx := context.Background()
	token := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(oauthSecret)},
	)
	githubClient := github.NewClient(oauth2.NewClient(ctx, token))

	handler := newEventHandler(ctx, githubClient, *organization, *owner, *path, *dependency)
	pluginServer := newServer(webhookSecret, handler)

	port := strconv.Itoa(*port)
	log.Printf("Starting server on port %s", port)

	http.Handle("/", pluginServer)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("failed to start server on port %s", port)
	}
}
