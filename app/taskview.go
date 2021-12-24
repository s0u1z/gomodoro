package main

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jinzhu/gorm"
	"github.com/rivo/tview"
	"github.com/s0u1z/gomodoro/models"
)

type Taskview struct {
	*tview.Flex

	tasks    []*models.Task
	tasklist *tview.List
	newtask  *tview.InputField

	activetask *models.Task

	taskcount int

	db *models.TaskModelDb
}

func NewTaskView(dbs *gorm.DB) *Taskview {

	view := Taskview{

		Flex:     tview.NewFlex().SetDirection(tview.FlexRow),
		tasklist: tview.NewList().ShowSecondaryText(false),
		newtask:  ConstructInputField("New Task [+]"),
		db:       models.NewTaskModel(dbs),
	}
	view.tasklist.SetSelectedBackgroundColor(tcell.ColorOrangeRed)
	view.newtask.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:

			view.AddNewTask()

		case tcell.KeyEsc:
			app.SetFocus(view.tasklist)
		}
	})

	view.AddItem(view.tasklist, 0, 1, true).AddItem(view.newtask, 1, 0, false)
	view.SetBorder(true).SetTitle("Task")

	view.LoadListItems(false)

	return &view
}

func (t *Taskview) AddNewTask() {

	name := t.newtask.GetText()

	_, err := t.db.AddTask(name)

	if err != nil {
		log.Panic(err)
	}

	t.tasks, err = t.db.AllTask()
	if err != nil {
		log.Panic(err)
	}

	t.AddToList(len(t.tasks)-1, true)

	t.newtask.SetText("")
}

func (t *Taskview) AddToList(i int, selectedItem bool) {

	if t.tasks[i].Name != "" {
		t.tasklist.AddItem("-"+t.tasks[i].Name, "", 0, func(idx int) func() {
			return func() {
				t.ActivateTask(idx)
			}
		}(i))

		if selectedItem {
			t.tasklist.SetCurrentItem(-1)
			t.ActivateTask(i)

		}
	}

}

func (t *Taskview) ActivateTask(i int) {

	t.activetask = t.tasks[i]

	subtaskview.LoadSubTask(t.activetask)

	app.SetFocus(t)
}

func (t *Taskview) addSection(name string) {
	t.tasklist.AddItem("[::d]"+name, "", 0, nil)
	t.tasklist.AddItem("[::d]"+strings.Repeat(string(tcell.RuneHLine), 25), "", 0, nil)
}

func (t *Taskview) AddTaskList() {

	t.addSection("List Of Tasks")
	var err error
	t.tasklist.Clear()
	t.taskcount = t.tasklist.GetItemCount()

	t.tasks, err = t.db.AllTask()
	if err != nil {
		log.Println(err)
	}

	for i := range t.tasks {

		t.AddToList(i, false)
	}
	t.tasklist.SetCurrentItem(1)

}

func (t *Taskview) LoadListItems(focus bool) {

	t.AddTaskList()

	if focus {
		app.SetFocus(t.tasklist)
	}
}

func (t *Taskview) GetActiveTask() *models.Task {

	return t.activetask
}

func (t *Taskview) HandleKeys(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:

		t.tasklist.SetCurrentItem(t.tasklist.GetCurrentItem() + 1)

		return nil
	case tcell.KeyDown:

		t.tasklist.SetCurrentItem(t.tasklist.GetCurrentItem() - 1)

		return nil
	case tcell.KeyCtrlN:
		app.SetFocus(t.newtask)

	}
	return event
}
