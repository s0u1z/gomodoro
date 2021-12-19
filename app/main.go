package main

import (
	"reflect"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/s0u1z/go-todo/data"
)

var (
	app           *tview.Application
	mainlayout    *tview.Grid
	contentlayout *tview.Flex
	appview       *AppView
)

type DetailsUi struct {
	*tview.Flex
	details *tview.TextView
	editor  *tview.TextViewWriter
	button  *tview.Button
}

func main() {

	app = tview.NewApplication()
	detailsui := DetailsUi{
		Flex:    tview.NewFlex().SetDirection(tview.FlexRow),
		details: tview.NewTextView(),
		editor:  &tview.TextViewWriter{},
		button:  tview.NewButton("Add Details "),
	}
	data.Init()
	footer := tview.NewTextView().SetBorder(true).SetTitle("Navigation")

	timer := tview.NewFlex().SetDirection(tview.FlexRow)
	timer.AddItem(tview.NewTextView().SetBorder(true).SetTitle("Pomodoro"), 0, 1, false).AddItem(tview.NewTextView().SetBorder(true).SetTitle("Notes"), 0, 1, false)
	details := detailsui.AddItem(tview.NewBox().SetBorder(true).SetTitle("Description"), 0, 1, false)

	mainlayout = tview.NewGrid().SetRows(4, 0, 4).SetColumns(40, 0, 40).AddItem(titlebar(), 0, 0, 1, 3, 0, 10, false).AddItem(footer, 2, 0, 1, 3, 30, 100, false)
	mainlayout.AddItem(layout(), 1, 0, 1, 1, 0, 90, false).AddItem(details, 1, 1, 1, 1, 0, 90, false).AddItem(timer, 1, 2, 1, 1, 0, 100, false)
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

func layout() *tview.Flex {

	appview = NewAppView()

	contentlayout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(appview, 0, 1, false)

	return contentlayout
}

func titlebar() *tview.Flex {
	titleText := tview.NewTextView().SetText("[lime::b]Gomodoro [::-]- Pomodoro Task Manager  [red::b]Version[lime::b] 0.0.1").SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	titlebar := tview.NewFlex().AddItem(titleText, 0, 1, false)

	return titlebar

}
