package server

import (
	"encoding/json"
	"log"
	"time"
)

// Custom time type to handle the datetime format
type CustomTime struct{ time.Time }

// UnmarshalJSON implements the json.Unmarshaler interface for CustomTime
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	// Remove quotes from the JSON string
	var s string
	json.Unmarshal(b, &s)
	// Parse time using the specified format
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	ct.Time = t
	log.Println(ct.Time)
	log.Println(s)
	return nil
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	Timezone  string  `json:"timezone"`
}

type CoreWeather struct {
	Datetime      string   `json:"datetime"`
	DatetimeEpoch int      `json:"datetimeEpoch"`
	Temp          float64  `json:"temp"`
	Feelslike     float64  `json:"feelslike"`
	Humidity      float64  `json:"humidity"`
	Preciprob     float64  `json:"precipprob"`
	Precipcover   float64  `json:"precipcover"`
	Preciptype    []string `json:"preciptype"`
	Windspeed     float64  `json:"windspeed"`
	Pressure      float64  `json:"pressure"`
	Uvindex       float64  `json:"uvindex"`
	Conditions    string   `json:"conditions"`
	Descpription  string   `json:"description"`
}

type Weather struct {
	CoreWeather
	Hours []interface{} `json:"hours,omitempty"`
}

type WeatherResponse struct {
	Location
	Days []Weather `json:"days"`
}

type ExternalWeatherResponse struct {
	QueryCost       int     `json:"queryCost"`
	ResolvedAddress string  `json:"resolvedAddress"`
	Tzoffset        float64 `json:"tzoffset"`
	Location
	Days []Weather `json:"days"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
