package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

const BASE_URL = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

func resolveURL(BASE_URL string, location string, time string, API_KEY string) string {
	buf := bytes.NewBufferString(BASE_URL)
	buf.WriteString(location + "/")
	buf.WriteString(time)

	structure := "?unitGroup=metric&elements=datetime%2Ctemp%2Cfeelslike%2Chumidity%2Cprecipprob%2Cprecipcover%2Cpreciptype%2Cwindspeed%2Cpressure%2Cuvindex%2Cconditions%2Cdescription"

	buf.WriteString(structure)

	switch time {
	case "today":
		buf.WriteString("&include=hours")
	default:
		buf.WriteString("&include=daily")
	}

	buf.WriteString("&key=" + API_KEY)
	buf.WriteString("&contentType=json")

	return buf.String()
}

func APICall(url string, ch chan<- []byte, errCh chan<- error) {
	res, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errCh <- errors.New("failed call")
		return
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		errCh <- err
		return
	}
	ch <- data

}
