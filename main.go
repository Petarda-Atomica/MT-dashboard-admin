package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var appLogo fyne.Resource
var globalDecryptionKey []byte

func main() {
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

	loginPage(w)
	w.ShowAndRun()
}
