package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetWeatherFromEndpoint(city string) (WeatherResponse, error) {
	var response WeatherResponse
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=imperial&appid=7d49905c92261c5feac800e61b06f302", city))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	return response, nil
}

func GetNewForcast(city string) (WeatherResponse, error) {
	log.Printf("Making external request for city %s", city)

	response, err := GetWeatherFromEndpoint(city)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *Server) GetCachedForcast(city string) (WeatherResponse, error) {
	log.Printf("Getting results for city %s from cached memory", city)

	cachedForcast, err := s.RedisStore.GetForcastByCity(city)
	if err != nil {
		return WeatherResponse{}, err
	}

	return cachedForcast, nil
}
