package forms

import (
	"context"
	"strconv"
	"time"

	grpcclient "github.com/a-nizam/persons-client/grpc"
	"github.com/a-nizam/persons-client/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PersonCardWindow struct {
	a              fyne.App
	window         fyne.Window
	parent         Refreshable
	action         int
	id             int64
	lblName        *widget.Label
	entryName      *widget.Entry
	lblBirthdate   *widget.Label
	entryBirthdate *widget.Entry
	btnSave        *widget.Button
}

const (
	ActionEdit int = iota
	ActionCreate
)

type Refreshable interface {
	Refresh()
}

func NewPersonCardWindow(a fyne.App, id int64, action int, rw Refreshable) *PersonCardWindow {
	w := PersonCardWindow{
		a:              a,
		action:         action,
		id:             id,
		parent:         rw,
		window:         a.NewWindow("Карточка сотрудника"),
		lblName:        widget.NewLabel("Имя"),
		entryName:      widget.NewEntry(),
		lblBirthdate:   widget.NewLabel("Дата рождения"),
		entryBirthdate: widget.NewEntry(),
	}

	var buttonText string
	switch action {
	case ActionCreate:
		buttonText = "Создать"
	case ActionEdit:
		buttonText = "Сохранить"
	}
	w.btnSave = widget.NewButton(buttonText, w.Save)
	w.window.Resize(fyne.NewSize(400, 300))
	content := container.NewVBox(w.lblName, w.entryName, w.lblBirthdate, w.entryBirthdate, layout.NewSpacer(), w.btnSave)
	w.window.SetContent(content)
	if action == ActionEdit {
		w.initFields()
	}
	return &w
}

func (w *PersonCardWindow) initFields() {
	client := grpcclient.New()
	timeout, err := time.ParseDuration(strconv.Itoa(Timeout) + "s")
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	person, err := client.GetPerson(ctx, w.id)
	if err != nil {
		dialog.ShowError(err, w.window)
	}
	w.entryName.SetText(person.Name)
	w.entryBirthdate.SetText(person.Birthdate.Format("2006-01-02"))
}

func (w *PersonCardWindow) Show() {
	w.window.Show()
}

func (w *PersonCardWindow) Save() {
	client := grpcclient.New()
	timeout, err := time.ParseDuration(strconv.Itoa(Timeout) + "s")
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	birtdate, err := time.Parse("2006-01-02", w.entryBirthdate.Text)
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}

	switch w.action {
	case ActionCreate:
		_, err := client.AddPerson(ctx, &models.Person{
			ID:        w.id,
			Name:      w.entryName.Text,
			Birthdate: birtdate,
		})
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		w.parent.Refresh()
		dialog.ShowInformation("Инфо", "Сотрудник успешно добавлен", w.window)

	case ActionEdit:
		err := client.EditPerson(ctx, &models.Person{
			ID:        w.id,
			Name:      w.entryName.Text,
			Birthdate: birtdate,
		})
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		w.parent.Refresh()
		dialog.ShowInformation("Инфо", "Сотрудник успешно сохранен", w.window)
	}
}
