package main

import (
	"html/template"
	"net/http"

	"go.uber.org/zap"
)

func returnStatus(w http.ResponseWriter) {
	s := NewMockStatus("mock status")
	t, err := template.ParseFiles("templates/template.html")
	if err != nil {
		logger.Debug("Unable to parse template",
			zap.Error(err),
		)
	}
	t.Execute(w, s)
}
