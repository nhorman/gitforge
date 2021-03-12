package bitbucketforge

import (
	"encoding/json"
	"time"
)

type LinkTuple struct {
	Raw    string `json:"raw,omitempty"`
	Markup string `json:"markup,omitempty"`
	HTML   string `json:"html,omitempty"`
	Href   string `json:"href,omitempty"`
	Type   string `json:"omitempty"`
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

func PrJsonToStruct(input []byte) (PullRequest, error) {
	var output PullRequest

	err := json.Unmarshal(input, &output)
	return output, err
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
	Values  []CommentValue `json:"values,-"`
	Page    int            `json:"page"`
	Size    int            `json:"size"`
}

func PrCommentsJsonToStruct(input []byte) (PRComments, error) {
	var output PRComments

	err := json.Unmarshal(input, &output)
	return output, err
}
