package server

import (
	"github.com/gorilla/mux"
	temporal "go.temporal.io/sdk/client"
	"log"
	"net/http"
)

type Server struct {
	TemporalClient *temporal.Client
	Port           string
	LogDirectory   string
}

func (s *Server) RegisterHandlersAndServe() error {

	router := mux.NewRouter()

	router.HandleFunc("/warningmessage", s.CreateWarningMessage).Methods("POST")

	println("Server listening on port " + s.Port)
	err := http.ListenAndServe(":"+s.Port, router)
	if err != nil {
		log.Fatal("HTTP server could not be started", err)
	}
	return err
}
