/*
 * Application Controller
 *
 */
package app

import (
	"fmt"
	"path/filepath"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/app/mode"
	"go.skymyer.dev/show-control/common"
	"go.skymyer.dev/show-control/db"
	"go.skymyer.dev/show-control/io"
	"go.skymyer.dev/show-control/io/novation"
	"go.uber.org/zap"

	_ "go.skymyer.dev/show-control/app/live"
	_ "go.skymyer.dev/show-control/app/program"
	//_ "go.skymyer.dev/show-control/app/calibrate"
	//_ "go.skymyer.dev/show-control/app/show"
)

var (
	Name     = "show-control"
	Database = "database.sqlite"
)

func New() (*Controller, error) {

	userDir, err := common.GetUserConfigDir(Name)
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(userDir, Database)
	repo, err := db.NewRepository(dbPath)
	if err != nil {
		return nil, err
	}

	midi, err := io.New(novation.DRIVER_NAME, novation.LAUNCHPAD_MINI_MK3)
	if err != nil {
		return nil, err
	}

	mode1Handler := mode.MustHandler("LIVE_MODE")
	mode2Handler := mode.MustHandler("DUMMY_MODE")
	mode3Handler := mode.MustHandler("DUMMY_MODE")
	mode4Handler := mode.MustHandler("PROGRAM_MODE")

	c := &Controller{
		repo:       repo,
		shutdownCh: make(chan bool, 1),
		io:         []io.Driver{midi},
		handlers: map[io.Mode]mode.Handler{
			io.MODE_1: mode1Handler,
			io.MODE_2: mode2Handler,
			io.MODE_3: mode3Handler,
			io.MODE_4: mode4Handler,
		},
		pages: map[io.Mode]map[io.Page]string{
			io.MODE_MAIN: {
				io.PAGE_1: mode1Handler.GetName(),
				io.PAGE_2: mode2Handler.GetName(),
				io.PAGE_3: mode3Handler.GetName(),
				io.PAGE_4: mode4Handler.GetName(),
			},
			io.MODE_1: mode1Handler.GetPages(),
			io.MODE_2: mode2Handler.GetPages(),
			io.MODE_3: mode3Handler.GetPages(),
			io.MODE_4: mode4Handler.GetPages(),
		},
	}
	return c, nil
}

type Controller struct {
	mode io.Mode
	page io.Page
	repo *db.Repository

	shutdownCh chan bool
	exitCount  int

	io       []io.Driver
	handlers map[io.Mode]mode.Handler
	pages    map[io.Mode]map[io.Page]string
}

func (c *Controller) Run() error {

	inputCh := make(chan io.InputEvent)
	defer close(inputCh)

	logger.Default.Debug("init io devices")
	for _, device := range c.io {
		if err := device.Open(inputCh); err != nil {
			return fmt.Errorf("open device: %v", err)
		}
	}

	logger.Default.Info("running")
	c.switchToMain()

	for {
		select {
		case <-c.shutdownCh:
			return c.shutdown()
		case in := <-inputCh:
			c.handleInput(in)
		}
	}
}

func (c *Controller) Stop() error {
	c.shutdownCh <- true
	return nil
}

func (c *Controller) shutdown() error {
	logger.Default.Debug("shutdown init")
	// Stop all handlers
	for _, h := range c.handlers {
		h.Stop()
	}

	// Clear all output
	c.sendEvents(io.ControlChangeEvent{}, io.ButtonChangeEvent{})

	// Close databse
	c.repo.Close()

	logger.Default.Debug("shutdown finished")
	return nil
}

func (c *Controller) handleInput(in io.InputEvent) {
	switch c.mode {
	case io.MODE_MAIN:
		switch {
		case in.IsControlPress():
			switch in.Control {
			case io.CTR_BACK:
				c.mainExit()
			case io.CTR_PAGE_1:
				c.switchMode(io.MODE_1)
			case io.CTR_PAGE_2:
				c.switchMode(io.MODE_2)
			case io.CTR_PAGE_3:
				c.switchMode(io.MODE_3)
			case io.CTR_PAGE_4:
				c.switchMode(io.MODE_4)
			}
		}
	default:
		switch {
		case in.IsControlPress():
			switch in.Control {
			case io.CTR_BACK:
				c.switchToMain()
			case io.CTR_PAGE_1:
				c.switchPage(io.PAGE_1)
			case io.CTR_PAGE_2:
				c.switchPage(io.PAGE_2)
			case io.CTR_PAGE_3:
				c.switchPage(io.PAGE_3)
			case io.CTR_PAGE_4:
				c.switchPage(io.PAGE_4)
			}
		case in.IsButton():
			for _, event := range c.handlers[c.mode].HandleInput(in) {
				c.sendEvents(event)
			}
		}
	}
}

func (c *Controller) mainExit() {
	c.exitCount++
	if c.exitCount >= 3 {
		logger.Default.Debug("main exit")
		c.Stop()
	}
}

func (c *Controller) switchToMain() {
	logger.Default.Debug("switch to main")
	c.exitCount = 0
	// Stop previous mode handler when going to main
	if h, ok := c.handlers[c.mode]; ok {
		h.Stop()
	}

	c.mode = io.MODE_MAIN
	c.page = io.PAGE_0

	c.sendEvents(
		c.newControlChange(c.mode, true),
		io.ButtonChangeEvent{},
	)
}

func (c *Controller) switchMode(mode io.Mode) {
	logger.Default.Debug("switch mode", zap.Int("mode", int(mode)))
	// Set new mode and default to page 1
	c.mode = mode
	c.page = io.PAGE_0

	// Start mode handler
	if err := c.handlers[c.mode].Start(); err != nil {
		logger.Default.Fatal("start handler", zap.Int("mode", int(mode)), zap.Error(err))
	}
	c.switchPage(io.PAGE_1)
}

func (c *Controller) switchPage(page io.Page) {
	// Skip if we are already on given page
	if page == c.page {
		return
	}

	// Skip if page does not exist
	if _, ok := c.pages[c.mode][page]; !ok {
		return
	}

	logger.Default.Debug("switch page", zap.Int("page", int(page)))

	c.page = page
	c.sendEvents(
		c.newControlChange(c.mode, true),
	)
	for _, event := range c.handlers[c.mode].SelectPage(c.page) {
		c.sendEvents(event)
	}
}

func (c *Controller) sendEvents(events ...interface{}) {
	for _, device := range c.io {
		device.Handle(events...)
	}
}

func (c *Controller) newControlChange(mode io.Mode, clear bool) io.ControlChangeEvent {
	var grid []io.GridControl

	// Page buttons
	for p, d := range c.pages[mode] {
		grid = append(grid, io.GridControl{
			Control:     io.PageToControl[p],
			Description: d,
		})
	}

	// Back button if not in main menu
	if mode != io.MODE_MAIN {
		grid = append(grid, io.GridControl{
			Control:     io.CTR_BACK,
			Description: "Back",
		})
	}

	return io.ControlChangeEvent{
		Mode:    c.mode,
		Page:    c.page,
		Partial: !clear,
		Grid:    grid,
	}
}
