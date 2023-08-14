package weather

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/tidwall/gjson"
)

var ErrLocationNotFound = errors.New("location not found")

type Weather struct {
	Temperature       float64
	WeatherStatus     string
	DescWeatherStatus string
	WindSpeed         float64
	Pressure          float64
	Humidity          int64
}

func GetWeatherByCity(city, apiKey string) (*Weather, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%v&appid=%v&units=metric", city, apiKey)
	return getWeather(url)
}

func GetWeatherByLocation(lat, long float64, apiKey string) (*Weather, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%v&units=metric", lat, long, apiKey)
	return getWeather(url)
}

func getWeather(url string) (*Weather, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	weatherJSON := string(body)

	if gjson.Get(weatherJSON, "cod").String() == "200" {
		weather := &Weather{
			Temperature:       gjson.Get(weatherJSON, "main.temp").Float(),
			WeatherStatus:     gjson.Get(weatherJSON, "weather.0.main").String(),
			DescWeatherStatus: gjson.Get(weatherJSON, "weather.0.description").String(),
			WindSpeed:         gjson.Get(weatherJSON, "wind.speed").Float(),
			Pressure:          math.Round(gjson.Get(weatherJSON, "main.pressure").Float()*0.750064*100) / 100,
			Humidity:          gjson.Get(weatherJSON, "main.humidity").Int(),
		}

		return weather, nil
	}

	return &Weather{}, fmt.Errorf("location not found: %w", ErrLocationNotFound)
}
