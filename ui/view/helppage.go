package forgeuiview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type HelpPage struct {
	flexlayout *tview.Flex
	helpbox    *tview.TextView
	helpstring string
	name       string
}

func NewHelpPage() WindowPage {
	flexl := tview.NewFlex().SetDirection(tview.FlexRow)
	helppage := &HelpPage{flexlayout: flexl}
	helppage.helpstring = ""
	helppage.name = ""
	helppage.helpbox = tview.NewTextView()
	helppage.helpbox.SetBorder(true)
	helppage.helpbox.SetTitle("Help")
	helppage.helpbox.SetTextAlign(tview.AlignCenter)
	helppage.flexlayout.AddItem(helppage.helpbox, 0, 1, true)

	return helppage
}

func (m *HelpPage) SetName(name string) {
	m.name = name
}

func (m *HelpPage) GetName() string {
	return m.name
}

func (m *HelpPage) GetWindowPrimitive() tview.Primitive {
	return m.flexlayout
}

func (m *HelpPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	// Any input closes the help window
	PopPage()
	return nil
}

func (m *HelpPage) PagePreDisplay() {
	m.helpbox.SetText(m.helpstring)
	return
}

func (m *HelpPage) PageDisplay() {
	return
}

func (m *HelpPage) PagePostDisplay() {
	return
}

func (m *HelpPage) SetPageInfo(info interface{}) {
	helptext := info.([]string)
	helptext = append(helptext, "Press any button to close help")
	m.helpstring = strings.Join(helptext, "\n")
}
