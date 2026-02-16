package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() error {
	var err error

	config := zap.NewProductionConfig()
	config.Encoding = "json" // structured logging
	config.OutputPaths = []string{"stdout"}

	Log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	_ = Log.Sync()
}
