package main

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func main() {

	b = NewBroker()
	go b.Start()
	//initial()

	config := NewConfig()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/v1/ws", wsHandler)
	port := fmt.Sprintf(":%d", config.GetPort())

	logger.Fatal("Unable to start server",
		zap.Error(http.ListenAndServe(port, nil)),
	)
}
