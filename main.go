package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Location struct {
	Name string `json:"name"`
	Country string `json:"country"`
}

type Condition struct {
	Text string `json:"text"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
	Condition Condition `json:"condition"`
}

type Hour struct {
	TimeEpoch int64 `json:"time_epoch"`
	TempC float64 `json:"temp_c"`
	Condition Condition `json:"condition"`
	ChanceOfRain float64 `json:"chance_of_rain"`
}

type Forecastday struct {
	Hour []Hour `json:"hour"`
}

type Forecast struct {
	Forecastday []Forecastday `json:"forecastday"`
}

type Weather struct {
	Location Location `json:"location"`
	Current Current `json:"current"`
	Forecast Forecast `json:"forecast"`
}

const API_KEY string = "3be491824e234b96aeb121640241702"

func main() {

	q := "istanbul"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	apiUrl := "https://api.weatherapi.com/v1/forecast.json?key=" + API_KEY + "&q=" + q +"&days=1&aqi=no&alerts=no"

	res, err := http.Get(apiUrl)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("The weather api is not available!")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {

		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s   🌡️ %.0fC    🌧️ %.0f%%    %s",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Println(message)
		} else {
			color.Blue(message)
		}
	}
}