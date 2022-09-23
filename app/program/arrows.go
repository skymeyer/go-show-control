package program

import (
	"go.skymyer.dev/show-control/io"
)

type CallButton io.Button
type ArrowButton io.Button

var (
	Call1 = CallButton(io.BTN_CALL_1)
	Call2 = CallButton(io.BTN_CALL_2)
	Call3 = CallButton(io.BTN_CALL_3)
	Call4 = CallButton(io.BTN_CALL_4)
	Call5 = CallButton(io.BTN_CALL_5)
	Call6 = CallButton(io.BTN_CALL_6)
	Call7 = CallButton(io.BTN_CALL_7)

	// TODO: add support for simultaneous up/down and left/right OR step control up/down
	//ArrowUp    ArrowButton = ArrowButton(io.BTN_ARROW_UP)
	//ArrowDown  ArrowButton = ArrowButton(io.BTN_ARROW_DOWN)
	ArrowLeft  ArrowButton = ArrowButton(io.BTN_ARROW_LEFT)
	ArrowRight ArrowButton = ArrowButton(io.BTN_ARROW_RIGHT)
)

type HandleArrow func(in ArrowButton, step uint8)

type ArrowCall struct {
	Button   CallButton
	Handle   HandleArrow
	Step     uint8
	Selected bool
}

type ArrowsCallControl struct {
	Selectors  []*ArrowCall
	showArrows bool
}

func (a *ArrowsCallControl) Handle(btn io.Button) bool {
	if !a.showArrows {
		return false
	}
	aBtn := ArrowButton(btn)
	if aBtn != ArrowLeft && aBtn != ArrowRight {
		return false
	}
	for _, call := range a.Selectors {
		if call.Selected {
			if call.Handle != nil {
				call.Handle(aBtn, call.Step)
				return true
			}
		}
	}
	return false
}

func (a *ArrowsCallControl) addArrowsSelectors(grid []io.GridButton) []io.GridButton {
	a.showArrows = false
	for _, call := range a.Selectors {
		grid = append(grid, io.GridButton{Button: io.Button(call.Button), Kind: io.BTN_KIND_CALL})
	}
	return grid
}

func (a *ArrowsCallControl) handleArrowSelect(btn io.Button, c []io.ButtonChangeEvent) []io.ButtonChangeEvent {
	if btn != io.BTN_CALL_1 && btn != io.BTN_CALL_2 && btn != io.BTN_CALL_3 && btn != io.BTN_CALL_4 &&
		btn != io.BTN_CALL_5 && btn != io.BTN_CALL_6 && btn != io.BTN_CALL_7 {
		return c
	}

	a.showArrows = false
	for _, call := range a.Selectors {
		var kind io.ButtonKind
		if io.Button(call.Button) == btn {
			if !call.Selected {
				a.showArrows = true
				call.Selected = true
				kind = io.BTN_KIND_CALL_SELECT
			} else {
				call.Selected = false
				kind = io.BTN_KIND_CALL
			}
		} else {
			call.Selected = false
			kind = io.BTN_KIND_CALL
		}
		c = append(c, io.ButtonChangeEvent{
			Partial: true,
			Grid: []io.GridButton{{
				Button: io.Button(call.Button),
				Kind:   kind,
			}},
		})
	}

	var arrowKind io.ButtonKind
	if a.showArrows {
		arrowKind = io.BTN_KIND_ARROW
	} else {
		arrowKind = io.BTN_KIND_OFF
	}

	c = append(c, io.ButtonChangeEvent{
		Partial: true,
		Grid: []io.GridButton{
			{
				Button: io.BTN_ARROW_LEFT,
				Kind:   arrowKind,
			},
			{
				Button: io.BTN_ARROW_RIGHT,
				Kind:   arrowKind,
			},
		},
	})

	return c
}
