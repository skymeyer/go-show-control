package io

type InputEvent struct {
	Action  Action
	Button  Button
	Control Control
}

func (i InputEvent) IsControl() bool {
	return i.Control > 0 && i.Button == 0
}

func (i InputEvent) IsControlPress() bool {
	return i.Control > 0 && i.Button == 0 && i.Action == ACTION_PRESS
}

func (i InputEvent) IsControlRelease() bool {
	return i.Control > 0 && i.Button == 0 && i.Action == ACTION_RELEASE
}

func (i InputEvent) IsButton() bool {
	return i.Button > 0 && i.Control == 0
}

func (i InputEvent) IsButtonPress() bool {
	return i.Button > 0 && i.Control == 0 && i.Action == ACTION_PRESS
}

func (i InputEvent) IsButtonRelease() bool {
	return i.Button > 0 && i.Control == 0 && i.Action == ACTION_RELEASE
}

type ControlChangeEvent struct {
	Mode    Mode
	Page    Page
	Grid    []GridControl
	Partial bool
}

type GridControl struct {
	Control     Control
	Description string
}

type ButtonChangeEvent struct {
	Grid    []GridButton
	Partial bool
}

type GridButton struct {
	Button      Button
	Kind        ButtonKind
	Description string
}
