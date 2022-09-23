package live

import (
	"time"

	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/app/mode"
	"go.skymyer.dev/show-control/io"
	"go.skymyer.dev/show-control/show"
)

func init() {
	mode.Register("LIVE_MODE", NewHandler)
}

func NewHandler() (mode.Handler, error) {
	return &Handler{
		name: "Live mode",
		pages: map[io.Page]string{
			io.PAGE_1: "Page 1",
			io.PAGE_2: "Page 2", // temp show
		},
		pageButtons: map[io.Page][]io.GridButton{},
		pagesInit: map[io.Page][]io.GridButton{
			io.PAGE_1: {},
			io.PAGE_2: {},
			io.PAGE_3: {},
			io.PAGE_4: {},
		},
		cueButtonsList: []io.Button{ // TODO - make dynamic from io device grid
			io.BTN_1_1, io.BTN_1_2, io.BTN_1_3, io.BTN_1_4, io.BTN_1_5, io.BTN_1_6, io.BTN_1_7, io.BTN_1_8,
			io.BTN_2_1, io.BTN_2_2, io.BTN_2_3, io.BTN_2_4, io.BTN_2_5, io.BTN_2_6, io.BTN_2_7, io.BTN_2_8,
			io.BTN_3_1, io.BTN_3_2, io.BTN_3_3, io.BTN_3_4, io.BTN_3_5, io.BTN_3_6, io.BTN_3_7, io.BTN_3_8,
			io.BTN_4_1, io.BTN_4_2, io.BTN_4_3, io.BTN_4_4, io.BTN_4_5, io.BTN_4_6, io.BTN_4_7, io.BTN_4_8,
		},
		cueButtons: make(map[io.Button]string),
	}, nil
}

type selectableSequence struct {
	Name     string
	Selected bool
}

type Handler struct {
	name        string
	page        io.Page
	pages       map[io.Page]string
	pageButtons map[io.Page][]io.GridButton
	pagesInit   map[io.Page][]io.GridButton
	show        *show.Controller
	executor    *show.Executor

	cueButtonsList []io.Button
	cueButtons     map[io.Button]string
}

func (h *Handler) Start() (err error) {
	// Start show control
	h.show, err = show.NewFromConfig(
		"etc/setup.yaml",
		show.WithFixtureLibrary("etc/fixtures.yaml"),
		show.WithTickRate(25*time.Millisecond),
	)
	if err != nil {
		return err
	}

	// Setup executor
	h.executor, err = h.show.NewExecutor("default")
	if err != nil {
		return err
	}
	err = h.executor.Load(
		"etc/cues.yaml",
		"etc/effects.yaml",
		"etc/sequences.yaml",
	)
	if err != nil {
		return err
	}

	// Initialize dynamic buttons
	var cueCnt int
	h.cueButtons = map[io.Button]string{}
	for _, cue := range h.executor.GetSequenceNames() {
		button := h.cueButtonsList[cueCnt]
		logger.Default.Debug("adding sequence button", zap.String("name", cue), zap.Int("button", int(button)))
		h.cueButtons[button] = cue
		cueCnt++
	}

	return h.show.Run()
}

func (h *Handler) Stop() error {
	if h.show != nil {
		if err := h.show.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) GetName() string {
	return h.name
}

func (h *Handler) GetPages() map[io.Page]string {
	return h.pages
}

func (h *Handler) SelectPage(page io.Page) (changes []io.ButtonChangeEvent) {
	h.page = page

	// Static init buttons
	grid := h.pagesInit[page]

	// Dynamic button init
	switch page {
	case io.PAGE_1:
		grid = h.addButtons(grid, h.cueButtons)
	case io.PAGE_2:
	case io.PAGE_3:
	case io.PAGE_4:
	}

	return append(changes, io.ButtonChangeEvent{
		Partial: false,
		Grid:    grid,
	})
}

func (h *Handler) HandleInput(in io.InputEvent) (changes []io.ButtonChangeEvent) {
	switch h.page {
	case io.PAGE_1:
		switch {
		case in.IsButtonPress():
			if seq, ok := h.cueButtons[in.Button]; ok {
				logger.Default.Debug("sequence enable", zap.String("name", seq))
				h.executor.EnableSequences([]show.ExecutorSequence{{Name: seq}})
			}
		}
	case io.PAGE_2:
		switch {
		case in.IsButtonPress():
			logger.Default.Debug("page 2 start time coded show")
			tcs, err := show.NewTimeCodedShow("etc/show.yaml", 25*time.Millisecond)
			if err != nil {
				logger.Default.Fatal("tcs load", zap.Error(err))
			}
			tcs.Run(h.executor)

		}
	}
	return changes
}

func (h *Handler) addButtons(grid []io.GridButton, b ...map[io.Button]string) []io.GridButton {
	for _, list := range b {
		for btn := range list {
			grid = append(grid, io.GridButton{Button: btn, Kind: io.BTN_KIND_DEFAULT})
		}
	}
	return grid
}
