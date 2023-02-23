package github

import (
	"context"
	"fmt"
	"github.com/coreos/go-semver/semver"
	gogithub "github.com/google/go-github/v50/github"
)

type Github interface {
	CreateNewTag(owner, repo, sha, releaseTag, preReleaseTag string) (string, error)
	CreateRelease(tag, preReleaseTag string) error
}
type github struct {
	client *gogithub.Client
	ctx    context.Context
	owner  string
	repo   string
}

func New(githubToken, owner, repo string) *github {
	ctx := context.Background()
	client := gogithub.NewTokenClient(ctx, githubToken)
	return &github{
		ctx:    ctx,
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

// CreateNewTag creates a new tag in the repository based on previous latest tag and increments its version or returns an error
func (g *github) CreateNewTag(sha, releaseTag, preReleaseTag string) (string, error) {
	tag := "0.0.1"
	latest, err := g.getLatestTag()
	if err != nil {
		if err.Error() == "no tags found" {
			return g.createNewTag(tag, sha)
		}
		return "", err
	}

	v, err := semver.NewVersion(latest)
	if err != nil {
		return "", err
	}

	if preReleaseTag != "" {
		v.PreRelease = semver.PreRelease(preReleaseTag)
	} else {
		v.PreRelease = ""
	}

	switch releaseTag {
	case "major":
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case "minor":
		v.Minor++
		v.Patch = 0
	case "patch":
		v.Patch++
	case "":
		v.Patch++
	default:
		v.Patch++
	}
	return g.createNewTag(v.String(), sha)
}

// createNewTag creates a new tag in the repository or returns an error
func (g *github) createNewTag(tag, sha string) (string, error) {
	_, _, err := g.client.Git.CreateTag(g.ctx, g.owner, g.repo, &gogithub.Tag{
		Tag:     &tag,
		Message: &tag,
		Object: &gogithub.GitObject{
			SHA:  &sha,
			Type: gogithub.String("commit"),
		},
	})
	if err != nil {
		return "", err
	}
	_, _, err = g.client.Git.CreateRef(g.ctx, g.owner, g.repo, &gogithub.Reference{
		Ref: gogithub.String(fmt.Sprintf("refs/tags/%s", tag)),
		Object: &gogithub.GitObject{
			SHA: &sha,
		},
	})
	if err != nil {
		return "", err
	}
	return tag, nil
}

// getLatestTag returns the latest tag of the repository or an error
func (g *github) getLatestTag() (string, error) {
	tags, _, err := g.client.Repositories.ListTags(g.ctx, g.owner, g.repo, nil)
	if err != nil {
		return "", err
	}
	var versions []*semver.Version
	if len(tags) > 0 {
		for _, tag := range tags {
			v, err := semver.NewVersion(tag.GetName())
			if err != nil {
				return "", err
			}
			versions = append(versions, v)
		}
	} else {
		return "", fmt.Errorf("no tags found")
	}
	semver.Sort(versions)
	latestTag := versions[len(versions)-1].String()
	return latestTag, nil
}

// CreateRelease creates a new release in the repository or returns an error
func (g *github) CreateRelease(tag, preReleaseTag string) error {
	preRelease := false
	if preReleaseTag != "" {
		preRelease = true
	}
	_, _, err := g.client.Repositories.CreateRelease(g.ctx, g.owner, g.repo, &gogithub.RepositoryRelease{
		TagName:    &tag,
		Name:       &tag,
		Body:       gogithub.String(""),
		Draft:      gogithub.Bool(false),
		Prerelease: gogithub.Bool(preRelease),
	})
	if err != nil {
		return err
	}
	return nil
}
