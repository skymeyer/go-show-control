package mode

import (
	"log"

	"go.skymyer.dev/show-control/io"
)

type Handler interface {
	Start() error
	Stop() error

	GetName() string
	GetPages() map[io.Page]string

	SelectPage(io.Page) []io.ButtonChangeEvent
	HandleInput(in io.InputEvent) []io.ButtonChangeEvent
}

type modeFactory func() (Handler, error)

var handlers = make(map[string]modeFactory)

func Register(name string, f modeFactory) {
	if _, exists := handlers[name]; exists {
		log.Fatalf("mode handler %q already exists", name)
	}
	handlers[name] = f
}

func MustHandler(name string) Handler {
	if _, ok := handlers[name]; !ok {
		log.Fatalf("invalid mode handler %q", name)
	}
	h, err := handlers[name]()
	if err != nil {
		log.Fatalf("load mode %q: %v", name, err)
	}
	return h

}
