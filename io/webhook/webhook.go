package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/io"
)

func init() {
	io.Register(DRIVER_NAME, New)
}

const (
	DRIVER_NAME = "webhook"
)

func New(device string) (io.Driver, error) {
	return &Webhook{}, nil
}

type Webhook struct {
	input     chan<- io.InputEvent
	eventLock sync.Mutex
	server    *http.Server
}

func (w *Webhook) Open(input chan<- io.InputEvent) error {

	w.input = input

	router := mux.NewRouter()
	io := router.PathPrefix("/io").Subrouter()
	io.HandleFunc("/control/{key}", w.controlHandler).Methods("POST")
	io.HandleFunc("/button/{key}", w.buttonHandler).Methods("POST")

	addr := ":8765" // FIXME use config
	w.server = &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Default.Info("webhook server startup", zap.String("addr", addr))
		err := w.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Default.Info("webhook server closed")
		} else if err != nil {
			logger.Default.Error("webhook server", zap.Error(err))
		}
	}()

	return nil
}

func (w *Webhook) controlHandler(wr http.ResponseWriter, req *http.Request) {
	var res bool
	vars := mux.Vars(req)
	if key, ok := controlMap[vars["key"]]; ok {
		logger.Default.Info("webhook control event", zap.String("key", vars["key"]))
		w.input <- io.InputEvent{
			Action:  io.ACTION_PRESS,
			Control: key,
		}
		res = true
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]bool{"ok": res})
}

func (w *Webhook) buttonHandler(wr http.ResponseWriter, req *http.Request) {
	var res bool
	vars := mux.Vars(req)
	if key, ok := buttonMap[vars["key"]]; ok {
		logger.Default.Info("webhook button event", zap.String("key", vars["key"]))
		w.input <- io.InputEvent{
			Action: io.ACTION_PRESS,
			Button: key,
		}
		res = true
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]bool{"ok": res})
}

func (w *Webhook) Handle(event ...interface{}) {
	w.eventLock.Lock()
	defer w.eventLock.Unlock()
}

func (w *Webhook) Close() error {
	if w.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return w.server.Shutdown(ctx)
	}
	return nil
}

func (w *Webhook) GetDevices() []string {
	return []string{}
}

func (w *Webhook) GetVersion() string {
	return "webhook"
}

func (w *Webhook) Sleep() error {
	return nil
}
func (w *Webhook) Wakeup() error {
	return nil
}
