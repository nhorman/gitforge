package forgeui

import (
	"git-forge/forge"
	"git-forge/ui/controller"
)

func RunUi(forge forge.ForgeUIModel) error {
	control, err2 := forgecontroller.NewForgeController(forge)
	if err2 != nil {
		return err2
	}

	return control.StartUi()
}
