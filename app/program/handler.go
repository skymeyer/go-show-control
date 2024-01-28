package program

import (
	"fmt"
	"time"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/app/mode"
	"go.skymyer.dev/show-control/io"
	"go.skymyer.dev/show-control/library/feature/dimmer"
	"go.skymyer.dev/show-control/library/feature/gobo"
	"go.skymyer.dev/show-control/library/feature/laser"
	"go.skymyer.dev/show-control/library/feature/pantilt"
	"go.skymyer.dev/show-control/library/feature/rgbw"
	"go.skymyer.dev/show-control/library/feature/strobe"
	"go.skymyer.dev/show-control/show"
	"go.uber.org/zap"
)

func init() {
	mode.Register("PROGRAM_MODE", NewHandler)
}

func NewHandler() (mode.Handler, error) {
	return &Handler{
		name: "Programming mode",
		pages: map[io.Page]string{
			io.PAGE_1: "RGBW",
			io.PAGE_2: "Colors & Gobos",
			io.PAGE_3: "Dimmers & Positions",
			io.PAGE_4: "Effects",
		},
		pagesInit: map[io.Page][]io.GridButton{
			io.PAGE_1: {
				// RGBW widgets
				{Button: io.BTN_1_1, Kind: io.BTN_KIND_WIDGET_SLIDER_RED},
				{Button: io.BTN_2_1, Kind: io.BTN_KIND_WIDGET_SLIDER_GREEN},
				{Button: io.BTN_3_1, Kind: io.BTN_KIND_WIDGET_SLIDER_BLUE},
				{Button: io.BTN_4_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE},
			},
			io.PAGE_2: {},
			io.PAGE_3: {
				{Button: io.BTN_1_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE}, // dimmer
				{Button: io.BTN_2_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE}, // strobe
				{Button: io.BTN_5_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE}, // pan
				{Button: io.BTN_6_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE}, // tilt
			},
			io.PAGE_4: {
				{Button: io.BTN_1_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE}, // pattern
				{Button: io.BTN_2_1, Kind: io.BTN_KIND_WIDGET_SLIDER_WHITE}, // size
				{Button: io.BTN_3_1, Kind: io.BTN_KIND_WIDGET_SLIDER_GREEN}, // angle
				{Button: io.BTN_4_1, Kind: io.BTN_KIND_WIDGET_SLIDER_GREEN}, // hangle
				{Button: io.BTN_5_1, Kind: io.BTN_KIND_WIDGET_SLIDER_GREEN}, // vangle
				{Button: io.BTN_6_1, Kind: io.BTN_KIND_WIDGET_SLIDER_BLUE},  // hpos
				{Button: io.BTN_7_1, Kind: io.BTN_KIND_WIDGET_SLIDER_BLUE},  // vpos
				{Button: io.BTN_8_1, Kind: io.BTN_KIND_WIDGET_SLIDER_RED},   // color
			},
		},
	}, nil
}

type selectableFixture struct {
	Name     string
	Feature  string
	Selected bool
}

type featureSelect struct {
	Name    string
	Feature string
	Kind    io.ButtonKind
}

type Handler struct {
	name      string
	page      io.Page
	pages     map[io.Page]string
	pagesInit map[io.Page][]io.GridButton
	show      *show.Controller

	// Page 1 - RGBW controls
	rgbFixtures map[io.Button]*selectableFixture
	rgbWidget   *RGBW

	// Page 2 - Colors & Gobos
	arrowsPage3 *ArrowsCallControl

	scanFixtures map[io.Button]*selectableFixture
	scanColors   map[io.Button]*featureSelect
	mhFixtures   map[io.Button]*selectableFixture
	mhColors     map[io.Button]*featureSelect
	mhGobos      map[io.Button]*featureSelect

	// Page 3 - Dimmers & Position
	sliderFixtures map[io.Button]*selectableFixture
	dimmer         *Slider
	strobe         *Slider

	moveFixtures map[io.Button]*selectableFixture
	panTilt      *PanTilt

	// Page 4 - Laser control
	arrowsPage4  *ArrowsCallControl
	laserPattern *Slider
	laserSize    *Slider
	laserAngle   *Slider
	laserHAngle  *Slider
	laserVAngle  *Slider
	laserHPos    *Slider
	laserVPos    *Slider
	laserColor   *Slider
}

