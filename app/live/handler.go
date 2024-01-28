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

const (
	TYPE_SCENE     = "scene"
	TYPE_DIMMER    = "dim"
	TYPE_DIMMER_FX = "dimfx"
	TYPE_POSITION  = "pos"
)

func NewHandler() (mode.Handler, error) {
	logger.Default.Debug("create new live handler")
	return &Handler{
		name: "Live mode",
		pages: map[io.Page]string{
			io.PAGE_1: "Page 1",
			io.PAGE_2: "Page 2",
			io.PAGE_3: "Page 3",
			io.PAGE_4: "Page 4",
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
			io.BTN_5_1, io.BTN_5_2, io.BTN_5_3, io.BTN_5_4, io.BTN_5_5, io.BTN_5_6, io.BTN_5_7, io.BTN_5_8,
			io.BTN_6_1, io.BTN_6_2, io.BTN_6_3, io.BTN_6_4, io.BTN_6_5, io.BTN_6_6, io.BTN_6_7, io.BTN_6_8,
			io.BTN_7_1, io.BTN_7_2, io.BTN_7_3, io.BTN_7_4, io.BTN_7_5, io.BTN_7_6, io.BTN_7_7, io.BTN_7_8,
			io.BTN_8_1, io.BTN_8_2, io.BTN_8_3, io.BTN_8_4, io.BTN_8_5, io.BTN_8_6, io.BTN_8_7, io.BTN_8_8,
		},
		cueButtonsPage1: make(map[io.Button]*show.Sequence),
		cueButtonsPage2: make(map[io.Button]*show.Sequence),
		cueButtonsPage3: make(map[io.Button]*show.Sequence),
		cueButtonsPage4: make(map[io.Button]*show.Sequence),

		activeDimmerFx: make(map[string]bool),
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

	cueButtonsPage1 map[io.Button]*show.Sequence
	cueButtonsPage2 map[io.Button]*show.Sequence
	cueButtonsPage3 map[io.Button]*show.Sequence
	cueButtonsPage4 map[io.Button]*show.Sequence

	activeScene    string
	activeDimmer   string
	activePosition string

	activeDimmerFx map[string]bool
}

func (h *Handler) Start() (err error) {
	logger.Default.Debug("load live handler show")
	// Start show control
	h.show, err = show.NewFromConfig(
		"etc/setup.yaml",
		show.WithFixtureLibrary("etc/fixtures.yaml"),
		show.WithTickRate(22*time.Millisecond),
	)
	logger.Default.Debug("live handler show loaded")

	if err != nil {
		logger.Default.Warn("start live handler", zap.Error(err))
		return err
	}

	// Setup executor
	logger.Default.Debug("setup executor")
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

	logger.Default.Debug("init live dynamic buttons")

	for _, seq := range h.executor.GetLiveSequences(io.PAGE_1) {
		seq.Loop = true
		buttonNum := seq.Button - 1
		button := h.cueButtonsList[buttonNum]
		logger.Default.Debug("adding sequence button", zap.String("name", seq.Name), zap.Int8("page", seq.Page), zap.Int("button", int(button)))
		h.cueButtonsPage1[button] = seq
	}

	for _, seq := range h.executor.GetLiveSequences(io.PAGE_2) {
		seq.Loop = true
		buttonNum := seq.Button - 1
		button := h.cueButtonsList[buttonNum]
		logger.Default.Debug("adding sequence button", zap.String("name", seq.Name), zap.Int8("page", seq.Page), zap.Int("button", int(button)))
		h.cueButtonsPage2[button] = seq
	}

	for _, seq := range h.executor.GetLiveSequences(io.PAGE_3) {
		seq.Loop = true
		buttonNum := seq.Button - 1
		button := h.cueButtonsList[buttonNum]
		logger.Default.Debug("adding sequence button", zap.String("name", seq.Name), zap.Int8("page", seq.Page), zap.Int("button", int(button)))
		h.cueButtonsPage3[button] = seq
	}

	for _, seq := range h.executor.GetLiveSequences(io.PAGE_4) {
		seq.Loop = true
		buttonNum := seq.Button - 1
		button := h.cueButtonsList[buttonNum]
		logger.Default.Debug("adding sequence button", zap.String("name", seq.Name), zap.Int8("page", seq.Page), zap.Int("button", int(button)))
		h.cueButtonsPage4[button] = seq
	}

	logger.Default.Debug("starting show")
	return h.show.Run()
}

func (h *Handler) Stop() error {
	h.activeScene = ""
	h.activeDimmer = ""
	h.activePosition = ""

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
		grid = h.addSequenceButtons(grid, h.cueButtonsPage1)
	case io.PAGE_2:
		grid = h.addSequenceButtons(grid, h.cueButtonsPage2)
	case io.PAGE_3:
		grid = h.addSequenceButtons(grid, h.cueButtonsPage3)
	case io.PAGE_4:
		grid = h.addSequenceButtons(grid, h.cueButtonsPage4)
	}

	return append(changes, io.ButtonChangeEvent{
		Partial: false,
		Grid:    grid,
	})
}

