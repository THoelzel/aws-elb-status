package main

import (
	"go.uber.org/zap"
	"time"
)

func startTicker() {
	for range time.NewTicker(time.Minute / 3).C {
		//go populateStatus()
	}
}

func initial() {
	//s := NewMockStatus("mock")
	s := NewStatus("mock")
	logger.Debug("tmp",
		zap.Any("val", s))
	go startTicker()
}
