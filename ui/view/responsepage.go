package forgeuiview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PageResponseInfo struct {
	Comment string
	HLID    string
}

type ResponsePage struct {
	app          *tview.Application
	mainflex     *tview.Flex
	commentbox   *tview.TextView
	responsebox  *tview.TextView
	commitbutton *tview.Button
	cancelbutton *tview.Button
	name         string
	comment      *PageResponseInfo
}

var respfocuslist []tview.Primitive = nil
var respfocusidx int = 0

func NewResponsePage(a *tview.Application) WindowPage {
	responsebox := tview.NewTextView()
	responsebox.SetBorder(true).SetTitle("Response")
	commentbox := tview.NewTextView()
	commentbox.SetBorder(true).SetTitle("Comment")
	mainflex := tview.NewFlex().SetDirection(tview.FlexRow)
	prflex := tview.NewFlex().SetDirection(tview.FlexRow)
	prflex.AddItem(commentbox, 0, 5, true)
	mainflex.AddItem(prflex, 0, 5, true)
	mainflex.AddItem(responsebox, 0, 5, true)

	buttonflex := tview.NewFlex().SetDirection(tview.FlexColumn)
	commitbutton := tview.NewButton("Post")
	cancelbutton := tview.NewButton("Cancel")
	commitbutton.Box.SetBorder(true)
	cancelbutton.Box.SetBorder(true)
	buttonflex.AddItem(commitbutton, 0, 1, true)
	buttonflex.AddItem(cancelbutton, 0, 1, true)
	mainflex.AddItem(buttonflex, 0, 1, true)
	responsepage := &ResponsePage{a, mainflex, commentbox, responsebox, commitbutton, cancelbutton, "", nil}

	respfocuslist = []tview.Primitive{responsebox, commitbutton, cancelbutton, commentbox}

	return responsepage
}

func (m *ResponsePage) SetName(name string) {
	m.name = name
}

func (m *ResponsePage) GetName() string {
	return m.name
}

func (m *ResponsePage) GetWindowPrimitive() tview.Primitive {
	return m.mainflex
}

func (m *ResponsePage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	runekey := event.Name()
	switch runekey {
	case "Rune[h]":
		helpwindow, _ := GetPage("help")
		helpwindow.SetPageInfo([]string{"H - This window",
			"Q - Exit Response editor"})
		PushPage("help")
		return nil
	case "Rune[q]":
		PopPage()
		return nil
	default:
		return event
	}
	return event
}

func (m *ResponsePage) PagePreDisplay() {
	m.commentbox.Clear()
	m.commentbox.SetText(m.comment.Comment)
	if m.comment.HLID != "" {
		m.commentbox.SetRegions(true)
		m.commentbox.Highlight(m.comment.HLID)
		m.commentbox.ScrollToHighlight()
	}
	return
}

func (m *ResponsePage) PageDisplay() {
	m.app.SetFocus(m.responsebox)
	return
}

func (m *ResponsePage) PagePostDisplay() {
	return
}

func (m *ResponsePage) SetPageInfo(data interface{}) {
	m.comment = data.(*PageResponseInfo)
	return
}
