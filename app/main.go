package main

import (
	"log"
	"reflect"

	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/jinzhu/gorm"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
	"github.com/s0u1z/gomodoro/models"
)

var (
	app           *tview.Application
	mainlayout    *tview.Grid
	contentlayout *tview.Flex
	appview       *AppView
)

type DetailsUi struct {
	*tview.Flex
	details       *tview.TextView
	detailsEditor *femto.View
	editor        *tview.TextViewWriter
	button        *tview.Button
}

func main() {

	app = tview.NewApplication()

	db, err := models.InitDB()
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	detailsui := DetailsUi{
		Flex:          tview.NewFlex().SetDirection(tview.FlexRow),
		details:       tview.NewTextView(),
		detailsEditor: femto.NewView(femto.NewBufferFromString("", "")),
		editor:        &tview.TextViewWriter{},
	}
	detailsui.detailsEditor.SetRuntimeFiles(runtime.Files)
	detailsui.detailsEditor.SetColorscheme(femto.ParseColorscheme("blue,red"))
	detailsui.detailsEditor.SetBorder(true)
	detailsui.detailsEditor.SetBorderColor(tcell.ColorAliceBlue)
	detailsui.detailsEditor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvt() {
			return event
		}
		switch event.Key() {
		case tcell.KeyESC:
			detailsui.details.SetText(detailsui.detailsEditor.Buf.String())
			detailsui.AddItem(detailsui.button.SetLabel("Update"), 1, 0, false)
			detailsui.detailsEditor.End()
			detailsui.detailsEditor.SetBorderColor(tcell.ColorDarkGrey)
			app.SetFocus(detailsui.details)

			return nil
		}
		return event
	})

	detailsui.button = tview.NewButton("edit").SetSelectedFunc(func() {

		detailsui.AddItem(detailsui.detailsEditor, 0, 3, false).SetBorder(true).SetTitle("Description")
		app.SetFocus(detailsui.detailsEditor)
		detailsui.RemoveItem(detailsui.button)

	})

	detailsui.details.SetBorder(true).SetTitle("Task Details")

	detailsui.AddItem(detailsui.details, 0, 1, false).AddItem(detailsui.button, 1, 0, false)

	footer := tview.NewTextView().SetBorder(true).SetTitle("Navigation")

	timer := tview.NewFlex().SetDirection(tview.FlexRow)
	timer.AddItem(tview.NewTextView().SetBorder(true).SetTitle("Pomodoro"), 0, 1, false).AddItem(tview.NewTextView().SetBorder(true).SetTitle("Notes"), 0, 1, false)

	mainlayout = tview.NewGrid().SetRows(4, 0, 4).SetColumns(40, 0, 40).AddItem(titlebar(), 0, 0, 1, 3, 0, 10, false).AddItem(footer, 2, 0, 1, 3, 30, 100, false)
	mainlayout.AddItem(layout(db), 1, 0, 1, 1, 0, 90, false).AddItem(detailsui, 1, 1, 1, 1, 0, 90, false).AddItem(timer, 1, 2, 1, 1, 0, 100, false)
	setKeyboardShortcuts()

	if err := app.SetRoot(mainlayout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func setKeyboardShortcuts() *tview.Application {

	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvt() {
			return event
		}

		// Global shortcuts
		switch unicode.ToLower(event.Rune()) {
		case 't':
			app.SetFocus(appview.tasklist)
			return nil
		case 's':
			app.SetFocus(appview.subtasklits)
			return nil
		}

		// Handle based on current focus
		switch {
		case appview.tasklist.HasFocus():
			event = appview.HandleShortcuts(event)
		case appview.subtasklits.HasFocus():
			event = appview.HandleShortcuts(event)

		}

		return event
	})
}

func ignoreKeyEvt() bool {
	textInputs := []string{"*tview.InputField", "*femto.View"}

	return InArray(reflect.TypeOf(app.GetFocus()).String(), textInputs)
}

func InArray(val interface{}, array interface{}) bool {
	return AtArrayPosition(val, array) != -1
}

func AtArrayPosition(val interface{}, array interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	}

	return
}

func layout(db *gorm.DB) *tview.Flex {

	appview = NewAppView(db)

	contentlayout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(appview, 0, 1, false)

	return contentlayout
}

func titlebar() *tview.Flex {
	titleText := tview.NewTextView().SetText("[lime::b]Gomodoro [::-]- Pomodoro Task Manager  [red::b]Version[lime::b] 0.0.1").SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	titlebar := tview.NewFlex().AddItem(titleText, 0, 1, false)

	return titlebar

}
