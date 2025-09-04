package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func loginPage(w fyne.Window) {
	// Make logo image
	img := canvas.NewImageFromResource(appLogo)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(400, 400))

	// Make password field
	passwordField := widget.NewEntry()
	passwordField.Password = true

	// Make submit button
	submitButton := widget.NewButtonWithIcon("Login", theme.Icon(theme.IconNameMailSend), nil)

	// Attach the login function
	loginSubmit := func() {
		// Retrieve key
		b, err := retrieveUSBKey([]byte(passwordField.Text))
		if err != nil {
			log.Println("Error reading USB key:", err)
			loginPage(w)
			return
		}

		// Copy key to gloabl var
		globalDecryptionKey = make([]byte, len(b))
		copy(globalDecryptionKey, b)

		// Zero out memmory of b
		for i := range b {
			b[i] = 0
		}
	}
	passwordField.OnSubmitted = func(s string) { loginSubmit() }
	submitButton.OnTapped = loginSubmit

	// Build UI
	w.SetContent(container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),

		container.New(
			layout.NewHBoxLayout(),
			layout.NewSpacer(),
			img,
			layout.NewSpacer(),
		),

		layout.NewSpacer(),

		container.New(
			layout.NewFormLayout(),
			widget.NewLabel("Password:"),
			passwordField,
		),
		submitButton,
	))
}
