package forgeuiview

import (
	"git-forge/ui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

type PRListPage struct {
	app      *tview.Application
	topflex  *tview.Flex
	listview *tview.List
	name     string
}

func NewPRListPage(a *tview.Application) WindowPage {
	topflexitem := tview.NewFlex().SetDirection(tview.FlexRow)

	topflexitem.AddItem(tview.NewFlex(), 0, 1, false)
	centerrow := tview.NewFlex()
	topflexitem.AddItem(centerrow, 0, 1, false)
	topflexitem.AddItem(tview.NewFlex(), 0, 1, false)

	centerrow.SetDirection(tview.FlexColumn)
	centerrow.AddItem(tview.NewFlex(), 0, 1, false)
	listflex := tview.NewFlex()
	centerrow.AddItem(listflex, 0, 3, false)
	centerrow.AddItem(tview.NewFlex(), 0, 1, false)

	listView := tview.NewList()
	listflex.AddItem(listView, 0, 1, true)
	listView.Box.SetBorder(true)
	listView.Box.SetTitle("PR List")
	listView.ShowSecondaryText(false)
	return &PRListPage{
		app:      a,
		topflex:  topflexitem,
		listview: listView,
		name:     "",
	}
}

func (m *PRListPage) SetName(name string) {
	m.name = name
}

func (m *PRListPage) GetName() string {
	return m.name
}

func (m *PRListPage) GetWindowPrimitive() tview.Primitive {
	return m.topflex
}

func (m *PRListPage) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	runekey := event.Name()
	switch runekey {
	case "Rune[h]":
		helpwindow, _ := GetPage("help")
		helpwindow.SetPageInfo([]string{"H - This window",
			"W - Toggle watch status of a PR"})
		PushPage("help")
		return nil
	case "Rune[q]":
		PopPage()
		return nil
	case "Rune[w]":
		var werr error
		var add bool = true
		//Do work to add page to watchlist here
		itemtext, pridstring := m.listview.GetItemText(m.listview.GetCurrentItem())
		if strings.HasPrefix(itemtext, "* ") == true {
			add = false
			m.listview.SetItemText(m.listview.GetCurrentItem(), strings.TrimPrefix(itemtext, "* "), pridstring)
		} else {
			add = true
			m.listview.SetItemText(m.listview.GetCurrentItem(), "* "+itemtext, pridstring)
		}
		model, _ := forgemodel.GetUiModel(nil)
		if add == true {
			werr = model.AddWatchPr(pridstring)
		} else {
			werr = model.DelWatchPr(pridstring)
		}
		if werr != nil {
			PopUpError(werr)
		}
		return nil
	default:
		return event
	}

	return event
}

func (m *PRListPage) PagePreDisplay() {
	var titletext string
	m.listview.Clear()
	model, _ := forgemodel.GetUiModel(nil)

	prs, err2 := model.GetAllPrTitles()
	if err2 == nil {
		for _, r := range prs {
			titletext = r.Title
			watched, _ := model.PrIsWatched(strconv.FormatInt(r.PrId, 10))
			if watched == true {
				titletext = "* " + r.Title
			}
			m.listview.AddItem(titletext, strconv.FormatInt(r.PrId, 10), 0, nil)
		}
	}
	return
}

func (m *PRListPage) PageDisplay() {
	m.app.SetFocus(m.listview)
	return
}

func (m *PRListPage) PagePostDisplay() {
	return
}

func (m *PRListPage) SetPageInfo(info interface{}) {
}
