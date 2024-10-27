package main

import (
	"github.com/a-nizam/persons-client/forms"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("github.com/a-nizam/persons")
	window := forms.NewMainWindow(a)
	window.Show()
	a.Run()
}
