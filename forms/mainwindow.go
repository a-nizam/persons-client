package forms

import (
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	grpcclient "github.com/a-nizam/persons-client/grpc"
	"github.com/a-nizam/persons-client/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	a      fyne.App
	window fyne.Window
	menu   *fyne.MainMenu
	table  *widget.Table
	data   []models.Person
	curRow int
}

const (
	empty   int = 0
	Timeout int = 5
)

func NewMainWindow(a fyne.App) *MainWindow {
	w := MainWindow{a: a}
	w.window = w.a.NewWindow("Список сотрудников")
	w.window.Resize(fyne.NewSize(600, 400))

	quitMenuItem := fyne.NewMenuItem("Выход", nil)
	quitMenuItem.IsQuit = true

	w.menu = fyne.NewMainMenu(
		fyne.NewMenu("Операции",
			fyne.NewMenuItem("Обновить", w.Refresh),
			fyne.NewMenuItem("Добавить", w.AddPerson),
			fyne.NewMenuItem("Редактировать", w.EditPerson),
			fyne.NewMenuItem("Удалить", w.RemovePerson),
			quitMenuItem,
		),
	)
	w.window.SetMainMenu(w.menu)
	w.initTable()
	content := container.NewStack(w.table)
	w.window.SetContent(content)
	return &w
}

func (w *MainWindow) initTable() {
	w.table = widget.NewTableWithHeaders(
		func() (rows int, cols int) {
			return len(w.data), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(c widget.TableCellID, o fyne.CanvasObject) {
			if len(w.data) == empty {
				return
			}
			switch c.Col {
			case 0:
				o.(*widget.Label).SetText(strconv.FormatInt(w.data[c.Row].ID, 10))
			case 1:
				o.(*widget.Label).SetText(w.data[c.Row].Name)
			case 2:
				o.(*widget.Label).SetText(w.data[c.Row].Birthdate.Format("2006-01-02"))
			}

		},
	)
	w.table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	w.table.UpdateHeader = func(c widget.TableCellID, o fyne.CanvasObject) {
		switch c.Col {
		case 0:
			o.(*widget.Label).SetText("№")
		case 1:
			o.(*widget.Label).SetText("Имя")
		case 2:
			o.(*widget.Label).SetText("Дата рождения")
		}
	}
	w.table.SetColumnWidth(0, 50)
	w.table.SetColumnWidth(1, 250)
	w.table.SetColumnWidth(2, 200)

	w.table.OnSelected = func(id widget.TableCellID) {
		w.curRow = id.Row
	}

	go w.Refresh()
}

func (w *MainWindow) AddPerson() {
	personCardWindow := NewPersonCardWindow(w.a, 0, ActionCreate, w)
	personCardWindow.Show()
}

func (w *MainWindow) GetID() (id int64, err error) {
	if w.curRow >= 0 && w.curRow < len(w.data) {
		id = w.data[w.curRow].ID
	} else {
		err = errors.New("не выбрана запись в таблице")
	}
	return
}

func (w *MainWindow) EditPerson() {
	id, err := w.GetID()
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	personCardWindow := NewPersonCardWindow(w.a, id, ActionEdit, w)
	personCardWindow.Show()
}

func (w *MainWindow) RemovePerson() {
	id, err := w.GetID()
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}

	client := grpcclient.New()
	timeout, err := time.ParseDuration(strconv.Itoa(Timeout) + "s")
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.RemovePerson(ctx, id)
	if err != nil {
		dialog.ShowError(err, w.window)
	}
	dialog.ShowInformation("Упешно", "Сотрудник успешно удален", w.window)
	w.Refresh()
}

func (w *MainWindow) Refresh() {
	w.data = nil
	var err error
	client := grpcclient.New()
	timeout, err := time.ParseDuration(strconv.Itoa(Timeout) + "s")
	if err != nil {
		dialog.ShowError(err, w.window)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := client.GetList(ctx)
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		birthdate, err := time.Parse("2006-01-02", person.Birthdate)
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		w.data = append(w.data, models.Person{
			ID:        person.ID,
			Name:      person.Name,
			Birthdate: birthdate,
		})
		w.table.Refresh()
	}
}

func (w *MainWindow) Show() {
	w.window.Show()
}
