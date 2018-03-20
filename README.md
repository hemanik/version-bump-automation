# Version Bump Automation

Version Bump Script for Maven projects leverages [Versions Maven Plugin](http://www.mojohaus.org/versions-maven-plugin/)
to update project dependencies specified in your project's pom.xml and GitHub API's to interact with GitHub to push and 
commit the changes to upstream repository.

The script does the following:

1. Lists the repositories by the Organization. 
2. Forks a repository from the Organization.
3. Caches the contents of the `pom.xml` to the temporary directory(`\tmp\temp`).
4. Executes the Versions Maven Plugin against this directory (`mvn versions:use-latest-versions`). 
Use property `-Dincludes=groupId:artifactId:type:classifier` to update specific dependency.  
5. Commits the changed file contents and pushes it to the remote repository.
6. Opens a PR to the base repository for the change.

Usage
-----
Assuming you have Golang installed execute the program as follows: 

`go run main.go github-api.go -org=openshiftio-vertx-boosters -owner=hemanik -dependency=io.rest-assured:*`

Command Line Arguments For The Script 
-------------------------------------

* `token`      : github auth-token with repo permissions enabled.
* `org`        : github organization to be searched for. 
* `owner`      : github user for the forked repository. 
* `filepath`   : path for the pom.xml file. By default takes the `pom.xml` in the root directory.
* `dependency` : `groupId:artifactId:type:classifier` combination to be updated.

