package githubforge

import (
	"context"
	"fmt"
	"git-forge/configset"
	"git-forge/forge"
	"strconv"
	"time"

	"github.com/google/go-github/v33/github"
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
		titles = append(titles, forge.PrTitle{*pr.Title, int64(*pr.Number)})
	}

	return titles, nil

}

func (f *GitHubForge) GetPr(idstring string) (*forge.PR, error) {
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

	prnum, _ := strconv.Atoi(idstring)
	pr, _, perr := client.PullRequests.Get(ctx, powner, pslug, prnum)
	if perr != nil {
		return nil, perr
	}

	retpr := forge.PR{
		Unread:       true,
		CurrentToken: pr.UpdatedAt.Format(time.UnixDate),
		Title:        *pr.Title,
		PrId:         int64(*pr.Number),
		PullSpec: forge.PrSpec{
			Source: forge.PrRemote{
				URL:        *pr.Head.Repo.GitURL,
				BranchName: *pr.Head.Ref,
			},
			Target: forge.PrRemote{
				URL:        *pr.Base.Repo.GitURL,
				BranchName: *pr.Base.Ref,
			},
		},
		Discussions: make([]forge.CommentData, 0),
	}

	comments, _, ierr := client.Issues.ListComments(ctx, powner, pslug, prnum, nil)
	if ierr != nil {
		return nil, ierr
	}

	for i := 0; i < len(comments); i++ {
		c := comments[i]
		newc := forge.CommentData{}
		newc.Id = int(*c.ID)
		newc.ParentId = 0         //Issue comments can't be nested
		newc.Type = forge.GENERAL //Issue comments are our General comments
		if c.User.Name != nil {
			newc.Author = *c.User.Name
		} else {
			newc.Author = *c.User.Login
		}
		newc.Content = *c.Body
		retpr.Discussions = append(retpr.Discussions, newc)
	}

	prc, _, perr := client.PullRequests.ListComments(ctx, powner, pslug, prnum, nil)
	if perr != nil {
		return nil, perr
	}

	for i := 0; i < len(prc); i++ {
		c := prc[i]
		newc := forge.CommentData{}
		newc.Id = int(*c.ID)
		if c.User.Name != nil {
			newc.Author = *c.User.Name
		} else {
			newc.Author = *c.User.Login
		}
		if c.InReplyTo != nil {
			newc.ParentId = int(*c.InReplyTo)
		} else {
			newc.ParentId = 0
		}
		newc.Type = forge.GENERAL //Issue comments are our General comments
		newc.Content = *c.Body
		retpr.Discussions = append(retpr.Discussions, newc)
	}

	return &retpr, nil
}

func (f *GitHubForge) RefreshPr(pr *forge.PR) (chan *forge.UpdatedPR, error) {
	return nil, fmt.Errorf("Not Implemented Yet")
}
