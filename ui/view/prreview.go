package forgeuiview

import (
	"git-forge/forge"
	//"git-forge/ui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type PRReviewPage struct {
	discussions *tview.TreeView
	commits     *tview.List
	display     *tview.TextView
	topflex     *tview.Flex
	pr          *forge.PR
	app         *tview.Application
	name        string
}

type DiscussionId struct {
	c forge.Discussion
	m *PRReviewPage
}

var focusList []tview.Primitive = nil
var focusidx int = 0

func NewPRReviewPage(a *tview.Application) WindowPage {
	PRPage := PRReviewPage{}

	PRPage.topflex = tview.NewFlex().SetDirection(tview.FlexRow)
	toprow := tview.NewFlex().SetDirection(tview.FlexColumn)
	PRPage.topflex.AddItem(toprow, 0, 1, true)
	bottomrow := tview.NewFlex()
	PRPage.topflex.AddItem(bottomrow, 0, 3, true)
	PRPage.discussions = tview.NewTreeView()
	PRPage.discussions.Box.SetTitle("Discussions")
	PRPage.discussions.Box.SetBorder(true)
	PRPage.commits = tview.NewList()
	PRPage.commits.Box.SetTitle("Commits")
	PRPage.commits.Box.SetBorder(true)
	PRPage.display = tview.NewTextView()
	PRPage.display.Box.SetBorder(true)
	toprow.AddItem(PRPage.discussions, 0, 1, true)
	toprow.AddItem(PRPage.commits, 0, 1, true)
	bottomrow.AddItem(PRPage.display, 0, 1, true)
	PRPage.app = a
	focusList = []tview.Primitive{PRPage.discussions, PRPage.commits}

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

func (m *PRReviewPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	runekey := event.Name()
	switch runekey {
	case "Rune[h]":
		helpwindow, _ := GetPage("help")
		helpwindow.SetPageInfo([]string{"H - This window",
			"Tab - Move between Discussion and Commit Pane",
			"Q - Back up to main window"})
		PushPage("help")
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
		data.m.display.SetText(data.c.Content)
		return
	})
	nodemap[0] = troot
	for _, c := range m.pr.Discussions {
		if c.Type == forge.INLINE {
			continue
		}
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

func (m *PRReviewPage) PagePreDisplay() {
	m.display.Box.SetTitle("Discussions for PR " + strconv.FormatInt(m.pr.PrId, 10) + ": " + m.pr.Title)
	m.display.Clear()
	focusidx = 0
	m.app.SetFocus(focusList[focusidx])
	m.populateDiscussions()
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
