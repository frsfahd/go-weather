package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var API_KEY = os.Getenv("API_KEY")

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// r.Use(gin.Logger())

	r.GET("/today/:location", s.GetWeatherToday)
	r.GET("/oneWeek/:location", s.GetWeatherOneWeek)

	r.GET("/client", s.CheckClient)

	r.GET("/test", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	return r
}

func (s *Server) GetWeatherToday(c *gin.Context) {
	location := c.Param("location")
	key := fmt.Sprintf("%s:today", location)

	url := resolveURL(BASE_URL, location, "today", API_KEY)

	dataCh := make(chan []byte)
	errCh := make(chan error)
	var externalData ExternalWeatherResponse
	var weather WeatherResponse
	var isCached bool = true

	d, err := s.db.GetString(key)
	if err != nil {
		//check cache first
		if errors.Is(err, redis.Nil) {
			isCached = false
			slog.Info("gotcha")
		} else {
			c.Error(err)

		}
	}

	// Get response data from cache
	if isCached {
		err = json.Unmarshal([]byte(d), &externalData)
		if err != nil {
			c.Error(err)
		}
		weather = WeatherResponse{Location: Location{Latitude: externalData.Latitude, Longitude: externalData.Longitude, Address: externalData.ResolvedAddress, Timezone: externalData.Timezone}, Days: externalData.Days}
		// response: 200 (from redis cache)
		c.IndentedJSON(http.StatusOK, Response{Message: "ok", Data: weather})
		return
	}

	// Get response data from API
	go APICall(url, dataCh, errCh)
	select {
	case data := <-dataCh:
		err := s.db.SetString(key, data)
		if err != nil {
			c.Error(err)
		}
		err = json.Unmarshal(data, &externalData)
		if err != nil {
			c.Error(err)
		}

		weather = WeatherResponse{Location: Location{Latitude: externalData.Latitude, Longitude: externalData.Longitude, Address: externalData.ResolvedAddress, Timezone: externalData.Timezone}, Days: externalData.Days}
		// response: 200 (directly from API)
		c.IndentedJSON(http.StatusOK, Response{Message: "ok", Data: weather})
	case err := <-errCh:
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "server error"})
	}

}

func (s *Server) GetWeatherOneWeek(c *gin.Context) {
	location := c.Param("location")
	key := fmt.Sprintf("%s:1Week", location)

	url := resolveURL(BASE_URL, location, "next7days", API_KEY)

	dataCh := make(chan []byte)
	errCh := make(chan error)
	var externalData ExternalWeatherResponse
	var weather WeatherResponse
	var isCached bool = true

	d, err := s.db.GetString(key)
	if err != nil {
		//check cache first
		if errors.Is(err, redis.Nil) {
			isCached = false
			slog.Info("gotcha")
		} else {
			c.Error(err)

		}
	}

	// Get response data from cache
	if isCached {
		err = json.Unmarshal([]byte(d), &externalData)
		if err != nil {
			c.Error(err)
		}
		weather = WeatherResponse{Location: Location{Latitude: externalData.Latitude, Longitude: externalData.Longitude, Address: externalData.ResolvedAddress, Timezone: externalData.Timezone}, Days: externalData.Days}
		// response: 200 (from redis cache)
		c.IndentedJSON(http.StatusOK, Response{Message: "ok", Data: weather})
		return
	}

	// Get response data from API
	go APICall(url, dataCh, errCh)
	select {
	case data := <-dataCh:
		err := s.db.SetString(key, data)
		if err != nil {
			c.Error(err)
		}
		err = json.Unmarshal(data, &externalData)
		if err != nil {
			c.Error(err)
		}

		weather = WeatherResponse{Location: Location{Latitude: externalData.Latitude, Longitude: externalData.Longitude, Address: externalData.ResolvedAddress, Timezone: externalData.Timezone}, Days: externalData.Days}
		// response: 200 (directly from API)
		c.IndentedJSON(http.StatusOK, Response{Message: "ok", Data: weather})
	case err := <-errCh:
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: "server error"})
	}
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "ok"

	c.JSON(http.StatusOK, resp)
}
func (s *Server) CheckClient(c *gin.Context) {
	resp := make(map[string]string)
	resp["remoteIP"] = c.RemoteIP()
	resp["clientIP"] = c.ClientIP()
	c.IndentedJSON(http.StatusOK, resp)
}
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
