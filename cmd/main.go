package main

import (
	"actions-semver-release/pkg/github"
	"actions-semver-release/pkg/util"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	repoOwner     = ""
	repoName      = ""
	commitSha     = ""
	githubToken   = ""
	releaseTag    = ""
	preReleaseTag = ""
	createRelease = false
)

func init() {
	repoPath := os.Getenv("GITHUB_REPOSITORY")
	if repoPath == "" {
		log.Fatalln("cannot get GITHUB_REPOSITORY env")
	}
	parts := strings.Split(repoPath, "/")
	repoOwner = parts[0]
	repoName = parts[1]
	if repoOwner == "" || repoName == "" {
		log.Fatalln("cannot get the necessary values")
	}
	commitSha = os.Getenv("GITHUB_SHA")
	if commitSha == "" {
		log.Fatalln("cannot get GITHUB_SHA env")
	}
	log.Println(fmt.Sprintf("repoOwner: %s, repoName: %s, commitSha: %s", repoOwner, repoName, commitSha))

	githubToken = os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatalln("cannot get GITHUB_TOKEN")
	}

	releaseTag = os.Getenv("RELEASE_TAG")
	preReleaseTag = os.Getenv("PRE_RELEASE_TAG")

	var err error
	createRelease, err = util.GetEnvBool(os.Getenv("CREATE_RELEASE"))
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	g := github.New(githubToken, repoOwner, repoName)
	tag, err := g.CreateNewTag(commitSha, releaseTag, preReleaseTag)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprint("created tag: ", tag))
	util.SetGithubOutput("tag", tag)
	if createRelease {
		err = g.CreateRelease(tag, preReleaseTag)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
