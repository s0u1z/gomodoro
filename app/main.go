package main

import (
	"log"

	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/jinzhu/gorm"
	"github.com/rivo/tview"
	"github.com/s0u1z/gomodoro/models"
)

var (
	app           *tview.Application
	mainlayout    *tview.Grid
	contentlayout *tview.Flex
	taskview      *Taskview
	subtaskview   *SubtaskView
	taskdetails   *TaskdetailsView
)

func main() {

	app = tview.NewApplication()

	db, err := models.InitDB()
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	footer := tview.NewTextView().SetBorder(true).SetTitle("Navigation")

	timer := tview.NewFlex().SetDirection(tview.FlexRow)
	timer.AddItem(tview.NewTextView().SetBorder(true).SetTitle("Pomodoro"), 0, 1, false).AddItem(tview.NewTextView().SetBorder(true).SetTitle("Notes"), 0, 1, false)

	mainlayout = tview.NewGrid().SetRows(4, 0, 4).SetColumns(40, 0, 40).AddItem(titlebar(), 0, 0, 1, 3, 0, 10, false).AddItem(footer, 2, 0, 1, 3, 30, 100, false)
	mainlayout.AddItem(layoutmanager(db), 1, 0, 1, 1, 0, 90, false).AddItem(taskdetails, 1, 1, 1, 1, 0, 90, false).AddItem(timer, 1, 2, 1, 1, 0, 100, false)
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
			app.SetFocus(taskview.tasklist)
			return nil
		case 's':
			app.SetFocus(subtaskview.subtasklits)
			return nil
		}

		// Handle based on current focus
		switch {
		case taskview.tasklist.HasFocus():
			event = taskview.HandleKeys(event)
		case subtaskview.subtasklits.HasFocus():
			event = subtaskview.HandleKeys(event)

		}

		return event
	})
}

func layoutmanager(db *gorm.DB) *tview.Flex {
	taskdetails = NewTaskdetailsView()
	taskview = NewTaskView(db)
	subtaskview = NewSubtaskView(db)

	contentlayout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(taskview, 0, 1, false).AddItem(subtaskview, 0, 1, false)

	return contentlayout
}

func titlebar() *tview.Flex {
	titleText := tview.NewTextView().SetText("[lime::b]Gomodoro [::-]- Pomodoro Task Manager  [red::b]Version[lime::b] 0.0.1").SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	titlebar := tview.NewFlex().AddItem(titleText, 0, 1, false)

	return titlebar

}
