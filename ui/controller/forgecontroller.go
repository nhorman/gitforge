package forgecontroller

import (
	"git-forge/forge"
	"git-forge/ui/model"
	"git-forge/ui/view"
	"github.com/rivo/tview"
)

type ForgeControl interface {
	StartUi() error
}

type ForgeUiControl struct {
	app *tview.Application
}

func NewForgeController(forge forge.ForgeUIModel) (ForgeControl, error) {
	forgemodel.GetUiModel(forge)
	control := &ForgeUiControl{}
	control.app = tview.NewApplication()
	return control, nil
}

func (f *ForgeUiControl) StartUi() error {
	return forgeuiview.DisplayTopLevelWindow(f.app, f)
}

func (f *ForgeUiControl) Quit() {
	f.app.Stop()
}
