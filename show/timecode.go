package show

import (
	"fmt"
	"time"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
	"go.uber.org/zap"
)

type TimeCodedShow struct {
	Name     string               `yaml:"name,omitempty"`
	Offset   time.Duration        `yaml:"offset,omitempty"`
	Schedule []*TimeCodedSequence `yaml:"schedule,omitempty"`
}

type TimeCodedSequence struct {
	// Required properties
	Sequences []string      `yaml:"seq,omitempty"`
	Start     time.Duration `yaml:"start,omitempty"`

	// Optional properties passed to the executor for the sequence list
	Stop           *time.Duration `yaml:"stop,omitempty"`
	BeatsPerMinute *uint8         `yaml:"bpm,omitempty"`
	Loop           *bool          `yaml:"loop,omitempty"`
}

func NewTimeCodedShow(file string, tickRate time.Duration) (*TimecodeShow, error) {

	var show = &TimeCodedShow{}
	if err := common.LoadFromFile(file, show); err != nil {
		return nil, err
	}

	var (
		queue         = NewQueue()
		previousStart time.Duration
	)
	for k, tcs := range show.Schedule {
		if tcs.Start < previousStart {
			return nil, fmt.Errorf("[%s] non sequential entry [%d] %v", show.Name, k, tcs)
		}
		if tcs.Stop != nil && *tcs.Stop < tcs.Start {
			return nil, fmt.Errorf("[%s] stop lower than start [%d] %v", show.Name, k, tcs)
		}

		// Only queue sequences after offset
		if tcs.Start-show.Offset < 0 {
			continue
		}

		// Correct start/stop based on offset
		tcs.Start = tcs.Start - show.Offset
		if tcs.Stop != nil {
			stop := *tcs.Stop - show.Offset
			tcs.Stop = &stop
		}

		queue.Push(tcs)
	}

	return &TimecodeShow{
		tickRate: tickRate,
		queue:    queue,
	}, nil
}

type TimecodeShow struct {
	shutdownCh chan bool
	queue      *Queue
	tickRate   time.Duration
	start      time.Time
}

func (tc *TimecodeShow) Run(executor *Executor) error {
	ticker := time.NewTicker(tc.tickRate)
	tc.shutdownCh = make(chan bool)
	tc.start = time.Now()

	go func() {
		defer ticker.Stop()
		defer logger.Default.Debug("tcshow terminated")

		logger.Default.Debug("tcshow run")
		var terminating bool

		for {
			select {
			case <-tc.shutdownCh:
				tc.shutdownCh = nil
				logger.Default.Debug("tcshow terminating")
				return
			case <-ticker.C:
				now := time.Now()

				for {
					next := tc.queue.Peek()

					if next == nil {
						terminating = true
						tc.shutdownCh <- true
						break
					}

					if tc.start.Add(next.Start).Before(now) {
						next = tc.queue.Poll()

						var duration time.Duration
						if next.Stop != nil {
							duration = *next.Stop - next.Start
						}

						var sequences []ExecutorSequence
						for _, name := range next.Sequences {
							es := ExecutorSequence{
								Name:           name,
								BeatsPerMinute: next.BeatsPerMinute,
								Loop:           next.Loop,
							}
							if duration > 0 {
								es.Duration = &duration
							}
							sequences = append(sequences, es)
						}
						executor.EnableSequences(sequences)
					} else {
						break
					}
				}

				if !terminating {
					execTime := time.Since(now)
					if execTime > tc.tickRate {
						logger.Default.Warn("tcshow out of range",
							zap.Duration("actual", execTime),
							zap.Duration("expected", tc.tickRate),
						)
					}
				}
			}
		}
	}()
	return nil
}

func (tc *TimecodeShow) Stop(e *Executor) error {
	if tc.shutdownCh != nil {
		tc.shutdownCh <- true
	}
	return nil
}
