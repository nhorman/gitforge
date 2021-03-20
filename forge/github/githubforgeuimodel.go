package githubforge

import (
	"context"
	"fmt"
	"git-forge/configset"
	"git-forge/forge"

	"github.com/google/go-github/v33/github"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/config"
)

func (f *GitHubForge) GetAllPrTitles() ([]forge.PrTitle, error) {

	cfg, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return nil, err
	}
	defer cfg.CommitConfig()
	fconfig, err := cfg.GetForgeRemoteSection()
	if err != nil {
		return nil, err
	}

	transport := &github.BasicAuthTransport{
		Username: f.cfg.User,
		Password: f.cfg.Pass,
	}

	client := github.NewClient(transport.Client())
	ctx := context.Background()
	_, pslug, powner, _ := getRepoSlugAndOwner(fconfig.Parent.Url)

	PRList := &github.PullRequestListOptions{}

	prs, _, lerr := client.PullRequests.List(ctx, powner, pslug, PRList)
	if lerr != nil {
		return nil, lerr
	}

	var titles []forge.PrTitle = make([]forge.PrTitle, 0)

	for _, pr := range prs {
		titles = append(titles, forge.PrTitle{*pr.Title, *pr.ID})
	}

	return titles, nil

}

func (f *GitHubForge) AddWatchPr(idstring string) error {
	return fmt.Errorf("Not Implemented Yet")
}
