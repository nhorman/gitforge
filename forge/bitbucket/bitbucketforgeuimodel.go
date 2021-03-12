package bitbucketforge

import (
	"fmt"
	"git-forge/configset"
	"git-forge/forge"
	"github.com/ktrysmt/go-bitbucket"
	"io/ioutil"
	"net/http"
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

	req, err := http.NewRequest("GET", "https://"+f.cfg.ApiBaseUrl+"/repositories/"+owner+"/"+slug+"/pullrequests/"+idstring, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch PR json: %s", err)
	}
	req.SetBasicAuth(f.cfg.User, f.cfg.Pass)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	pullrequest, err := PrJsonToStruct(body)
	if err != nil {
		return nil, err
	}

	retpr := forge.PR{
		Title: pullrequest.Title,
		PrId:  int64(pullrequest.ID),
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
		Discussions: make([]forge.Discussion, 0),
	}

	creq, cerr := http.NewRequest("GET", "https://"+f.cfg.ApiBaseUrl+"/repositories/"+owner+"/"+slug+"/pullrequests/"+idstring+"/comments", nil)
	if cerr != nil {
		return nil, fmt.Errorf("Unable to fetch PR json: %s", cerr)
	}
	creq.SetBasicAuth(f.cfg.User, f.cfg.Pass)
	cresp, crerr := http.DefaultClient.Do(creq)
	if crerr != nil {
		return nil, crerr
	}
	defer cresp.Body.Close()

	cbody, _ := ioutil.ReadAll(cresp.Body)

	comments, commenterr := PrCommentsJsonToStruct(cbody)
	if commenterr != nil {
		return nil, commenterr
	}

	for i := 0; i < len(comments.Values); i++ {
		c := comments.Values[i]
		if c.Deleted == true {
			continue
		}
		newcomment := forge.Discussion{}
		newcomment.Id = c.ID
		newcomment.ParentId = c.Parent.ID
		newcomment.Author = c.User.DisplayName
		if c.Inline.Path == "" {
			newcomment.Type = forge.GENERAL
		} else {
			newcomment.Type = forge.INLINE
			newcomment.Inline.Path = c.Inline.Path
			newcomment.Inline.Offset = c.Inline.To
		}
		newcomment.Content = c.Content.Raw
		retpr.Discussions = append(retpr.Discussions, newcomment)
	}

	return &retpr, nil
}