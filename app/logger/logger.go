package logger

import (
	"go.uber.org/zap"
)

var (
	Default *zap.Logger
)

func init() {
	Default, _ = zap.NewProduction()
}