func (h *Handler) Start() (err error) {
	// FIXME: dynamically load from config
	h.rgbFixtures = map[io.Button]*selectableFixture{
		io.BTN_8_1: {Name: "strip", Feature: "rgb8"},
		io.BTN_8_2: {Name: "strip", Feature: "rgb7"},
		io.BTN_8_3: {Name: "strip", Feature: "rgb6"},
		io.BTN_8_4: {Name: "strip", Feature: "rgb5"},
		io.BTN_8_5: {Name: "strip", Feature: "rgb4"},
		io.BTN_8_6: {Name: "strip", Feature: "rgb3"},
		io.BTN_8_7: {Name: "strip", Feature: "rgb2"},
		io.BTN_8_8: {Name: "strip", Feature: "rgb1"},
		io.BTN_7_1: {Name: "wl", Feature: "rgb"},
		io.BTN_7_3: {Name: "w1", Feature: "rgb"},
		io.BTN_7_4: {Name: "w2", Feature: "rgb"},
		io.BTN_7_5: {Name: "w3", Feature: "rgb"},
		io.BTN_7_6: {Name: "w4", Feature: "rgb"},
		io.BTN_7_8: {Name: "wr", Feature: "rgb"},
	}

	h.rgbWidget = &RGBW{
		Red:   []io.Button{io.BTN_1_1, io.BTN_1_2, io.BTN_1_3, io.BTN_1_4, io.BTN_1_5, io.BTN_1_6, io.BTN_1_7, io.BTN_1_8},
		Green: []io.Button{io.BTN_2_1, io.BTN_2_2, io.BTN_2_3, io.BTN_2_4, io.BTN_2_5, io.BTN_2_6, io.BTN_2_7, io.BTN_2_8},
		Blue:  []io.Button{io.BTN_3_1, io.BTN_3_2, io.BTN_3_3, io.BTN_3_4, io.BTN_3_5, io.BTN_3_6, io.BTN_3_7, io.BTN_3_8},
		White: []io.Button{io.BTN_4_1, io.BTN_4_2, io.BTN_4_3, io.BTN_4_4, io.BTN_4_5, io.BTN_4_6, io.BTN_4_7, io.BTN_4_8},
	}

	h.scanFixtures = map[io.Button]*selectableFixture{
		io.BTN_1_1: {Name: "s1"},
		io.BTN_2_1: {Name: "s2"},
		io.BTN_2_2: {Name: "s3"},
		io.BTN_1_2: {Name: "s4"},
	}

	h.scanColors = map[io.Button]*featureSelect{
		io.BTN_1_3: {Name: "white", Feature: "color", Kind: io.BTN_KIND_COLOR_WHITE},
		io.BTN_1_4: {Name: "yellow", Feature: "color", Kind: io.BTN_KIND_COLOR_YELLOW},
		io.BTN_1_5: {Name: "red", Feature: "color", Kind: io.BTN_KIND_COLOR_RED},
		io.BTN_1_6: {Name: "green", Feature: "color", Kind: io.BTN_KIND_COLOR_GREEN},
		io.BTN_1_7: {Name: "blue", Feature: "color", Kind: io.BTN_KIND_COLOR_BLUE},
		io.BTN_2_3: {Name: "tri", Feature: "color", Kind: io.BTN_KIND_DUAL_GREEN_BLUE},
		io.BTN_2_4: {Name: "purple", Feature: "color", Kind: io.BTN_KIND_COLOR_MAGENTA},
		io.BTN_2_5: {Name: "orange", Feature: "color", Kind: io.BTN_KIND_COLOR_ORANGE},
		io.BTN_2_6: {Name: "cyan", Feature: "color", Kind: io.BTN_KIND_DUAL_MAGENTA_CYAN},
	}

	h.mhFixtures = map[io.Button]*selectableFixture{
		io.BTN_8_3: {Name: "mh1"},
		io.BTN_8_4: {Name: "mh2"},
		io.BTN_8_5: {Name: "mh3"},
		io.BTN_8_6: {Name: "mh4"},
	}

	h.mhGobos = map[io.Button]*featureSelect{
		io.BTN_4_1: {Name: "none", Feature: "figure", Kind: io.BTN_KIND_COLOR_GRAY},
		io.BTN_4_2: {Name: "flower", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_4_3: {Name: "circles", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_4_4: {Name: "triangle", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_4_5: {Name: "star", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_4_6: {Name: "sun", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_4_7: {Name: "whirlpool", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_4_8: {Name: "glass", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_2: {Name: "flower-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_3: {Name: "circles-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_4: {Name: "triangle-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_5: {Name: "star-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_6: {Name: "sun-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_7: {Name: "whirlpool-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
		io.BTN_5_8: {Name: "glass-dither", Feature: "figure", Kind: io.BTN_KIND_DEFAULT},
	}

	h.mhColors = map[io.Button]*featureSelect{
		io.BTN_6_1: {Name: "white", Feature: "color", Kind: io.BTN_KIND_COLOR_WHITE},
		io.BTN_6_2: {Name: "red", Feature: "color", Kind: io.BTN_KIND_COLOR_RED},
		io.BTN_6_3: {Name: "orange", Feature: "color", Kind: io.BTN_KIND_COLOR_ORANGE},
		io.BTN_6_4: {Name: "yellow", Feature: "color", Kind: io.BTN_KIND_COLOR_YELLOW},
		io.BTN_6_5: {Name: "green", Feature: "color", Kind: io.BTN_KIND_COLOR_GREEN},
		io.BTN_6_6: {Name: "blue", Feature: "color", Kind: io.BTN_KIND_COLOR_BLUE},
		io.BTN_6_7: {Name: "magenta", Feature: "color", Kind: io.BTN_KIND_COLOR_MAGENTA},
		io.BTN_6_8: {Name: "cyan", Feature: "color", Kind: io.BTN_KIND_COLOR_CYAN},
		io.BTN_7_2: {Name: "magenta-cyan", Feature: "color", Kind: io.BTN_KIND_DUAL_MAGENTA_CYAN},
		io.BTN_7_3: {Name: "magenta-blue", Feature: "color", Kind: io.BTN_KIND_DUAL_MAGENTA_BLUE},
		io.BTN_7_4: {Name: "green-blue", Feature: "color", Kind: io.BTN_KIND_DUAL_GREEN_BLUE},
		io.BTN_7_5: {Name: "green-yellow", Feature: "color", Kind: io.BTN_KIND_DUAL_GREEN_YELLOW},
		io.BTN_7_6: {Name: "orange-yellow", Feature: "color", Kind: io.BTN_KIND_DUAL_ORANGE_YELLOW},
		io.BTN_7_7: {Name: "orange-red", Feature: "color", Kind: io.BTN_KIND_DUAL_ORANGE_RED},
	}

	h.dimmer = &Slider{
		Buttons: []io.Button{io.BTN_1_1, io.BTN_1_2, io.BTN_1_3, io.BTN_1_4, io.BTN_1_5, io.BTN_1_6, io.BTN_1_7, io.BTN_1_8},
	}

	h.strobe = &Slider{
		Buttons: []io.Button{io.BTN_2_1, io.BTN_2_2, io.BTN_2_3, io.BTN_2_4, io.BTN_2_5, io.BTN_2_6, io.BTN_2_7, io.BTN_2_8},
	}

	h.sliderFixtures = map[io.Button]*selectableFixture{
		io.BTN_3_1: {Name: "wl"},
		io.BTN_3_2: {Name: "s1"},
		io.BTN_4_2: {Name: "s2"},
		io.BTN_3_3: {Name: "w1"},
		io.BTN_3_4: {Name: "w2"},
		io.BTN_3_5: {Name: "w3"},
		io.BTN_3_6: {Name: "w4"},
		io.BTN_4_7: {Name: "s3"},
		io.BTN_3_7: {Name: "s4"},
		io.BTN_3_8: {Name: "wr"},
		io.BTN_4_3: {Name: "mh1"},
		io.BTN_4_4: {Name: "mh2"},
		io.BTN_4_5: {Name: "mh3"},
		io.BTN_4_6: {Name: "mh4"},
		io.BTN_4_8: {Name: "strip"},
	}

	h.moveFixtures = map[io.Button]*selectableFixture{
		io.BTN_7_2: {Name: "s1"},
		io.BTN_8_2: {Name: "s2"},
		io.BTN_8_3: {Name: "mh1"},
		io.BTN_8_4: {Name: "mh2"},
		io.BTN_8_5: {Name: "mh3"},
		io.BTN_8_6: {Name: "mh4"},
		io.BTN_8_7: {Name: "s3"},
		io.BTN_7_7: {Name: "s4"},
		io.BTN_8_8: {Name: "strip"},
	}

	h.panTilt = &PanTilt{
		ButtonsPan:  []io.Button{io.BTN_5_1, io.BTN_5_2, io.BTN_5_3, io.BTN_5_4, io.BTN_5_5, io.BTN_5_6, io.BTN_5_7, io.BTN_5_8},
		ButtonsTilt: []io.Button{io.BTN_6_1, io.BTN_6_2, io.BTN_6_3, io.BTN_6_4, io.BTN_6_5, io.BTN_6_6, io.BTN_6_7, io.BTN_6_8},
	}

	h.arrowsPage3 = &ArrowsCallControl{
		Selectors: []*ArrowCall{
			{Button: Call1, Handle: h.dimmer.HandleArrow, Step: 1},
			{Button: Call2, Handle: h.strobe.HandleArrow, Step: 1},
			{Button: Call5, Handle: h.panTilt.HandleArrowPan, Step: 100},
			{Button: Call6, Handle: h.panTilt.HandleArrowTilt, Step: 100},
		},
	}

	// PAGE 4

	h.laserPattern = &Slider{
		Buttons: []io.Button{io.BTN_1_1, io.BTN_1_2, io.BTN_1_3, io.BTN_1_4, io.BTN_1_5, io.BTN_1_6, io.BTN_1_7, io.BTN_1_8},
	}
	h.laserSize = &Slider{
		Buttons: []io.Button{io.BTN_2_1, io.BTN_2_2, io.BTN_2_3, io.BTN_2_4, io.BTN_2_5, io.BTN_2_6, io.BTN_2_7, io.BTN_2_8},
	}
	h.laserAngle = &Slider{
		Buttons: []io.Button{io.BTN_3_1, io.BTN_3_2, io.BTN_3_3, io.BTN_3_4, io.BTN_3_5, io.BTN_3_6, io.BTN_3_7, io.BTN_3_8},
	}
	h.laserHAngle = &Slider{
		Buttons: []io.Button{io.BTN_4_1, io.BTN_4_2, io.BTN_4_3, io.BTN_4_4, io.BTN_4_5, io.BTN_4_6, io.BTN_4_7, io.BTN_4_8},
	}
	h.laserVAngle = &Slider{
		Buttons: []io.Button{io.BTN_5_1, io.BTN_5_2, io.BTN_5_3, io.BTN_5_4, io.BTN_5_5, io.BTN_5_6, io.BTN_5_7, io.BTN_5_8},
	}
	h.laserHPos = &Slider{
		Buttons: []io.Button{io.BTN_6_1, io.BTN_6_2, io.BTN_6_3, io.BTN_6_4, io.BTN_6_5, io.BTN_6_6, io.BTN_6_7, io.BTN_6_8},
	}
	h.laserVPos = &Slider{
		Buttons: []io.Button{io.BTN_7_1, io.BTN_7_2, io.BTN_7_3, io.BTN_7_4, io.BTN_7_5, io.BTN_7_6, io.BTN_7_7, io.BTN_7_8},
	}
	h.laserColor = &Slider{
		Buttons: []io.Button{io.BTN_8_1, io.BTN_8_2, io.BTN_8_3, io.BTN_8_4, io.BTN_8_5, io.BTN_8_6, io.BTN_8_7, io.BTN_8_8},
	}

	h.arrowsPage4 = &ArrowsCallControl{
		Selectors: []*ArrowCall{
			{Button: Call1, Handle: h.laserPattern.HandleArrow, Step: 1},
			{Button: Call2, Handle: h.laserSize.HandleArrow, Step: 1},
			{Button: Call3, Handle: h.laserAngle.HandleArrow, Step: 1},
			{Button: Call4, Handle: h.laserHAngle.HandleArrow, Step: 1},
			{Button: Call5, Handle: h.laserVAngle.HandleArrow, Step: 1},
			{Button: Call6, Handle: h.laserHPos.HandleArrow, Step: 1},
			{Button: Call7, Handle: h.laserVPos.HandleArrow, Step: 1},
			// No Call8 for bottom slider color
		},
	}

	// Start show control
	h.show, err = show.NewFromConfig(
		"etc/setup.yaml",
		show.WithFixtureLibrary("etc/fixtures.yaml"),
		show.WithTickRate(22*time.Millisecond),
	)
	if err != nil {
		return fmt.Errorf("show: %v", err)
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
		grid = h.addSelectableFixtures(grid, h.rgbFixtures)
	case io.PAGE_2:
		grid = h.addSelectableFixtures(grid, h.scanFixtures, h.mhFixtures)
		grid = h.addfeatureSelectors(grid, h.scanColors, h.mhGobos, h.mhColors)
	case io.PAGE_3:
		grid = h.addSelectableFixtures(grid, h.sliderFixtures)
		grid = h.addSelectableFixtures(grid, h.moveFixtures)
		grid = h.arrowsPage3.addArrowsSelectors(grid)
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
			// Select fixture button handling
			changes = h.handleFixtureSelect(in.Button, changes, h.rgbFixtures)

			// Handle color selection
			if h.rgbWidget.Handle(in.Button) {
				// Grab new values
				feature := rgbw.Feature{
					Red:   &h.rgbWidget.RedValue,
					Green: &h.rgbWidget.GreenValue,
					Blue:  &h.rgbWidget.BlueValue,
					White: &h.rgbWidget.WhiteValue,
				}
				// Apply new values to selected fixtures
				for _, fixture := range h.rgbFixtures {
					if len(fixture.Name) > 0 && len(fixture.Feature) > 0 && fixture.Selected {
						h.show.SetFixtureFeature(show.Feature{
							Name:    fixture.Feature,
							Fixture: fixture.Name,
							Config:  feature,
						})
					}
				}
			}
		}
	case io.PAGE_2:
		switch {
		case in.IsButtonPress():
			// Select fixture button handling
			changes = h.handleFixtureSelect(in.Button, changes, h.scanFixtures, h.mhFixtures)

			// Handle selecting gobo
			h.handlefeatureSelect(in.Button, h.scanFixtures, h.scanColors)
			h.handlefeatureSelect(in.Button, h.mhFixtures, h.mhGobos, h.mhColors)
		}
	case io.PAGE_3:
		switch {
		case in.IsButtonPress():
			// Select fixture button handling
			changes = h.handleFixtureSelect(in.Button, changes, h.sliderFixtures)
			changes = h.handleFixtureSelect(in.Button, changes, h.moveFixtures)

			// Handle arrow call control
			changes = h.arrowsPage3.handleArrowSelect(in.Button, changes)

			// Handle arrow usage (FIXME, combine handle/action into slider,rgbw objects ...)
			arrows := h.arrowsPage3.Handle(in.Button)

			// Handle dimmer selection
			if h.dimmer.Handle(in.Button) || arrows {
				// Apply new values to selected fixtures
				for _, fixture := range h.sliderFixtures {
					if len(fixture.Name) > 0 && fixture.Selected {
						value := uint16(h.dimmer.Value) * 256
						h.show.SetFixtureFeature(show.Feature{
							Name:    "dimmer",
							Fixture: fixture.Name,
							Config: dimmer.Feature{
								Value: &value,
							},
						})
					}
				}
			}

			// Handle strobe selection
			if h.strobe.Handle(in.Button) || arrows {
				// Apply new values to selected fixtures
				for _, fixture := range h.sliderFixtures {
					if len(fixture.Name) > 0 && fixture.Selected {
						h.show.SetFixtureFeature(show.Feature{
							Name:    "strobe",
							Fixture: fixture.Name,
							Config: strobe.Feature{
								Value: h.strobe.Value,
							},
						})
					}
				}
			}

			// Handle tilt selection
			if h.panTilt.Handle(in.Button) || arrows {
				speed := uint16(255)
				// Apply new values to selected fixtures
				for _, fixture := range h.moveFixtures {
					if len(fixture.Name) > 0 && fixture.Selected {
						h.show.SetFixtureFeature(show.Feature{
							Fixture: fixture.Name,
							Name:    "position",
							Config: pantilt.Feature{
								Pan:   &h.panTilt.PanValue,
								Tilt:  &h.panTilt.TiltValue,
								Speed: &speed,
							},
						})
					}
				}
			}
		}
	case io.PAGE_4:
		switch {
		case in.IsButtonPress():
			// Handle arrow call control
			changes = h.arrowsPage4.handleArrowSelect(in.Button, changes)

			// Handle arrow usage (FIXME, combine handle/action into slider,rgbw objects ...)
			h.arrowsPage4.Handle(in.Button)

			h.laserPattern.Handle(in.Button)
			h.laserSize.Handle(in.Button)
			h.laserAngle.Handle(in.Button)
			h.laserHAngle.Handle(in.Button)
			h.laserVAngle.Handle(in.Button)
			h.laserHPos.Handle(in.Button)
			h.laserVPos.Handle(in.Button)
			h.laserColor.Handle(in.Button)

			laserEnabled := true
			if h.laserPattern.Value == 0 && h.laserColor.Value == 0 {
				laserEnabled = false
			}

			logger.Default.Debug("laser",
				zap.Uint8("pattern", h.laserPattern.Value),
				zap.Uint8("size", h.laserSize.Value),
				zap.Uint8("angle", h.laserAngle.Value),
				zap.Uint8("hangle", h.laserHAngle.Value),
				zap.Uint8("vangle", h.laserVAngle.Value),
				zap.Uint8("hpos", h.laserHPos.Value),
				zap.Uint8("vpos", h.laserVPos.Value),
				zap.Uint8("color", h.laserColor.Value),
			)

			h.show.SetFixtureFeature(show.Feature{
				Name:    "manual",
				Fixture: "laser-chauvet", // laser-rgb-1
				Config: laser.Feature{
					//Dynamic: laserEnabled,
					Static:  laserEnabled,
					Color:   h.laserColor.Value,
					Pattern: h.laserPattern.Value,
					Size:    h.laserSize.Value,
					Angle:   h.laserAngle.Value,
					HAngle:  h.laserHAngle.Value,
					VAngle:  h.laserVAngle.Value,
					HPos:    h.laserHPos.Value,
					VPos:    h.laserVPos.Value,
				},
			})
		}
	}
	return changes
}

func (h *Handler) addSelectableFixtures(grid []io.GridButton, s ...map[io.Button]*selectableFixture) []io.GridButton {
	for _, list := range s {
		for btn, fixture := range list {
			kind := io.BTN_KIND_DEFAULT
			if fixture.Selected {
				kind = io.BTN_KIND_DEFAULT_SELECT
			}
			grid = append(grid, io.GridButton{Button: btn, Kind: kind})
		}
	}
	return grid
}

func (h *Handler) handlefeatureSelect(button io.Button, f map[io.Button]*selectableFixture, g ...map[io.Button]*featureSelect) {
	for _, list := range g {
		if g, ok := list[button]; ok {
			for _, fixture := range f {
				if len(fixture.Name) > 0 && len(g.Feature) > 0 && len(g.Name) > 0 && fixture.Selected {
					h.show.SetFixtureFeature(show.Feature{
						Name:    g.Feature,
						Fixture: fixture.Name,
						Config: gobo.Feature{
							Gobo: g.Name,
						},
					})
				}
			}
		}
	}
}

func (h *Handler) addfeatureSelectors(grid []io.GridButton, s ...map[io.Button]*featureSelect) []io.GridButton {
	for _, list := range s {
		for btn, g := range list {
			grid = append(grid, io.GridButton{Button: btn, Kind: g.Kind})
		}
	}
	return grid
}

func (h *Handler) handleFixtureSelect(btn io.Button, c []io.ButtonChangeEvent, s ...map[io.Button]*selectableFixture) []io.ButtonChangeEvent {
	for _, list := range s {
		if fixture, ok := list[btn]; ok {
			var kind io.ButtonKind
			if fixture.Selected {
				fixture.Selected = false
				kind = io.BTN_KIND_DEFAULT
			} else {
				fixture.Selected = true
				kind = io.BTN_KIND_DEFAULT_SELECT
			}
			c = append(c, io.ButtonChangeEvent{
				Partial: true,
				Grid: []io.GridButton{{
					Button: btn,
					Kind:   kind,
				},
				},
			})
		}
	}
	return c
}
