package mode

import (
	"go.skymyer.dev/show-control/io"
)

func init() {
	Register("DUMMY_MODE", NewDummyHandler)
}

func NewDummyHandler() (Handler, error) {
	return &DummyHandler{}, nil
}

type DummyHandler struct{}

func (h *DummyHandler) Start() error {
	return nil
}

func (h *DummyHandler) Stop() error {
	return nil
}

func (h *DummyHandler) GetName() string {
	return "Dummy Mode"
}

func (h *DummyHandler) GetPages() map[io.Page]string {
	return map[io.Page]string{io.PAGE_1: "Page 1"}
}

func (h *DummyHandler) SelectPage(page io.Page) (changes []io.ButtonChangeEvent) {
	return changes
}

func (h *DummyHandler) HandleInput(in io.InputEvent) (changes []io.ButtonChangeEvent) {
	return changes
}
