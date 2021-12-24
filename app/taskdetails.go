package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

type TaskdetailsView struct {
	*tview.Flex
	details       *tview.TextView
	detailsEditor *femto.View
	editor        *tview.TextViewWriter
	button        *tview.Button
}

func NewTaskdetailsView() *TaskdetailsView {

	view := TaskdetailsView{
		Flex:          tview.NewFlex().SetDirection(tview.FlexRow),
		details:       tview.NewTextView(),
		detailsEditor: femto.NewView(femto.NewBufferFromString("", "")),
		editor:        &tview.TextViewWriter{},
	}
	view.detailsEditor.SetRuntimeFiles(runtime.Files)
	view.detailsEditor.SetColorscheme(femto.ParseColorscheme("blue,red"))
	view.detailsEditor.SetBorder(true)
	view.detailsEditor.SetBorderColor(tcell.ColorAliceBlue)

	view.detailsEditor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvt() {
			return event
		}
		switch event.Key() {
		case tcell.KeyESC:
			view.details.SetText(view.detailsEditor.Buf.String())
			view.AddItem(view.button.SetLabel("Update"), 1, 0, false)
			view.detailsEditor.End()
			view.detailsEditor.SetBorderColor(tcell.ColorDarkGrey)
			app.SetFocus(view.details)

			return nil
		}
		return event
	})

	view.button = tview.NewButton("edit").SetSelectedFunc(func() {

		view.AddItem(view.detailsEditor, 0, 3, false).SetBorder(true).SetTitle("Description")
		app.SetFocus(view.detailsEditor)
		view.RemoveItem(view.button)

	})

	view.details.SetBorder(true).SetTitle("Task Details")

	view.AddItem(view.details, 0, 1, false).AddItem(view.button, 1, 0, false)

	return &view
}
