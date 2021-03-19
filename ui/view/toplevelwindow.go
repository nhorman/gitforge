package forgeuiview

import (
	"fmt"
	//"git-forge/ui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TopLevelControls interface {
	Quit()
}

var tlcontrols TopLevelControls = nil

type WindowPage interface {
	SetName(name string)
	GetName() string
	GetWindowPrimitive() tview.Primitive
	HandleInput(event *tcell.EventKey) *tcell.EventKey
	PagePreDisplay()
	PageDisplay()
	PagePostDisplay()
	SetPageInfo(interface{})
}

func inputEventHandler(event *tcell.EventKey) *tcell.EventKey {
	var pageret *tcell.EventKey = nil
	cpagename, _ := mainwindowpages.GetFrontPage()
	pageops, err := GetPage(cpagename)
	if err == nil {
		pageret = pageops.HandleInput(event)
	}
	if pageret != nil {
		// Global key ops
		runekey := event.Name()
		switch runekey {
		case "Rune[q]":
			tlcontrols.Quit()
			return nil
		default:
			return event
		}
	}
	return nil
}

func SwitchPage(newpage string, oldpage string) error {
	newpagewindow := pageregistry[newpage]
	newpagewindow.PagePreDisplay()
	mainwindowpages.SwitchToPage(newpage)
	newpagewindow.PageDisplay()
	if oldpage != "" {
		oldpagewindow := pageregistry[oldpage]
		oldpagewindow.PagePostDisplay()
	}
	return nil
}

var pagestack []WindowPage = make([]WindowPage, 0)
var mainwindowpages *tview.Pages = nil

func PushPage(newpage string) error {
	cpagename, _ := mainwindowpages.GetFrontPage()
	pagestack = append(pagestack, pageregistry[cpagename])
	return SwitchPage(newpage, cpagename)
}

func PopPage() (WindowPage, error) {
	cpagename, _ := mainwindowpages.GetFrontPage()
	if len(pagestack) == 0 {
		return nil, fmt.Errorf("Stack is Empty")
	}
	poppage := pagestack[len(pagestack)-1]
	pagestack = pagestack[0 : len(pagestack)-1]
	SwitchPage(poppage.GetName(), cpagename)
	return poppage, nil
}

func PeekPage() (WindowPage, error) {
	if len(pagestack) == 0 {
		return nil, fmt.Errorf("Stack Is Empty")
	}
	peekpage := pagestack[len(pagestack)-1]
	return peekpage, nil
}

var pageregistry map[string]WindowPage = make(map[string]WindowPage)

func RegisterPage(name string, window WindowPage, resize bool, visible bool) error {
	if _, ok := pageregistry[name]; ok {
		return fmt.Errorf("Page already registered")
	}
	pageregistry[name] = window
	window.SetName(name)
	mainwindowpages.AddPage(name, window.GetWindowPrimitive(), resize, visible)
	return nil
}

func GetPage(name string) (WindowPage, error) {
	if p, ok := pageregistry[name]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("No such page %s\n", name)
}

func DisplayTopLevelWindow(a *tview.Application, c TopLevelControls) error {
	tlcontrols = c
	mainwindowpages = tview.NewPages()

	mainpage := NewMainPage(a)
	RegisterPage("main", mainpage, true, true)
	helppage := NewHelpPage()
	RegisterPage("help", helppage, true, false)

	prlistpage := NewPRListPage(a)
	RegisterPage("prlist", prlistpage, true, false)

	errorpage := NewErrorPage()
	RegisterPage("error", errorpage, true, false)

	reviewpage := NewPRReviewPage(a)
	RegisterPage("prreview", reviewpage, true, false)

	responsepage := NewResponsePage(a)
	RegisterPage("response", responsepage, true, false)

	a.SetInputCapture(inputEventHandler)
	a.SetRoot(mainwindowpages, true)
	a.Run()
	return nil
}
