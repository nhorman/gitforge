package forgeuiview

import (
	"git-forge/ui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type MainPage struct {
	app      *tview.Application
	mainflex *tview.Flex
	issuebox *tview.Box
	prbox    *tview.List
	name     string
}

func NewMainPage(a *tview.Application) WindowPage {
	prbox := tview.NewList()
	prbox.SetChangedFunc(func(index int, maintext string, secondtext string, shortcut rune) {
		if index == 0 {
			if prbox.GetItemCount() > 1 {
				prbox.SetCurrentItem(1)
			}
		}
	})

	prbox.SetSelectedFunc(func(index int, maintext string, secondtext string, shortcut rune) {
		model, _ := forgemodel.GetUiModel(nil)
		pr, err := model.GetLocalPr(secondtext)
		if err != nil {
			PopUpError(err)
			return
		}
		rpage, _ := GetPage("prreview")
		rpage.SetPageInfo(pr)
		PushPage("prreview")
		return
	})
	prbox.SetBorder(true).SetTitle("Git Forge Watched PR Reviews")
	issuebox := tview.NewBox().SetBorder(true).
		SetTitle("Git Forge Watched Issues")
	mainflex := tview.NewFlex().SetDirection(tview.FlexRow)
	prflex := tview.NewFlex().SetDirection(tview.FlexRow)
	prflex.AddItem(prbox, 0, 5, true)
	mainflex.AddItem(prflex, 0, 1, true)
	mainflex.AddItem(issuebox, 0, 1, true)

	mainpage := &MainPage{a, mainflex, issuebox, prbox, ""}

	mainpage.PagePreDisplay()
	return mainpage
}

func (m *MainPage) SetName(name string) {
	m.name = name
}

func (m *MainPage) GetName() string {
	return m.name
}

func (m *MainPage) GetWindowPrimitive() tview.Primitive {
	return m.mainflex
}

func (m *MainPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	runekey := event.Name()
	switch runekey {
	case "Rune[h]":
		helpwindow, _ := GetPage("help")
		helpwindow.SetPageInfo([]string{"H - This window",
			"P - List all PRs for this project",
			"Q - Quit"})
		PushPage("help")
		return nil
	case "Rune[p]":
		PushPage("prlist")
		return nil
	default:
		return event
	}
	return event
}

func (m *MainPage) PagePreDisplay() {
	m.prbox.ShowSecondaryText(false)
	m.prbox.Clear()
	model, _ := forgemodel.GetUiModel(nil)
	prs, _ := model.GetWatchedPrs()

	m.prbox.AddItem("PR                TITLE", "-1", 0, nil)
	for _, p := range prs {
		m.prbox.AddItem(strconv.FormatInt(p.PrId, 10)+"                "+p.Title, strconv.FormatInt(p.PrId, 10), 0, nil)
	}
	if len(prs) > 0 {
		m.prbox.SetCurrentItem(1)
	}

	m.app.SetFocus(m.prbox)
	return
}

func (m *MainPage) PageDisplay() {
	return
}

func (m *MainPage) PagePostDisplay() {
	return
}

func (m *MainPage) SetPageInfo(interface{}) {
	return
}
