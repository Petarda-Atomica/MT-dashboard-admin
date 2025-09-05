package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var appLogo fyne.Resource
var globalDecryptionKey []byte
var appList []mtApp = []mtApp{
	{
		name: "Setări",
		icon: theme.Icon(theme.IconNameSettings),
		f:    APP_settings,
	},
}

func main() {
	// Make sure the key is zeroed out before quitting
	defer func() {
		for i := range globalDecryptionKey {
			globalDecryptionKey[i] = 0
		}
	}()

	// Load logo
	var err error
	appLogo, err = fyne.LoadResourceFromPath("logo.png")
	if err != nil {
		log.Panic(err)
	}

	a := app.New()
	w := a.NewWindow("MT Admin")
	w.SetIcon(appLogo)
	w.SetMaster()

	loginPage(w, "")
	w.ShowAndRun()
}
