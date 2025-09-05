package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type mtApp struct {
	f    func(fyne.Window, func())
	name string
	icon fyne.Resource
}

func searchMTApps(apps []mtApp, query string) []mtApp {
	var results []mtApp
	query = strings.ToLower(query)

	for _, app := range apps {
		if strings.Contains(strings.ToLower(app.name), query) {
			results = append(results, app)
		}
	}
	return results
}

func appSelectPage(w fyne.Window, appList []mtApp) {
	// Make search bar
	searchBar := widget.NewEntry()
	searchBar.PlaceHolder = "Caută unelte..."

	// Make the container for all the apps
	appsContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 150)))

	// Build UI
	w.SetContent(container.New(
		layout.NewVBoxLayout(),
		searchBar,
		appsContainer,
	))

	// Create a function to update the apps displayed
	updateAppList := func(query string) {
		// Search for the apps
		apps := searchMTApps(appList, query)

		appsContainer.RemoveAll()
		for _, this := range apps {
			// Create icon
			emptyBTN := widget.NewButton("", func() { this.f(w, func() { appSelectPage(w, appList) }) })
			appIcon := canvas.NewImageFromResource(this.icon)
			appIcon.ScaleMode = canvas.ImageScaleSmooth
			appIcon.FillMode = canvas.ImageFillContain
			appIcon.SetMinSize(fyne.NewSize(100, 100))
			btn := container.NewStack(emptyBTN, appIcon)

			// Create app title
			appTitle := widget.NewLabel(this.name)
			appTitle.Alignment = fyne.TextAlignCenter
			appTitle.Truncation = fyne.TextTruncateEllipsis

			// Add app to container
			appsContainer.Add(
				container.New(
					layout.NewVBoxLayout(),
					btn,
					appTitle,
				),
			)
		}

		appsContainer.Refresh()
	}

	// Use function from above
	searchBar.OnChanged = updateAppList
	updateAppList("")
}