func (h *Handler) handleSelect(seq *show.Sequence) {
	switch seq.Type {
	case TYPE_SCENE:
		if seq.Name != h.activeScene {
			logger.Default.Debug("set scene", zap.String("new", seq.Name), zap.String("old", h.activeScene))
			if h.activeScene != "" {
				h.executor.DisableSequences([]string{h.activeScene})
			}
			h.executor.EnableSequences([]show.ExecutorSequence{{Name: seq.Name}})
			h.activeScene = seq.Name
		}
	case TYPE_DIMMER:
		if seq.Name != h.activeDimmer {
			logger.Default.Debug("set dimmer", zap.String("new", seq.Name), zap.String("old", h.activeDimmer))
			if h.activeDimmer != "" {
				h.executor.DisableSequences([]string{h.activeDimmer})
			}
			h.executor.EnableSequences([]show.ExecutorSequence{{Name: seq.Name}})
			h.activeDimmer = seq.Name
		}

	case TYPE_POSITION:
		if seq.Name != h.activePosition {
			logger.Default.Debug("set position", zap.String("new", seq.Name), zap.String("old", h.activePosition))
			if h.activePosition != "" {
				h.executor.DisableSequences([]string{h.activePosition})
			}
			h.executor.EnableSequences([]show.ExecutorSequence{{Name: seq.Name}})
			h.activePosition = seq.Name
		}

	case TYPE_DIMMER_FX:
		if _, ok := h.activeDimmerFx[seq.Name]; ok {
			delete(h.activeDimmerFx, seq.Name)
			logger.Default.Debug("del dimmer fx", zap.String("name", seq.Name))
			h.executor.DisableSequences([]string{seq.Name})
		} else {
			h.activeDimmerFx[seq.Name] = true
			logger.Default.Debug("add dimmer fx", zap.String("name", seq.Name))
			h.executor.EnableSequences([]show.ExecutorSequence{{Name: seq.Name}})
		}

	default:
		logger.Default.Warn("unknown sequence type", zap.String("type", seq.Type))
	}
}

func (h *Handler) HandleInput(in io.InputEvent) (changes []io.ButtonChangeEvent) {
	switch h.page {
	case io.PAGE_1:
		switch {
		case in.IsButtonPress():
			if seq, ok := h.cueButtonsPage1[in.Button]; ok {
				h.handleSelect(seq)
			}
		}
	case io.PAGE_2:
		switch {
		case in.IsButtonPress():
			//logger.Default.Debug("page 2 start time coded show")
			//tcs, err := show.NewTimeCodedShow("etc/show.yaml", 25*time.Millisecond)
			//if err != nil {
			//	logger.Default.Fatal("tcs load", zap.Error(err))
			//}
			//tcs.Run(h.executor)
			if seq, ok := h.cueButtonsPage2[in.Button]; ok {
				h.handleSelect(seq)
			}
		}
	case io.PAGE_3:
		switch {
		case in.IsButtonPress():
			if seq, ok := h.cueButtonsPage3[in.Button]; ok {
				h.handleSelect(seq)
			}
		}
	case io.PAGE_4:
		switch {
		case in.IsButtonPress():
			if seq, ok := h.cueButtonsPage4[in.Button]; ok {
				h.handleSelect(seq)
			}
		}
	}
	return changes
}

func (h *Handler) addSequenceButtons(grid []io.GridButton, b ...map[io.Button]*show.Sequence) []io.GridButton {
	for _, list := range b {
		for btn := range list {
			grid = append(grid, io.GridButton{Button: btn, Kind: io.BTN_KIND_DEFAULT})
		}
	}
	return grid
}

func (h *Handler) addButtons(grid []io.GridButton, b ...map[io.Button]string) []io.GridButton {
	for _, list := range b {
		for btn := range list {
			grid = append(grid, io.GridButton{Button: btn, Kind: io.BTN_KIND_DEFAULT})
		}
	}
	return grid
}
