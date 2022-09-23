package calibrate

import (
	"go.skymyer.dev/show-control/app/mode"
	"go.skymyer.dev/show-control/io"
)

func init() {
	mode.Register("CALIBRATE_MODE", NewHandler)
}

func NewHandler() (mode.Handler, error) {
	return &Handler{
		name: "Calibration mode",
		pages: map[io.Page]string{
			io.PAGE_1: "Page 1",
			io.PAGE_2: "Page 2",
			io.PAGE_3: "Page 3",
		},
		pageButtons: map[io.Page][]io.GridButton{},
	}, nil
}

type Handler struct {
	name        string
	pages       map[io.Page]string
	pageButtons map[io.Page][]io.GridButton
}

func (h *Handler) Start() error {
	return nil
}

func (h *Handler) Stop() error {
	return nil
}

func (h *Handler) GetName() string {
	return h.name
}

func (h *Handler) GetPages() map[io.Page]string {
	return h.pages
}

func (h *Handler) SelectPage(page io.Page) (changes []io.ButtonChangeEvent) {
	return changes
}

func (h *Handler) HandleInput(in io.InputEvent) (changes []io.ButtonChangeEvent) {
	return changes
}
