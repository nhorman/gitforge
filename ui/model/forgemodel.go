package forgemodel

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
	//"gopkg.in/src-d/go-git.v4/config"
	"encoding/json"
	"git-forge/forge"
	"os/exec"
	"strconv"
	"strings"
)

type ForgeUiModel struct {
	Forge forge.ForgeUIModel
}

type ForgeUiOpts interface {
	GetAllPrTitles() ([]forge.PrTitle, error)
	AddWatchPr(idstring string) error
	DelWatchPr(idstring string) error
	PrIsWatched(idstring string) (bool, error)
	GetWatchedPrs() ([]string, []string, error)
	GetLocalPr(idstring string) (*forge.PR, error)
	GetPrInlineContent(*forge.PR, *forge.Discussion) (string, error)
	RefreshPr(*forge.PR, func(*forge.PR, forge.UpdateResult)) error
}

func NewUiModel(forge forge.ForgeUIModel) (*ForgeUiModel, error) {
	model := &ForgeUiModel{
		Forge: forge,
	}
	return model, nil
}

var _internalmodel *ForgeUiModel = nil

func GetUiModel(forge forge.ForgeUIModel) (*ForgeUiModel, error) {
	if _internalmodel == nil {
		_internalmodel, _ = NewUiModel(forge)
	}
	return _internalmodel, nil
}

func (f *ForgeUiModel) GetAllPrTitles() ([]forge.PrTitle, error) {

	titles, err := f.Forge.GetAllPrTitles()
	if err != nil {
		return nil, err
	}
	return titles, nil
}

func (f *ForgeUiModel) AddWatchPr(idstring string) error {
	// TODO SHOULD BE DOING FETCH AND NOTE ADD HERE
	// JUST NEED TO GET PULL INFO FROM FORGE DRIVER
	pr, err := f.Forge.GetPr(idstring)
	if err != nil {
		return err
	}
	return f.addWatchPr(pr)
}

func (f *ForgeUiModel) addWatchPr(pr *forge.PR) error {

	jsonout, jerr := json.Marshal(pr)
	if jerr != nil {
		return jerr
	}

	pullcmd := exec.Command("git", "fetch", pr.PullSpec.Source.URL, pr.PullSpec.Source.BranchName+":refs/prs/"+strconv.FormatInt(pr.PrId, 10))
	pullerr := pullcmd.Run()
	if pullerr != nil {
		return pullerr
	}
	notecmd := exec.Command("git", "notes", "add", "-m", string(jsonout), "refs/prs/"+strconv.FormatInt(pr.PrId, 10))
	noteerr := notecmd.Run()
	if noteerr != nil {
		unwindcmd := exec.Command("git", "update-ref", "-d", "refs/prs"+strconv.FormatInt(pr.PrId, 10))
		unwinderr := unwindcmd.Run()
		if unwinderr != nil {
			return unwinderr
		}
	}
	return nil
}

func (f *ForgeUiModel) DelWatchPr(idstring string) error {
	// This is forge indepent for now,so do all the work here
	cmd := exec.Command("git", "notes", "remove", "refs/prs/"+idstring)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd2 := exec.Command("git", "update-ref", "-d", "refs/prs/"+idstring)
	err2 := cmd2.Run()
	if err2 != nil {
		return err2
	}
	return nil
}

func (f *ForgeUiModel) PrIsWatched(idstring string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "refs/prs/"+idstring)
	err := cmd.Run()
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (f *ForgeUiModel) GetWatchedPrs() ([]forge.PrTitle, error) {
	//var titles []string = make([]string, 0)
	var prs []forge.PrTitle = make([]forge.PrTitle, 0)
	var ids []string = make([]string, 0)
	cfg := &forge.ForgeObj{}
	repo, err := cfg.OpenLocalRepo()
	if err != nil {
		return nil, err
	}

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().String()
		if strings.Contains(name, "refs/prs/") == true {
			ids = append(ids, strings.TrimPrefix(name, "refs/prs/"))
		}
		return nil
	})

	for _, id := range ids {
		pr, terr := f.GetLocalPr(id)
		if terr != nil {
			return nil, terr
		}
		prs = append(prs, forge.PrTitle{Title: pr.Title, PrId: pr.PrId})
	}

	return prs, nil
}

func (f *ForgeUiModel) GetLocalPr(idstring string) (*forge.PR, error) {
	var pr forge.PR
	cmd := exec.Command("git", "notes", "show", "refs/prs/"+idstring)

	jsonout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	err2 := json.Unmarshal(jsonout, &pr)
	if err2 != nil {
		return nil, err2
	}
	return &pr, nil
}

func (f *ForgeUiModel) GetPrInlineContent(pr *forge.PR, d *forge.Discussion) (string, error) {
	cmd := exec.Command("git", "show", "refs/prs/"+strconv.FormatInt(pr.PrId, 10)+":"+d.Inline.Path)
	content, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var newcontent []string = make([]string, 0)
	contentlines := strings.Split(string(content), "\n")
	newcontent = append(newcontent, contentlines[0:d.Inline.Offset-1]...)
	newcontent = append(newcontent, contentlines[d.Inline.Offset]+"[\"comment\"]")
	newcontent = append(newcontent, d.Content+"[\"\"]")
	newcontent = append(newcontent, contentlines[d.Inline.Offset+1:]...)

	ret := strings.Join(newcontent, "\n")
	return ret, nil
}

func (f *ForgeUiModel) RefreshPr(pr *forge.PR, complete func(pr *forge.PR, result *forge.UpdatedPR, data interface{}), data interface{}) error {
	schan, serr := f.Forge.RefreshPr(pr)
	if serr != nil {
		return serr
	}

	go func(pr *forge.PR) {
		retcode := <-schan
		complete(pr, retcode, data)
		switch retcode.Result {
		case forge.UPDATE_CURRENT:
			// Nothing to do
		case forge.UPDATE_REPULL:
			f.DelWatchPr(strconv.FormatInt(pr.PrId, 10))
			f.addWatchPr(retcode.Pr)
			retcode.Result = forge.UPDATE_CURRENT
			complete(pr, retcode, data)
		case forge.UPDATE_FINISHED:
			// Nothing to do here, we do this below
		case forge.UPDATE_FAILED:
			// Nothing to do here
		default:
		}
		retcode.Result = forge.UPDATE_FINISHED
		complete(pr, retcode, data)
	}(pr)

	return nil
}
