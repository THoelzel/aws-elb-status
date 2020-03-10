package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	lc := zap.NewDevelopmentConfig()
	lc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var err error
	logger, err = lc.Build()

	if err != nil {
		fmt.Println("Unable to initiate logger")
		return
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()
}
