package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetWeatherFromEndpoint(city string) (WeatherResponse, error) {
	var response WeatherResponse
	response.ZipCode = city
	apiKey := os.Getenv("WeatherAPIKey")
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&units=imperial&appid=%s", city, apiKey))
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

func GetNewForcast(cityZip string) (WeatherResponse, error) {
	log.Printf("Making external request for city %s", cityZip)

	response, err := GetWeatherFromEndpoint(cityZip)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *Server) GetCachedForcast(city string) (WeatherResponse, error) {
	cachedForcast, err := s.RedisStore.GetForcastByCity(city)
	if err != nil {
		return WeatherResponse{}, err
	}

	return cachedForcast, nil
}
