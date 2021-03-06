package forgeuiview

import (
	"git-forge/forge"
	"git-forge/ui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	//"gopkg.in/src-d/go-git.v4/plumbing/object"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type StatusItems struct {
	approved  *tview.TextView
	rows      *tview.Flex
	toprow    *tview.Flex
	bottomrow *tview.Flex
}

func NewStatusArea() *StatusItems {
	newstatus := StatusItems{}
	newstatus.rows = tview.NewFlex().SetDirection(tview.FlexRow)
	newstatus.toprow = tview.NewFlex().SetDirection(tview.FlexColumn)
	newstatus.bottomrow = tview.NewFlex().SetDirection(tview.FlexColumn)
	newstatus.rows.AddItem(newstatus.toprow, 0, 1, true)
	newstatus.rows.AddItem(newstatus.bottomrow, 0, 1, true)
	newstatus.approved = tview.NewTextView()
	newstatus.toprow.AddItem(newstatus.approved, 0, 1, true)
	return &newstatus
}

type PRReviewPage struct {
	toprow      *tview.Flex
	bottomrow   *tview.Flex
	statusrow   *tview.Flex
	discussions *tview.TreeView
	commits     *tview.TreeView
	tdisplay    *tview.TextView
	ldisplay    *tview.List
	displaylist bool
	selcomment  *forge.CommentData
	topflex     *tview.Flex
	statusarea  *StatusItems
	pr          *forge.PR
	app         *tview.Application
	name        string
}

type DiscussionId struct {
	c forge.CommentData
	m *PRReviewPage
}

var focusList []tview.Primitive = nil
var focusidx int = 0

func NewPRReviewPage(a *tview.Application) WindowPage {
	PRPage := PRReviewPage{}

	PRPage.topflex = tview.NewFlex().SetDirection(tview.FlexRow)
	statusrow := tview.NewFlex()
	statusrow.Box.SetBorder(true)
	statusrow.Box.SetTitle("Status")
	PRPage.topflex.AddItem(statusrow, 0, 1, true)
	PRPage.statusrow = statusrow
	toprow := tview.NewFlex().SetDirection(tview.FlexColumn)
	PRPage.topflex.AddItem(toprow, 0, 3, true)
	bottomrow := tview.NewFlex()
	PRPage.topflex.AddItem(bottomrow, 0, 5, true)
	PRPage.discussions = tview.NewTreeView()
	PRPage.discussions.Box.SetTitle("Discussions")
	PRPage.discussions.Box.SetBorder(true)
	PRPage.commits = tview.NewTreeView()
	PRPage.commits.Box.SetTitle("Commits")
	PRPage.commits.Box.SetBorder(true)
	PRPage.tdisplay = tview.NewTextView()
	PRPage.tdisplay.Box.SetBorder(true)
	PRPage.ldisplay = tview.NewList()
	PRPage.ldisplay.Box.SetBorder(true)
	PRPage.displaylist = false
	toprow.AddItem(PRPage.discussions, 0, 1, true)
	toprow.AddItem(PRPage.commits, 0, 1, true)
	bottomrow.AddItem(PRPage.tdisplay, 0, 1, true)
	PRPage.toprow = toprow
	PRPage.bottomrow = bottomrow
	PRPage.app = a
	PRPage.selcomment = nil
	PRPage.statusarea = NewStatusArea()
	PRPage.statusrow.AddItem(PRPage.statusarea.rows, 0, 1, true)
	focusList = []tview.Primitive{PRPage.discussions, PRPage.commits, bottomrow}

	return &PRPage
}

func (m *PRReviewPage) SetName(name string) {
	m.name = name
}

func (m *PRReviewPage) GetName() string {
	return m.name
}

func (m *PRReviewPage) GetWindowPrimitive() tview.Primitive {
	return m.topflex
}

