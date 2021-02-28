package forgeuiview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type ErrorPage struct {
	flexlayout  *tview.Flex
	errorbox    *tview.TextView
	errorstring string
	name        string
}

func NewErrorPage() WindowPage {
	flexl := tview.NewFlex().SetDirection(tview.FlexRow)
	errorpage := &ErrorPage{flexlayout: flexl}
	errorpage.errorstring = ""
	errorpage.name = ""
	errorpage.errorbox = tview.NewTextView()
	errorpage.errorbox.SetBorder(true)
	errorpage.errorbox.SetTitle("Error")
	errorpage.errorbox.SetTextAlign(tview.AlignCenter)
	errorpage.flexlayout.AddItem(errorpage.errorbox, 0, 1, true)

	return errorpage
}

func (m *ErrorPage) SetName(name string) {
	m.name = name
}

func (m *ErrorPage) GetName() string {
	return m.name
}

func (m *ErrorPage) GetWindowPrimitive() tview.Primitive {
	return m.flexlayout
}

func (m *ErrorPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Any input closes the help window
	PopPage()
	return event
}

func (m *ErrorPage) PagePreDisplay() {
	m.errorbox.SetText(m.errorstring)
	return
}

func (m *ErrorPage) PageDisplay() {
	return
}

func (m *ErrorPage) PagePostDisplay() {
	return
}

func (m *ErrorPage) SetPageInfo(info interface{}) {
	errortext := info.([]string)
	errortext = append(errortext, "Press any button to close error window")
	m.errorstring = strings.Join(errortext, "\n")
}
