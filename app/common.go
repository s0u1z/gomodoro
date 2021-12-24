package main

import (
	"reflect"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ConstructInputField(placeholder string) *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder(placeholder).
		SetPlaceholderTextColor(tcell.ColorDarkSlateBlue).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorLightBlue)
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
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				return
			}
		}
	}

	return
}
