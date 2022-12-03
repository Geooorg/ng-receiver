package server

import (
	"net/http"
)

func (s *Server) CreateWarningMessage(w http.ResponseWriter, r *http.Request) {

	println("received warning message")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
