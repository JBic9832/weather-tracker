package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) handleGetWeather(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	city := vars["city"]

	log.Println()

	exists, err := s.RedisStore.DoesEntryExist(city)
	if err != nil {
		return err
	}

	if exists {
		cachedData, err := s.GetCachedForcast(city)
		if err != nil {
			return err
		}

		return EncodeJSON(w, http.StatusOK, cachedData)
	}

	newData, err := GetNewForcast(city)
	if err != nil {
		return err
	}

	if err := s.RedisStore.StoreForcastByCity(newData); err != nil {
		return err
	}

	return EncodeJSON(w, http.StatusOK, newData)

}
