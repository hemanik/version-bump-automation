# Version Bump Automation

Version Bump Script for Maven projects leverages [Versions Maven Plugin]
(http://www.mojohaus.org/versions-maven-plugin/)
to update project dependencies specified in your project's pom.xml and GitHub API's to interact with 
GitHub to push and commit the changes to upstream repository.

The script does the following:

1. Lists the repositories by the Organization. 
2. Forks a repository from the Organization.
3. Starts a server and listens for any push events received on the repository.
3. Caches the contents of the `pom.xml` to the temporary directory(`\tmp\temp`).
4. Executes the Versions Maven Plugin against this directory (`mvn versions:use-latest-versions`). 
Use property `-Dincludes=groupId:artifactId:type:classifier` to update specific dependency.  
5. Commits the changed file contents and pushes it to the remote repository.
6. Opens a PR to the base repository for the change.

Configuring Webhook 
-------------------

Using ngrok

To deliver webhooks, Github needs a publicly reachable address or DNS name. We will use ngrok to 
temporarily expose the development server to the internet.

First, [download ngrok](https://ngrok.com/download). Then, expose port 8888 to the web. We will use 
port 8888 to run our Go webserver.

`~/Downloads/ngrok http 8888`

You’ll see a new screen with a random Forwarding https address. Mine is https://11f018ee.ngrok.io, but 
yours will have a different subdomain. Copy this address, we will use it to configure the github 
webhook URL.

Keep the ngrok session running. If you restart it, you’ll get a completely different URL, requiring 
you to reconfigure the Github settings.

Configure Github to send webhooks

In a repo’s settings(you must own the repository) page, click the “Webhook” section.

* Click the Add Webhook button.
* In the Payload URL field, type your ngrok https address, and add `/` as the path.
Example: https://11f018ee.ngrok.io/
* In the Secret field insert your hmac secret token. 

Select the list of events you’re interested in. In my case, I selected the Push category, which will 
trigger a Webhook whenever someone pushes to the repository.

Usage
-----

Configuring tokens for use in the program:

* store the webhook secret selected while configuring webhook in a file (Eg. `hmac.token`).
* store the github auth token with repo permissions enabled in a file (Eg. `oauth.token`). 

Starting the ngrok server:

`~/Downloads/ngrok http 8888`

Assuming you have Golang installed, server started webhook and tokens configured execute the program 
as follows: 

`go run main.go -githubTokenFile=config/oauth.token -webhookSecretFile=config/hmac.token 
-repo=openshift-deployment-testing -owner=hemanik -dependency=io.rest-assured:*`

Whenever a push event is trigger on the repository, the updater script is triggered which updates the 
dependencies and opens a new commit and a pull request with updated dependencies.

Command Line Arguments For The Script 
-------------------------------------

* `port`              : port on which the server listens. (Default is 8888)
* `githubTokenFile`   : config file that stores github auth-token with repo permissions enabled.
* `webhookSecretFile` : config file that stores webhook secret.
* `org`               : github organization to be searched for. 
* `owner`             : github user for the forked repository. 
* `repo`              : github repo for with the webhook is enabled.
* `dependency`        : `groupId:artifactId:type:classifier` combination to be updated. (If not set 
all dependencies will be updated.)
* `filepath`          : path for the pom.xml file. By default takes the `pom.xml` in the root 
directory.