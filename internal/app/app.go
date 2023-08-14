package app

import (
	"fmt"
	"log"
	"time"

	"github.com/realPointer/weather-bot/config"
	"github.com/realPointer/weather-bot/internal/weather"
	tele "gopkg.in/telebot.v3"
)

func initBot(cfg *config.Config) *tele.Bot {
	if cfg.WeatherToken == "" || cfg.TelegramToken == "" {
		log.Fatal("Please set telegram and weather tokens")
	}

	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func Run(cfg *config.Config) {
	b := initBot(cfg)

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnWeather := menu.Text("🌤 Weather by City")
	btnLocation := menu.Location("🌍 Weather by Location")

	menu.Reply(
		menu.Row(btnWeather),
		menu.Row(btnLocation),
	)

	replyMarkup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnCancel := replyMarkup.Text("Cancel")

	replyMarkup.Reply(
		replyMarkup.Row(btnCancel),
	)

	b.Handle("/start", func(c tele.Context) error {
		message := "Hello! 🤗\n\n" +
			"This bot can help you find out the weather for the city or your location\n" +
			"/weather <city> - Get the weather for a city\n" +
			"Or you can simply click on the buttons you need 😊\n\n" +
			"To pinpoint a city, you can enter a two- or three-digit country code separated by a comma\n" +
			"For example \"Moscow, RU\" or \"Moscow, USA\""
		return c.Send(message, menu)
	})

	b.Handle(&btnWeather, func(c tele.Context) error {
		return c.Send("Please enter a city name:", replyMarkup)
	})

	messageText := "Weather information for this location\n" +
		"───────────────────────────\n" +
		"🌡 Temperature: %v℃\n" +
		"🌤 Weather Status: %v (%v)\n" +
		"💨 Wind: %vmps\n" +
		"🌬 Pressure: %vmmHg\n" +
		"💧 Humidity: %v%%"

	b.Handle(tele.OnLocation, func(c tele.Context) error {
		weather, err := weather.GetWeatherByLocation(float64(c.Message().Location.Lat), float64(c.Message().Location.Lng), cfg.WeatherToken)
		if err != nil {
			return c.Send("Error getting weather")
		}

		err = c.Send(fmt.Sprintf(messageText, weather.Temperature, weather.WeatherStatus, weather.DescWeatherStatus, weather.WindSpeed, weather.Pressure, weather.Humidity))
		if err != nil {
			return err
		}

		return c.Send("Choose an option:", menu)
	})

	b.Handle("/weather", func(c tele.Context) error {
		city := c.Message().Payload

		if city == "" {
			return c.Send("Please enter a city name:", replyMarkup)
		}

		weather, err := weather.GetWeatherByCity(city, cfg.WeatherToken)
		if err != nil {
			return c.Send("Error getting weather")
		}

		err = c.Send(fmt.Sprintf(messageText, weather.Temperature, weather.WeatherStatus, weather.DescWeatherStatus, weather.WindSpeed, weather.Pressure, weather.Humidity))
		if err != nil {
			return err
		}

		return c.Send("Choose an option:", menu)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		if c.Message().Text == "Cancel" {
			return c.Send("Choose an option:", menu)
		}

		weather, err := weather.GetWeatherByCity(c.Message().Text, cfg.WeatherToken)
		if err != nil {
			return c.Send("Error getting weather")
		}

		err = c.Send(fmt.Sprintf(messageText, weather.Temperature, weather.WeatherStatus, weather.DescWeatherStatus, weather.WindSpeed, weather.Pressure, weather.Humidity))
		if err != nil {
			return err
		}

		return c.Send("Choose an option:", menu)
	})

	b.Start()
}
