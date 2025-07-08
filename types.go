package main

type Forcast struct {
	Type        string `json:"main"`
	Description string `json:"description"`
}

type Main struct {
	Temp float32 `json:"temp"`
}

type Coord struct {
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}

type WeatherResponse struct {
	Temp        Main      `json:"main"`
	Weather     []Forcast `json:"weather"`
	City        string    `json:"name"`
	Coordinates Coord     `json:"coord"`
	ZipCode     string    `json:"zipcode"`
}