func (m *PRReviewPage) HandleComment(newcomment bool) {
	var comment *os.File = nil
	var err error
	var commentname string = ""
	var oldcomment *forge.CommentData = m.selcomment
	var respcomment string = ""
	var pathlinestring string = ""
	if m.displaylist == true {
		itemidx := m.ldisplay.GetCurrentItem()
		var lines []string = make([]string, 0)
		for i := itemidx - 3; i > 0 && i < itemidx; i++ {
			newline, _ := m.ldisplay.GetItemText(i)
			lines = append(lines, newline)
		}
		respcomment = strings.Join(lines, "\n")
		_, pathlinestring = m.ldisplay.GetItemText(itemidx)
	} else {
		respcomment = m.tdisplay.GetText(true)
	}
	comment, err = ioutil.TempFile("", "GITFORGE")
	if err != nil {
		PopUpError(err)
	}
	defer os.Remove(comment.Name())
	comment.Write([]byte(respcomment))
	commentname = comment.Name()

	response, err := ioutil.TempFile("", "GITFORGERESPONSE")
	if err != nil {
		PopUpError(err)
	}
	defer os.Remove(response.Name())

	m.app.Suspend(func() {
		command := os.Getenv("GITFORGE_EDITOR")
		cmd := exec.Command(command, response.Name(), commentname)
		fmt.Printf("%s\n", cmd.String())
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		fmt.Printf("CMD COMPLETES\n")
		if err != nil {
			fmt.Printf("Failed to start editor: %s\n", err)
			fmt.Printf("press enter to return to review")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	})
	responseText, _ := ioutil.ReadFile(response.Name())
	response.Close()
	comment.Close()
	model, _ := forgemodel.GetUiModel(nil)
	if newcomment == true {
		oldcomment = nil
	}
	newcommentdata := &forge.CommentData{}
	newcommentdata.Type = forge.GENERAL
	newcommentdata.Content = string(responseText)
	if m.displaylist == true {
		newcommentdata.Type = forge.COMMIT
		stringinfo := strings.Split(pathlinestring, ":")
		newcommentdata.Commit = stringinfo[0]
		newcommentdata.Path = stringinfo[1]
		newcommentdata.Offset, _ = strconv.Atoi(stringinfo[2])
	}

	//TODO: Determine New comment type here based on oldcomment type?
	ret := model.PostComment(m.pr, oldcomment, newcommentdata)
	if ret != nil {
		PopUpError(ret)
	}
}

func (m *PRReviewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	runekey := event.Name()
	switch runekey {
	case "Rune[a]":
		model, _ := forgemodel.GetUiModel(nil)
		model.ApprovePR(m.pr)
	case "Rune[h]":
		helpwindow, _ := GetPage("help")
		helpwindow.SetPageInfo([]string{"H - This window",
			"Tab - Move between Discussion/Commit/Display Pane",
			"A - Toggle Approval status for this PR",
			"C - Start a new comment thread",
			"R - Respond To selected Comment",
			"Q - Back up to main window"})
		PushPage("help")
		return nil
	case "Rune[c]":
		m.HandleComment(true)
	case "Rune[r]":
		m.HandleComment(false)
		return nil
	case "Rune[q]":
		PopPage()
		return nil
	case "Tab":
		focusidx = (focusidx + 1) % len(focusList)
		m.app.SetFocus(focusList[focusidx])
		return nil
	}

	return event
}

func (m *PRReviewPage) populateDiscussions() {
	var nodemap map[int]*tview.TreeNode = make(map[int]*tview.TreeNode, 0)
	troot := tview.NewTreeNode("Discussions")
	var parent *tview.TreeNode
	var ok bool
	var current *tview.TreeNode = nil

	m.discussions.SetRoot(troot)
	m.discussions.SetTopLevel(1)
	m.discussions.SetSelectedFunc(func(node *tview.TreeNode) {
		data := node.GetReference().(*DiscussionId)
		data.m.bottomrow.Clear()
		data.m.bottomrow.AddItem(data.m.tdisplay, 0, 1, true)
		data.m.displaylist = false
		if data.c.Type == forge.GENERAL {
			data.m.tdisplay.SetRegions(false)
			data.m.tdisplay.SetText(data.c.Content)
			data.m.selcomment = &data.c
		} else if data.c.Type == forge.INLINE {
			data.m.tdisplay.SetRegions(true)
			model, _ := forgemodel.GetUiModel(nil)
			content, _ := model.GetPrInlineContent(data.m.pr, &data.c)
			data.m.tdisplay.SetText(content)
			data.m.tdisplay.Highlight("comment")
			data.m.tdisplay.ScrollToHighlight()
			data.m.selcomment = &data.c
		}
		return
	})
	nodemap[0] = troot
	for _, c := range m.pr.Discussions {
		parent, ok = nodemap[c.ParentId]
		if ok == false {
			return
		}
		var contentlen int = len(c.Content)
		if contentlen > 80 {
			contentlen = 80
		}
		shortcontent := c.Content[0:contentlen]
		child := tview.NewTreeNode(c.Author + " : " + shortcontent).SetSelectable(true)
		child.SetReference(&DiscussionId{c, m})
		parent.AddChild(child)
		_, ok = nodemap[c.Id]
		if ok == false {
			nodemap[c.Id] = child
		}
		if current == nil {
			current = child
			m.discussions.SetCurrentNode(child)
		}
	}

}

type CommentThread struct {
	Data     *forge.CommentData
	Hash     string
	Parent   *CommentThread
	Children []*CommentThread
	Node     *tview.TreeNode
	HLID     string
	Content  string
	Offset   int
	Path     string
}

func buildCommitComentThreadStrings(root *CommentThread, thread int, level int, idx int, basetabs string) (string, error) {
	var comments []string = make([]string, 0)
	//Capture our content
	if root.Data != nil {
		localcontent := fmt.Sprintf("[\"comment%d.%d.%d\"]%s[\"\"]", thread, level, idx, root.Data.Content)
		root.HLID = fmt.Sprintf("comment%d.%d.%d", thread, level, idx)
		localtabbedcontent := basetabs + strings.Replace(localcontent, "\n", basetabs, -1) + "\n"
		content := strings.Split(localtabbedcontent, "\n")
		comments = append(comments, content...)
	}
	newbasetabs := basetabs + "\t"
	var newidx int = 1
	for _, kid := range root.Children {
		kidcontent, _ := buildCommitComentThreadStrings(kid, thread, level+1, newidx, newbasetabs)
		newidx = newidx + 1
		comments = append(comments, kidcontent)
	}

	return strings.Join(comments, "\n"), nil
}

func setContentForChildren(root *CommentThread, content string) {

	root.Content = content
	root.Node.SetReference(root)
	for _, kid := range root.Children {
		setContentForChildren(kid, content)
	}
}

func insertThreadIntoCommit(commit string, threadcontent string, t *CommentThread) string {
	var foundpath bool = false
	var foundregion bool = false

	commitlines := strings.Split(commit, "\n")
	var idx int = 0
	for _, cl := range commitlines {
		var offset int = 0
		var offsetsize int = 0
		words := strings.Split(cl, " ")
		if words[0] == "diff" && words[2][2:] == t.Path {
			foundpath = true
		}
		if foundpath == true && words[0] == "@@" {
			offsetbit := strings.Split(words[2], ",")
			offset, _ = strconv.Atoi(strings.Trim(offsetbit[0], "+"))
			offsetsize, _ = strconv.Atoi(offsetbit[1])
			if (t.Offset >= offset) && (t.Offset <= offset+offsetsize) {
				foundregion = true
			}
		}

		if foundpath == true && foundregion == true {
			delta := t.Offset - offset
			idx = idx + delta + 1
			break
		}
		idx = idx + 1
	}

	outputcontent := make([]string, 0)
	outputcontent = append(outputcontent, commitlines[0:idx]...)
	outputcontent = append(outputcontent, threadcontent)
	if idx < len(commitlines) {
		outputcontent = append(outputcontent, commitlines[idx+1:]...)
	}

	return strings.Join(outputcontent, "\n")

}

func (m *PRReviewPage) populateCommitComments(child *tview.TreeNode, c *forge.Commit, allcommits []forge.CommentData) {
	var nodemap map[int]*CommentThread = make(map[int]*CommentThread, 0)
	var ids []int = make([]int, 0)

	model, _ := forgemodel.GetUiModel(nil)
	nodemap[0] = &CommentThread{nil, c.Hash, nil, make([]*CommentThread, 0), child, "", "", 0, ""}

	root := nodemap[0]
	ids = append(ids, 0)

	//iterate over our id list and find threads at each level
	for len(ids) != 0 {
		n := len(ids) - 1
		id := ids[n]
		ids = ids[:n]
		for i := 0; i < len(allcommits); i++ {
			idx := &allcommits[i]
			if idx.ParentId == id {
				//this commit has the current id as a parent, so
				//its a child
				parent := nodemap[id]
				newchild := &CommentThread{idx, c.Hash, parent, make([]*CommentThread, 0), tview.NewTreeNode(idx.Content), "", "", idx.Offset, idx.Path}
				parent.Node.AddChild(newchild.Node)
				parent.Children = append(parent.Children, newchild)
				ids = append(ids, idx.Id)
				nodemap[idx.Id] = newchild
			}
		}
	}

	root.Content = model.GetCommitData(c.Hash)
	root.HLID = ""
	root.Node.SetReference(root)
	// now we have a tree of comments for a given commit
	// now bundle each path down the tree into its own slice of strings
	var j int = 1
	for _, t := range root.Children {
		commentthread, _ := buildCommitComentThreadStrings(t, j, 1, 1, "\n\t")
		threadcontent := insertThreadIntoCommit(root.Content, commentthread, t)
		setContentForChildren(t, threadcontent)
		j = j + 1
	}
}

func (m *PRReviewPage) populateCommits() {
	model, _ := forgemodel.GetUiModel(nil)
	troot := tview.NewTreeNode("Commits")
	var first *tview.TreeNode = nil
	parent := troot
	m.commits.SetRoot(troot)
	m.commits.SetTopLevel(1)

	m.commits.SetSelectedFunc(func(node *tview.TreeNode) {
		data := node.GetReference().(*CommentThread)
		if data.HLID != "" {
			// existing comments display with the textview
			//we don't need to worry about line numbers
			m.bottomrow.Clear()
			m.bottomrow.AddItem(m.tdisplay, 0, 1, true)
			m.tdisplay.Clear()
			m.displaylist = false
			m.tdisplay.SetRegions(false)
			m.tdisplay.SetText(data.Content)
			m.selcomment = data.Data
			m.tdisplay.SetRegions(true)
			m.tdisplay.Highlight(data.HLID)
			m.tdisplay.ScrollToHighlight()
		} else {
			//commits themselves are list views so
			//we can select individual lines
			m.bottomrow.Clear()
			m.bottomrow.AddItem(m.ldisplay, 0, 1, true)
			m.ldisplay.Clear()
			m.displaylist = true
			contentlines := strings.Split(data.Content, "\n")
			var nextnum int = -1
			var currentnum int = -1
			var commitpath string = ""
			var linestate int = 0
			for _, l := range contentlines {
				if strings.HasPrefix(l, "diff --git") == true {
					diffline := strings.Split(l, " ")
					commitpath = strings.TrimLeft(diffline[3], "b/")
					currentnum = -1
					linestate = 0
				}

				if strings.HasPrefix(l, "@@") == true {
					//This is a diff line, so we can update
					//our line number
					litems := strings.Split(l, " ")
					var loffsetgroup []string
					if strings.HasPrefix(l, "@@@") == true {
						loffsetgroup = strings.Split(litems[3], ",")
					} else {
						loffsetgroup = strings.Split(litems[2], ",")
					}
					loffset := strings.TrimLeft(loffsetgroup[0], "+")
					nextnum, _ = strconv.Atoi(loffset)
					currentnum = -1
					linestate = 1
				}
				if linestate > 1 {
					m.ldisplay.AddItem("("+strconv.Itoa(currentnum)+") "+l, data.Hash+":"+commitpath+":"+strconv.Itoa(currentnum), 0, nil)
					m.ldisplay.ShowSecondaryText(false)
				} else {
					numstring := strconv.Itoa(nextnum)
					space := strings.Repeat(" ", len(numstring))
					m.ldisplay.AddItem(space+l, "", 0, nil)
					m.ldisplay.ShowSecondaryText(false)
				}
				currentnum = nextnum
				nextnum = nextnum + 1
				if linestate == 1 {
					linestate = 2
				}
			}
		}
	})

	for _, c := range m.pr.Commits {
		var line string = c.Hash
		commit, cerr := model.GetCommit(c.Hash)
		if cerr == nil {
			title := strings.Split(commit.Message, "\n")
			line = line + " - " + title[0]
		}
		child := tview.NewTreeNode(line).SetSelectable(true)

		m.populateCommitComments(child, &c, c.Comments)
		parent.AddChild(child)
		if first == nil {
			first = child
		}
	}
	m.commits.SetCurrentNode(first)
}

func (m *PRReviewPage) PagePreDisplay() {
	m.tdisplay.Box.SetTitle("Discussions for PR " + strconv.FormatInt(m.pr.PrId, 10) + ": " + m.pr.Title)
	m.tdisplay.Clear()
	focusidx = 0
	m.app.SetFocus(focusList[focusidx])
	m.populateDiscussions()
	m.populateCommits()
	var status string = "Approval Status: "
	switch m.pr.Approved {
	case forge.UNKNOWN:
		status = status + "UNKNOWN"
	case forge.UNAPPROVED:
		status = status + "NOT APPROVED"
	case forge.APPROVED:
		status = status + "APPROVED"
	}
	m.statusarea.approved.SetText(status)
	return
}

func (m *PRReviewPage) PageDisplay() {
	return
}

func (m *PRReviewPage) PagePostDisplay() {
	return
}

func (m *PRReviewPage) SetPageInfo(info interface{}) {
	m.pr = info.(*forge.PR)
}
