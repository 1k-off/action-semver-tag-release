package main

import (
	"actions-semver-release/pkg/github"
	"actions-semver-release/pkg/util"
	"flag"
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
	assets        []string
	repoRoot      = ""
)

func init() {
	repoRoot = os.Getenv("GITHUB_WORKSPACE")

	repoOwner = *flag.String("repo-owner", "", "Repository owner name")
	repoName = *flag.String("repo-name", "", "Repository name")
	commitSha = *flag.String("commit-sha", "", "Commit SHA")
	githubToken = *flag.String("github-token", "", "GitHub token")
	releaseTag = *flag.String("release-tag", "", "Release tag")
	preReleaseTag = *flag.String("pre-release-tag", "", "Pre-release tag")
	createRelease = *flag.Bool("create-release", false, "Create release")
	assetsStr := *flag.String("assets", "", "Assets")
	assets = util.GetStringAsArray(assetsStr)
}

func getEnvIfFlagEmpty() {
	if repoOwner == "" || repoName == "" {
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
	}
	if commitSha == "" {
		commitSha = os.Getenv("GITHUB_SHA")
		if commitSha == "" {
			log.Fatalln("cannot get GITHUB_SHA env")
		}
	}
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			log.Fatalln("cannot get GITHUB_TOKEN")
		}
	}
	if releaseTag == "" {
		releaseTag = os.Getenv("RELEASE_TAG")
	}
	if preReleaseTag == "" {
		preReleaseTag = os.Getenv("PRE_RELEASE_TAG")
	}

	if !createRelease {
		var err error
		createRelease, err = util.GetEnvBool("CREATE_RELEASE")
		if err != nil {
			log.Fatalln(err)
		}
	}

	if createRelease {
		if len(assets) == 0 {
			assets = util.GetEnvArray("ASSETS")
		}
	}
}
func main() {
	flag.Parse()
	getEnvIfFlagEmpty()
	log.Println(fmt.Sprintf("repoOwner: %s, repoName: %s, commitSha: %s", repoOwner, repoName, commitSha))

	g := github.New(githubToken, repoOwner, repoName, repoRoot)
	tag, err := g.CreateNewTag(commitSha, releaseTag, preReleaseTag)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprint("created tag: ", tag))
	util.SetGithubOutput("tag", tag)
	if createRelease {
		releaseId, err := g.CreateRelease(tag, preReleaseTag)
		if err != nil {
			log.Fatalln(err)
		}
		if len(assets) > 0 {
			err = g.UploadReleaseAssets(releaseId, assets)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
