package forgeuiview

import ()

func PopUpError(err error) {
	errorwindow, _ := GetPage("error")

	var msg []string = make([]string, 0)
	msg = append(msg, err.Error())
	errorwindow.SetPageInfo(msg)
	PushPage("error")
}
