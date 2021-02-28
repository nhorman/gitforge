package githubforge

import (
	"context"
	"fmt"
	"git-forge/configset"

	"github.com/google/go-github/v33/github"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/config"
)

func (f *GitHubForge) GetAllPrTitles() ([]string, []int64, error) {

	cfg, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return nil, nil, err
	}
	defer cfg.CommitConfig()
	fconfig, err := cfg.GetForgeRemoteSection()
	if err != nil {
		return nil, nil, err
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
		return nil, nil, lerr
	}

	var titles []string = make([]string, 0)
	var ids []int64 = make([]int64, 0)
	for _, pr := range prs {
		titles = append(titles, *pr.Title)
		ids = append(ids, *pr.ID)
	}

	return titles, ids, nil

}

func (f *GitHubForge) AddWatchPr(idstring string) error {
	return fmt.Errorf("Not Implemented Yet")
}
