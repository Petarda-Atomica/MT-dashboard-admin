package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func APP_settings(w fyne.Window, backf func()) {
	// Make separator
	separator := canvas.NewRectangle(theme.Color(theme.ColorNameDisabled))
	separator.SetMinSize(fyne.NewSize(1, 5))

	// Make change password header
	changePassHeader := canvas.NewText("Schimbă parola", theme.Color(theme.ColorNameForeground))
	changePassHeader.TextSize = theme.Size(theme.SizeNameHeadingText)

	// Make change password entry
	changePassEntry := widget.NewEntry()
	changePassEntry.Password = true

	// Make chnage password button
	changePassBtn := widget.NewButton("Confirmă", func() {
		// Check password is valid
		if changePassEntry.Validate() != nil {
			return
		}

		// Change password on USB
		writeBlob([]byte(changePassEntry.Text), globalDecryptionKey, USBBlobPath)

		// Clear password entry
		changePassEntry.Text = ""
		changePassEntry.Refresh()
	})

	// Build UI
	w.SetContent(container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewStackLayout(),
			canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
			container.New(
				layout.NewHBoxLayout(),
				widget.NewButtonWithIcon("Înapoi", theme.Icon(theme.IconNameNavigateBack), backf),
			),
		),

		changePassHeader,
		container.New(
			layout.NewFormLayout(),
			changePassBtn,
			changePassEntry,
		),
		separator,
	))
}
