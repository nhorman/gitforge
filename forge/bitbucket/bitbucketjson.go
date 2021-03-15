package bitbucketforge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LinkTuple struct {
	Raw    string `json:"raw,omitempty"`
	Markup string `json:"markup,omitempty"`
	HTML   string `json:"html,omitempty"`
	Href   string `json:"href,omitempty"`
	Type   string `json:"type,omitempty"`
}

type PullRequest struct {
	Rendered struct {
		Description LinkTuple `json: "description"`
		Title       LinkTuple `json:"title"`
	} `json:"rendered"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Links       struct {
		Decline        LinkTuple `json:"decline"`
		Diffstat       LinkTuple `json:"diffstat"`
		Commits        LinkTuple `json:"commits"`
		Self           LinkTuple `json:"self"`
		Comments       LinkTuple `json:"comments"`
		Merge          LinkTuple `json:"merge"`
		HTML           LinkTuple `json:"html"`
		Activity       LinkTuple `json:"activity"`
		RequestChanges LinkTuple `json:"request-changes"`
		Diff           LinkTuple `json:"diff"`
		Approve        LinkTuple `json:"approve"`
		Statuses       LinkTuple `json:"statuses"`
	} `json:"links"`
	Title             string        `json:"title"`
	CloseSourceBranch bool          `json:"close_source_branch"`
	Reviewers         []interface{} `json:"reviewers"`
	ID                int           `json:"id"`
	Destination       struct {
		Commit struct {
			Hash  string `json:"hash"`
			Type  string `json:"type"`
			Links struct {
				Self LinkTuple `json:"self"`
				HTML LinkTuple `json:"html"`
			} `json:"links"`
		} `json:"commit"`
		Repository struct {
			Links struct {
				Self   LinkTuple `json:"self"`
				HTML   LinkTuple `json:"html"`
				Avatar LinkTuple `json:"avatar"`
			} `json:"links"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"repository"`
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
	CreatedOn time.Time `json:"created_on"`
	Summary   LinkTuple `json:"summary"`
	Source    struct {
		Commit struct {
			Hash  string `json:"hash"`
			Type  string `json:"type"`
			Links struct {
				Self LinkTuple `json:"self"`
				HTML LinkTuple `json:"html"`
			} `json:"links"`
		} `json:"commit"`
		Repository struct {
			Links struct {
				Self   LinkTuple `json:"self"`
				HTML   LinkTuple `json:"html"`
				Avatar LinkTuple `json:"avatar"`
			} `json:"links"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"repository"`
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	CommentCount int           `json:"comment_count"`
	State        string        `json:"state"`
	TaskCount    int           `json:"task_count"`
	Participants []interface{} `json:"participants"`
	Reason       string        `json:"reason"`
	UpdatedOn    time.Time     `json:"updated_on"`
	Author       struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self   LinkTuple `json:"self"`
			HTML   LinkTuple `json:"html"`
			Avatar LinkTuple `json:"avatar"`
		} `json:"links"`
		Nickname  string `json:"nickname"`
		Type      string `json:"type"`
		AccountID string `json:"account_id"`
	} `json:"author"`
	MergeCommit interface{} `json:"merge_commit"`
	ClosedBy    interface{} `json:"closed_by"`
}

func GetPrFromBitBucket(baseUrl string, owner string, slug string, user string, pass string, idstring string) (*PullRequest, error) {

	req, err := http.NewRequest("GET", "https://"+baseUrl+"/repositories/"+owner+"/"+slug+"/pullrequests/"+idstring, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch PR json: %s", err)
	}
	req.SetBasicAuth(user, pass)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var output PullRequest

	err = json.Unmarshal(body, &output)
	return &output, err
}

type CommentValue struct {
	Links struct {
		Self LinkTuple `json:"self"`
		HTML LinkTuple `json:"html"`
	} `json:"links,omitempty"`
	Deleted     bool `json:"deleted"`
	Pullrequest struct {
		Type  string `json:"type"`
		ID    int    `json:"id"`
		Links struct {
			Self LinkTuple `json:"self"`
			HTML LinkTuple `json:"html"`
		} `json:"links"`
		Title string `json:"title"`
	} `json:"pullrequest"`
	Content   LinkTuple `json:"content"`
	CreatedOn time.Time `json:"created_on"`
	User      struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self   LinkTuple `json:"self"`
			HTML   LinkTuple `json:"html"`
			Avatar LinkTuple `json:"avatar"`
		} `json:"links"`
		Nickname  string `json:"nickname"`
		Type      string `json:"type"`
		AccountID string `json:"account_id"`
	} `json:"user"`
	UpdatedOn    time.Time `json:"updated_on"`
	Type         string    `json:"type"`
	ID           int       `json:"id"`
	CommentLinks struct {
		Self LinkTuple `json:"self"`
		Code LinkTuple `json:"code"`
		HTML LinkTuple `json:"html"`
	} `json:"links,omitempty"`
	Inline struct {
		To   int         `json:"to"`
		From interface{} `json:"from"`
		Path string      `json:"path"`
	} `json:"inline,omitempty"`
	Parent struct {
		ID    int `json:"id"`
		Links struct {
			Self LinkTuple `json:"self"`
			HTML LinkTuple `json:"html"`
		} `json:"links"`
	} `json:"parent,omitempty"`
	ParentComentLinks struct {
		Self LinkTuple `json:"self"`
		Code LinkTuple `json:"code"`
		HTML LinkTuple `json:"html"`
	} `json:"links,omitempty"`
}

type PRComments struct {
	Pagelen int            `json:"pagelen"`
	Next    string         `json:"next,omitempty"`
	Values  []CommentValue `json:"values,-"`
	Page    int            `json:"page"`
	Size    int            `json:"size"`
}

func GetAllPrCommentsFromBitBucket(baseurl, owner, slug, user, pass, idstring string, cb func(prj *PRComments, data interface{}), mydata interface{}) error {

	creq, cerr := http.NewRequest("GET", "https://"+baseurl+"/repositories/"+owner+"/"+slug+"/pullrequests/"+idstring+"/comments", nil)
	if cerr != nil {
		return fmt.Errorf("Unable to fetch PR json: %s", cerr)
	}
	creq.SetBasicAuth(user, pass)
	cresp, crerr := http.DefaultClient.Do(creq)
	if crerr != nil {
		return crerr
	}

	cbody, _ := ioutil.ReadAll(cresp.Body)

	var output PRComments
	err := json.Unmarshal(cbody, &output)
	if err != nil {
		return err
	}

	cresp.Body.Close()
	cb(&output, mydata)

	for output.Next != "" {
		creq, cerr := http.NewRequest("GET", output.Next, nil)
		if cerr != nil {
			return fmt.Errorf("Unable to fetch PR json: %s", cerr)
		}
		creq.SetBasicAuth(user, pass)
		cresp, crerr := http.DefaultClient.Do(creq)
		if crerr != nil {
			return crerr
		}

		output.Next = ""
		cbody, _ := ioutil.ReadAll(cresp.Body)
		cresp.Body.Close()
		err := json.Unmarshal(cbody, &output)
		if err != nil {
			return err
		}
		cb(&output, mydata)
	}
	return nil
}

type PRCommits struct {
	Pagelen int    `json:"pagelen"`
	Next    string `json:"next,omitempty"`
	Values  []struct {
		Hash       string `json:"hash"`
		Repository struct {
			Links struct {
				Self   LinkTuple `json:"self"`
				HTML   LinkTuple `json:"html"`
				Avatar LinkTuple `json:"avatar"`
			} `json:"links"`
			Type     string `json:"type"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"repository"`
		Links struct {
			Self     LinkTuple `json:"self"`
			Comments LinkTuple `json:"comments"`
			Patch    LinkTuple `json:"patch"`
			HTML     LinkTuple `json:"html"`
			Diff     LinkTuple `json:"diff"`
			Approve  LinkTuple `json:"approve"`
			Statuses LinkTuple `json:"statuses"`
		} `json:"links"`
		Author struct {
			Raw  string `json:"raw"`
			Type string `json:"type"`
			User struct {
				DisplayName string `json:"display_name"`
				UUID        string `json:"uuid"`
				Links       struct {
					Self   LinkTuple `json:"self"`
					HTML   LinkTuple `json:"html"`
					Avatar LinkTuple `json:"avatar"`
				} `json:"links"`
				Nickname  string `json:"nickname"`
				Type      string `json:"type"`
				AccountID string `json:"account_id"`
			} `json:"user"`
		} `json:"author"`
		Summary LinkTuple `json:"summary"`
		Parents []struct {
			Hash  string `json:"hash"`
			Type  string `json:"type"`
			Links struct {
				Self LinkTuple `json:"self"`
				HTML LinkTuple `json:"html"`
			} `json:"links"`
		} `json:"parents"`
		Date    time.Time `json:"date"`
		Message string    `json:"message"`
		Type    string    `json:"type"`
	} `json:"values"`
	Page int `json:"page"`
}

func GetAllPrCommitsFromBitBucket(baseurl, owner, slug, user, pass, idstring string, cb func(prc *PRCommits, data interface{}), mydata interface{}) error {

	creq, cerr := http.NewRequest("GET", "https://"+baseurl+"/repositories/"+owner+"/"+slug+"/pullrequests/"+idstring+"/commits", nil)
	if cerr != nil {
		return fmt.Errorf("Unable to fetch PR json: %s", cerr)
	}
	creq.SetBasicAuth(user, pass)
	cresp, crerr := http.DefaultClient.Do(creq)
	if crerr != nil {
		return crerr
	}
	defer cresp.Body.Close()

	cbody, _ := ioutil.ReadAll(cresp.Body)

	var output PRCommits

	err := json.Unmarshal(cbody, &output)
	if err != nil {
		return err
	}
	cb(&output, mydata)

	for output.Next != "" {
		creq, cerr := http.NewRequest("GET", output.Next, nil)
		if cerr != nil {
			return fmt.Errorf("Unable to fetch PR commits json: %s", cerr)
		}
		creq.SetBasicAuth(user, pass)
		cresp, crerr := http.DefaultClient.Do(creq)
		if crerr != nil {
			return crerr
		}

		output.Next = ""
		cbody, _ := ioutil.ReadAll(cresp.Body)
		cresp.Body.Close()
		err := json.Unmarshal(cbody, &output)
		if err != nil {
			return err
		}
		cb(&output, mydata)
	}

	return nil
}

type PrCommitComments struct {
	Pagelen int `json:"pagelen"`
	Values  []struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Code struct {
				Href string `json:"href"`
			} `json:"code"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
		Deleted bool `json:"deleted"`
		Commit  struct {
			Hash  string `json:"hash"`
			Type  string `json:"type"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
		} `json:"commit"`
		Content struct {
			Raw    string `json:"raw"`
			Markup string `json:"markup"`
			HTML   string `json:"html"`
			Type   string `json:"type"`
		} `json:"content"`
		CreatedOn time.Time `json:"created_on"`
		User      struct {
			DisplayName string `json:"display_name"`
			UUID        string `json:"uuid"`
			Links       struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
			Nickname  string `json:"nickname"`
			Type      string `json:"type"`
			AccountID string `json:"account_id"`
		} `json:"user"`
		Inline struct {
			To   int         `json:"to"`
			From interface{} `json:"from"`
			Path string      `json:"path"`
		} `json:"inline"`
		UpdatedOn time.Time `json:"updated_on"`
		Type      string    `json:"type"`
		ID        int       `json:"id"`
	} `json:"values"`
	Page int    `json:"page"`
	Size int    `json:"size"`
	Next string `json:"next, omitempty"`
}

func GetAllPrCommitCommentsFromBitBucket(url, user, pass string, cb func(prc *PrCommitComments, data interface{}), mydata interface{}) error {

	creq, cerr := http.NewRequest("GET", url, nil)
	if cerr != nil {
		return fmt.Errorf("Unable to fetch PR Commit Comments json: %s", cerr)
	}
	creq.SetBasicAuth(user, pass)
	cresp, crerr := http.DefaultClient.Do(creq)
	if crerr != nil {
		return crerr
	}
	defer cresp.Body.Close()

	cbody, _ := ioutil.ReadAll(cresp.Body)

	var output PrCommitComments

	err := json.Unmarshal(cbody, &output)
	if err != nil {
		return err
	}
	cb(&output, mydata)

	for output.Next != "" {
		creq, cerr := http.NewRequest("GET", output.Next, nil)
		if cerr != nil {
			return fmt.Errorf("Unable to fetch PR commits json: %s", cerr)
		}
		creq.SetBasicAuth(user, pass)
		cresp, crerr := http.DefaultClient.Do(creq)
		if crerr != nil {
			return crerr
		}

		output.Next = ""
		cbody, _ := ioutil.ReadAll(cresp.Body)
		cresp.Body.Close()
		err := json.Unmarshal(cbody, &output)
		if err != nil {
			return err
		}
		cb(&output, mydata)
	}

	return nil
}
