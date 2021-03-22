package bitbucketforge

import (
	"fmt"
	"git-forge/configset"
	"git-forge/forge"
	"github.com/ktrysmt/go-bitbucket"
	"time"
)

func (f *BitBucketForge) GetAllPrTitles() ([]forge.PrTitle, error) {
	// now get us our auth token for the bitbucket api
	c := bitbucket.NewBasicAuth(f.cfg.User, f.cfg.Pass)

	cfg, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return nil, err
	}
	defer cfg.CommitConfig()

	fconfig, ferr := cfg.GetForgeRemoteSection()
	if ferr != nil {
		return nil, fmt.Errorf("Forge config is busted: %s\n", ferr)
	}

	_, slug, owner, _ := getRepoSlugAndOwner(fconfig.Parent.Url)

	propts := &bitbucket.PullRequestsOptions{
		Owner:    owner,
		RepoSlug: slug,
	}
	prs, rerr := c.Repositories.PullRequests.Gets(propts)
	if rerr != nil {
		return nil, rerr
	}
	var retprs []forge.PrTitle = make([]forge.PrTitle, 0)
	prmap := prs.(map[string]interface{})
	var count int = int(prmap["pagelen"].(float64))
	if int(prmap["size"].(float64)) < count {
		count = int(prmap["size"].(float64))
	}
	for i := 0; i < count; i++ {
		retprs = append(retprs, forge.PrTitle{
			Title: prmap["values"].([]interface{})[i].(map[string]interface{})["title"].(string),
			PrId:  int64(prmap["values"].([]interface{})[i].(map[string]interface{})["id"].(float64)),
		})
	}
	return retprs, nil
}

func (f *BitBucketForge) GetPr(idstring string) (*forge.PR, error) {

	cfg, err := gitconfigset.NewForgeConfigSet()
	if err != nil {
		return nil, err
	}
	defer cfg.CommitConfig()

	fconfig, ferr := cfg.GetForgeRemoteSection()
	if ferr != nil {
		return nil, fmt.Errorf("Forge config is busted: %s\n", ferr)
	}

	_, slug, owner, _ := getRepoSlugAndOwner(fconfig.Parent.Url)

	pullrequest, err := GetPrFromBitBucket(f.cfg.ApiBaseUrl, owner, slug, f.cfg.User, f.cfg.Pass, idstring)
	if err != nil {
		return nil, err
	}

	retpr := forge.PR{
		Unread:       true,
		CurrentToken: pullrequest.UpdatedOn.Format(time.UnixDate),
		Title:        pullrequest.Title,
		PrId:         int64(pullrequest.ID),
		PullSpec: forge.PrSpec{
			Source: forge.PrRemote{
				URL:        pullrequest.Source.Repository.Links.HTML.Href,
				BranchName: pullrequest.Source.Branch.Name,
			},
			Target: forge.PrRemote{
				URL:        pullrequest.Destination.Repository.Links.HTML.Href,
				BranchName: pullrequest.Destination.Branch.Name,
			},
		},
		Discussions: make([]forge.CommentData, 0),
	}

	commenterr := GetAllPrCommentsFromBitBucket(f.cfg.ApiBaseUrl, owner, slug, f.cfg.User, f.cfg.Pass, idstring, func(comments *PRComments, data interface{}) {
		myretpr := data.(*forge.PR)
		for i := 0; i < len(comments.Values); i++ {
			c := comments.Values[i]
			if c.Deleted == true {
				continue
			}
			newcomment := forge.CommentData{}
			newcomment.Id = c.ID
			newcomment.ParentId = c.Parent.ID
			newcomment.Author = c.User.DisplayName
			if c.Inline.Path == "" {
				newcomment.Type = forge.GENERAL
			} else {
				newcomment.Type = forge.INLINE
				newcomment.Path = c.Inline.Path
				newcomment.Offset = c.Inline.To
			}
			newcomment.Content = c.Content.Raw
			retpr.Discussions = append(myretpr.Discussions, newcomment)
		}
	}, &retpr)

	if commenterr != nil {
		return nil, commenterr
	}

	commiterr := GetAllPrCommitsFromBitBucket(f.cfg.ApiBaseUrl, owner, slug, f.cfg.User, f.cfg.Pass, idstring, func(commits *PRCommits, data interface{}) {
		myretpr := data.(*forge.PR)
		myretpr.Commits = make([]forge.Commit, 0)
		for i := 0; i < len(commits.Values); i++ {
			c := commits.Values[i]
			newcommit := forge.Commit{}
			newcommit.Comments = make([]forge.CommentData, 0)
			newcommit.Hash = c.Hash
			GetAllPrCommitCommentsFromBitBucket(c.Links.Comments.Href, f.cfg.User, f.cfg.Pass, func(ccomments *PrCommitComments, data interface{}) {
				mynewcommit := data.(*forge.Commit)
				for j := 0; j < len(ccomments.Values); j++ {
					cc := ccomments.Values[j]
					if cc.Deleted == true {
						continue
					}
					newcomitcomment := forge.CommentData{}
					newcomitcomment.Id = cc.ID
					newcomitcomment.ParentId = cc.Parent.ID
					newcomitcomment.Author = cc.User.DisplayName
					newcomitcomment.Content = cc.Content.Raw
					newcomitcomment.Path = cc.Inline.Path
					newcomitcomment.Offset = cc.Inline.To
					mynewcommit.Comments = append(mynewcommit.Comments, newcomitcomment)
				}
			}, &newcommit)
			myretpr.Commits = append(retpr.Commits, newcommit)
		}
	}, &retpr)

	if commiterr != nil {
		return nil, commiterr
	}
	return &retpr, nil
}
