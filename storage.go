package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	RedisDB *redis.Client
	Context context.Context
}

func NewStorage(addr string) *Storage {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	return &Storage{
		RedisDB: redisClient,
		Context: ctx,
	}
}

func FormatCity(city string) string {
	s := strings.ToLower(city)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "+", "")

	return s

}

func (s *Storage) DoesEntryExist(cityName string) (bool, error) {
	city := FormatCity(cityName)
	log.Println("Checking for city in cache:", city)
	exists, err := s.RedisDB.Exists(s.Context, city).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (s *Storage) StoreForcastByCity(forcast WeatherResponse) error {
	forcastJson, err := json.Marshal(forcast)
	if err != nil {
		return err
	}

	city := forcast.ZipCode
	log.Printf("Caching result for: %s\n", city)

	if err := s.RedisDB.Set(s.Context, city, forcastJson, 20*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetForcastByCity(cityZip string) (WeatherResponse, error) {
	log.Printf("Getting cached result for: %s\n", cityZip)

	var weatherData WeatherResponse
	data, err := s.RedisDB.Get(s.Context, cityZip).Bytes()
	if err != nil {
		return weatherData, err
	}

	err = json.Unmarshal(data, &weatherData)
	if err != nil {
		return weatherData, err
	}

	return weatherData, nil
}
