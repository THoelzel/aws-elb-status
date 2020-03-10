package main

import "time"

func startTicker() {
	for range time.NewTicker(time.Minute / 3).C {
		//go populateStatus()
	}
}

func initial() {
	//s := NewMockStatus("mock")
	go startTicker()
}
