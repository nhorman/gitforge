package forgeuiview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponsePage struct {
	app         *tview.Application
	mainflex    *tview.Flex
	commentbox  *tview.TextView
	responsebox *tview.TextView
	name        string
	comment     string
}

func NewResponsePage(a *tview.Application) WindowPage {
	responsebox := tview.NewTextView()
	responsebox.SetBorder(true).SetTitle("Response")
	commentbox := tview.NewTextView()
	commentbox.SetBorder(true).SetTitle("Comment")
	mainflex := tview.NewFlex().SetDirection(tview.FlexRow)
	prflex := tview.NewFlex().SetDirection(tview.FlexRow)
	prflex.AddItem(commentbox, 0, 5, true)
	mainflex.AddItem(prflex, 0, 1, true)
	mainflex.AddItem(responsebox, 0, 1, true)

	responsepage := &ResponsePage{a, mainflex, commentbox, responsebox, "", ""}

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
			"Q - Quit"})
		PushPage("help")
		return nil
	default:
		return event
	}
	return event
}

func (m *ResponsePage) PagePreDisplay() {
	m.commentbox.SetText(m.comment)
	return
}

func (m *ResponsePage) PageDisplay() {
	return
}

func (m *ResponsePage) PagePostDisplay() {
	return
}

type PageResponseInfo struct {
	Comment string
}

func (m *ResponsePage) SetPageInfo(data interface{}) {
	m.comment = data.(*PageResponseInfo).Comment
	return
}
