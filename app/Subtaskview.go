package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/jinzhu/gorm"
	"github.com/rivo/tview"
	"github.com/s0u1z/gomodoro/models"
)

type SubtaskView struct {
	*tview.Flex
	subtasks      []*models.SubTask
	subtasklits   *tview.List
	newsubtask    *tview.InputField
	activesubtask *models.SubTask

	//Init ModelDB

	sdb *models.SubTaskModelDb
}

func NewSubtaskView(dbs *gorm.DB) *SubtaskView {

	view := SubtaskView{
		Flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		subtasklits: tview.NewList().ShowSecondaryText(false),
		newsubtask:  ConstructInputField("New Sub Task [+]"),
		sdb:         models.NewSubTaskModel(dbs),
	}
	view.subtasklits.SetSelectedBackgroundColor(tcell.ColorDarkMagenta)
	view.newsubtask.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:

			view.AddSubTask()

		case tcell.KeyEsc:
			app.SetFocus(view.subtasklits)
		}
	})

	view.AddItem(view.subtasklits, 0, 1, true).AddItem(view.newsubtask, 1, 0, false)
	view.SetBorder(true).SetTitle("Subtask")

	return &view

}

func (t *SubtaskView) LoadSubTask(task *models.Task) {

	subtask, err := t.sdb.GetSubtaskByTask(task)
	if err != nil {
		log.Panic(err)
	}

	t.SetSubTask(subtask)

}

func (t *SubtaskView) Clear() {

	t.subtasklits.Clear()
	t.subtasks = nil
	t.activesubtask = nil

}

func (t *SubtaskView) AddSubTask() {

	name := t.newsubtask.GetText()

	activetask := taskview.GetActiveTask()

	subtask, err := t.sdb.AddSubtask(name, *activetask)

	if err != nil {
		log.Panic(err)
	}

	t.subtasks = append(t.subtasks, subtask)

	t.AddSubTaskToList(len(t.subtasks) - 1)

	t.newsubtask.SetText("")
}

func (t *SubtaskView) SetSubTask(sub []*models.SubTask) {
	t.Clear()

	t.subtasks = sub

	for _, subtask := range t.subtasks {

		t.subtasks = append(t.subtasks, subtask)
		t.subtasklits.AddItem(subtask.Name, "", 0, func() {})

	}
}

func (t *SubtaskView) AddSubTaskToList(i int) {
	t.subtasklits.AddItem(t.subtasks[i].Name, "", 0, func() {})
}

func (t *SubtaskView) HandleKeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:

		t.subtasklits.SetCurrentItem(t.subtasklits.GetCurrentItem() - 1)

		return nil
	case tcell.KeyDown:

		t.subtasklits.SetCurrentItem(t.subtasklits.GetCurrentItem() - 1)

		return nil

	case tcell.KeyCtrlS:
		app.SetFocus(t.newsubtask)
		return nil
	}
	return event
}
