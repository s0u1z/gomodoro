package main

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gomodoro/data"
	"github.com/rivo/tview"
)

var activetask *data.Task

type AppView struct {
	*tview.Flex

	tasks    data.Tasks
	tasklist *tview.List
	newtask  *tview.InputField

	taskcount int

	subtasks      data.SubTasks
	subtasklits   *tview.List
	newsubtask    *tview.InputField
	activesubtask *data.SubTask
}

func NewAppView() *AppView {

	view := AppView{

		Flex: tview.NewFlex().SetDirection(tview.FlexRow),

		tasklist: tview.NewList().ShowSecondaryText(false),
		newtask:  makeLightTextInput("New Task [+]"),

		subtasklits: tview.NewList().ShowSecondaryText(false),
		newsubtask:  makeLightTextInput("New Sub Task [+]"),
	}

	view.tasklist.SetSelectedBackgroundColor(tcell.ColorOrangeRed)

	view.subtasklits.SetSelectedBackgroundColor(tcell.ColorDarkMagenta)

	view.newtask.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:

			view.AddNewTask()

		case tcell.KeyEsc:
			app.SetFocus(view.tasklist)
		}
	})

	view.newsubtask.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:

			view.AddSubTask()

		case tcell.KeyEsc:
			app.SetFocus(view.subtasklits)
		}
	})

	taskview := tview.NewFlex().SetDirection(tview.FlexRow)
	taskview.AddItem(view.tasklist, 0, 1, true).AddItem(view.newtask, 1, 0, false)
	taskview.SetBorder(true).SetTitle("Task")

	subtaskview := tview.NewFlex().SetDirection(tview.FlexRow)
	subtaskview.AddItem(view.subtasklits, 0, 1, true).AddItem(view.newsubtask, 1, 0, false)
	subtaskview.SetBorder(true).SetTitle("Subtask")

	view.AddItem(taskview, 0, 1, false).AddItem(subtaskview, 0, 1, false)

	view.LoadListItems(false)

	return &view
}

func (t *AppView) AddNewTask() {

	name := t.newtask.GetText()

	task := data.CreateTask(name)

	err := data.AddTaskToList(*task)

	if err != nil {
		log.Panic(err)
	}

	t.tasks, err = data.GetAllTask()
	if err != nil {
		log.Panic(err)
	}

	t.AddToList(len(t.tasks)-1, true)

	t.newtask.SetText("")
}

func (t *AppView) AddToList(i int, selectedItem bool) {

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

func (t *AppView) ActivateTask(i int) {

	activetask = t.tasks[i]

	t.LoadSubTask(*activetask)

	app.SetFocus(t)
}

func (t *AppView) LoadSubTask(task data.Task) {

	subtask, err := data.GetSubTaskByName(task)
	if err != nil {
		log.Panic(err)
	}

	t.SetSubTask(subtask)

}

func (t *AppView) Clear() {

	t.subtasklits.Clear()
	t.subtasks = nil
	t.activesubtask = nil

}

func (t *AppView) AddSubTask() {

	name := t.newsubtask.GetText()

	activetask := activetask

	subtask := data.CreateSubTask(name)

	err := data.AddSubTask(*subtask, *activetask)
	if err != nil {
		log.Panic(err)
	}

	t.subtasks = append(t.subtasks, subtask)

	t.AddSubTaskToList(len(t.subtasks) - 1)

	t.newsubtask.SetText("")
}

func (t *AppView) SetSubTask(sub data.SubTasks) {
	t.Clear()

	t.subtasks = sub

	for _, subtask := range t.subtasks {

		t.subtasks = append(t.subtasks, subtask)
		t.subtasklits.AddItem(subtask.Name, "", 0, func() {})

	}
}

func (t *AppView) AddSubTaskToList(i int) {
	t.subtasklits.AddItem(t.subtasks[i].Name, "", 0, func() {})
}

func (t *AppView) HandleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		if t.tasklist.HasFocus() {
			t.tasklist.SetCurrentItem(t.tasklist.GetCurrentItem() + 1)
		}
		t.subtasklits.SetCurrentItem(t.subtasklits.GetCurrentItem() - 1)

		return nil
	case tcell.KeyDown:
		if t.tasklist.HasFocus() {
			t.tasklist.SetCurrentItem(t.tasklist.GetCurrentItem() - 1)
		}
		t.subtasklits.SetCurrentItem(t.subtasklits.GetCurrentItem() - 1)

		return nil
	case tcell.KeyCtrlN:
		app.SetFocus(t.newtask)
	case tcell.KeyCtrlS:
		app.SetFocus(t.newsubtask)
		return nil
	}
	return event
}

func (t *AppView) addSection(name string) {
	t.tasklist.AddItem("[::d]"+name, "", 0, nil)
	t.tasklist.AddItem("[::d]"+strings.Repeat(string(tcell.RuneHLine), 25), "", 0, nil)
}

func (t *AppView) AddTaskList() {

	t.addSection("List Of Tasks")
	var err error
	t.tasklist.Clear()
	t.subtasklits.Clear()

	t.taskcount = t.tasklist.GetItemCount()

	t.tasks, err = data.GetAllTask()
	if err != nil {
		log.Println(err)
	}

	for i := range t.tasks {

		t.AddToList(i, false)
	}
	t.tasklist.SetCurrentItem(1)

}

func (t *AppView) LoadListItems(focus bool) {

	t.AddTaskList()

	if focus {
		app.SetFocus(t.tasklist)
	}
}

func makeLightTextInput(placeholder string) *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder(placeholder).
		SetPlaceholderTextColor(tcell.ColorDarkSlateBlue).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorLightBlue)
}
