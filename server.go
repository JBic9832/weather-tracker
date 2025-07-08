package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	RedisStore *Storage
}

type ServerError struct {
	Error string
}

type handlerFunction func(http.ResponseWriter, *http.Request) error

func makeHandlerFunc(f handlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			EncodeJSON(w, http.StatusBadRequest, ServerError{Error: err.Error()})

		}
	}
}

func EncodeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func NewServer(listenAddr string) *Server {
	rdb := NewStorage("localhost:6379")
	return &Server{
		listenAddr: listenAddr,
		RedisStore: rdb,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	log.Printf("Server is live at http://localhost%s", s.listenAddr)

	router.Handle("/weather/{city}", makeHandlerFunc(s.handleGetWeather))

	return http.ListenAndServe(s.listenAddr, router)
}
