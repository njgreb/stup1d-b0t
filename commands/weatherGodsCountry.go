package commands

import "github.com/njgreb/stup1d-b0t/weather"

func WeatherGodsCountry(command string) string {
	weather, err := weather.GetPlainvilleWeather()

	if err != nil {
		weather = "I failed to get Gods weather :("
	}

	return weather
}
