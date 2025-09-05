package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

func loginPage(w fyne.Window, receivedError string) {
	// Make error message
	errMsg := container.New(
		layout.NewStackLayout())
	errMsgText := widget.NewLabel("")
	if receivedError != "" {
		errMsgText.Text = receivedError
		errMsg.Add(canvas.NewRectangle(colornames.Red))
		errMsg.Add(container.NewHScroll(errMsgText))
	}

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
			switch err.Error() {
			case "no valid blobs":
				loginPage(w, "Introdu un stick USB valid!")

			case "wrong password":
				loginPage(w, "Parola greșită!")
			}
			return
		}

		// Copy key to gloabl var
		globalDecryptionKey = make([]byte, len(b))
		copy(globalDecryptionKey, b)

		// Zero out memmory of b
		for i := range b {
			b[i] = 0
		}

		// Go to app select page
		log.Println("Successfully decrypted USB key!")
		appSelectPage(w, appList)
	}
	passwordField.OnSubmitted = func(s string) { loginSubmit() }
	submitButton.OnTapped = loginSubmit

	// Build UI
	w.SetContent(container.New(
		layout.NewVBoxLayout(),
		errMsg,

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
