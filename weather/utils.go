package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/njgreb/stup1d-b0t/cache"
)

func UpdateWeather() string {
	return "nothing"
}

func SetUserWeather(user string, location string) (string, error) {

	// store the value in redis
	err := cache.Set(user, location, 0)
	if err != nil {
		return "Failed to save preferred weather location.", err
	}

	weatherString, err := GetWeather(location)

	return "Preferred weather location set, latest weather: " + weatherString, nil
}

func GetWeather(location string) (string, error) {
	// see if we have the result in the cache
	val := cache.Get(location)
	if val != "" {
		return val, nil
	}

	// city name based search
	//res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + location + "&appid=84a719ec00c69a35d7821a0ae543b545&units=imperial")
	// zip code based search (USA)
	weatherUrl := "http://api.openweathermap.org/data/2.5/weather?zip=" + location + ",US&appid=84a719ec00c69a35d7821a0ae543b545&units=imperial"
	spew.Dump(weatherUrl)
	res, err := http.Get(weatherUrl)
	body, err := ioutil.ReadAll(res.Body)

	var weather_instance Weather_main
	json.Unmarshal(body, &weather_instance)
	spew.Dump(body)

	if len(weather_instance.Weather) == 0 {
		return "Failed to load weather :(", nil
	}

	fmt.Printf("%s\n", err)

	windDirectionText := "West"

	switch {
	case weather_instance.Wind.Deg > 270 || (weather_instance.Wind.Deg >= 0 && weather_instance.Wind.Deg < 45):
		windDirectionText = "North"
		break
	case weather_instance.Wind.Deg >= 45 && weather_instance.Wind.Deg < 135:
		windDirectionText = "East"
		break
	case weather_instance.Wind.Deg >= 135 && weather_instance.Wind.Deg < 215:
		windDirectionText = "South"
		break
	}

	return_string := fmt.Sprintf("%s, %.1fF | High: %.1fF | Low: %.1fF | Humidity: %d%% | Wind: %.1fmph @ %s (%d deg) | %s",
		weather_instance.Weather[0].MainW,
		weather_instance.Main.Temp,
		weather_instance.Main.TempMax,
		weather_instance.Main.TempMin,
		weather_instance.Main.Humidity,
		weather_instance.Wind.Speed,
		windDirectionText,
		weather_instance.Wind.Deg,
		weather_instance.Name)

	// store the value in redis
	err = cache.Set(location, return_string, 1*time.Minute)
	if err != nil {
		return "", err
	}

	if err != nil {
		fmt.Println(err)
		return_string = "I derped yo"
	}

	return return_string, nil
}

func GetPlainvilleWeather() (string, error) {
	return GetWeather("Plainville,KS USA")
}

type coord_weather struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

type weather_weather struct {
	Id          int    `json:"id,omitempty"`
	MainW       string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type main_weather struct {
	Temp      float64 `json:"temp,omitempty"`
	FeelsLike float64 `json:"feels_like,omitempty"`
	TempMin   float64 `json:"temp_min,omitempty"`
	TempMax   float64 `json:"temp_max,omitempty"`
	Pressure  int     `json:"pressure,omitempty"`
	Humidity  int     `json:"humidity,omitempty"`
}

type wind_weather struct {
	Speed float64 `json:"speed,omitempty"`
	Deg   int     `json:"deg,omitempty"`
	Gust  float64 `json:"gust,omitempty"`
}

type clouds_weather struct {
	All int `json:"all,omitempty"`
}

type sys_weather struct {
	TypeField int    `json:"type,omitempty"`
	Id        int    `json:"lid,omitempty"`
	Country   string `json:"country,omitempty"`
	Sunrise   int    `json:"sunrise,omitempty"`
	Sunset    int    `json:"sunset,omitempty"`
}

type Weather_main struct {
	Coord      coord_weather     `json:"coord"`
	Weather    []weather_weather `json:"weather"`
	Base       string            `json:"base"`
	Main       main_weather      `json:"main"`
	Visibility int               `json:"visibility"`
	Wind       wind_weather      `json:"wind"`
	Clouds     clouds_weather    `json:"clouds"`
	Dt         int               `json:"dt"`
	Sys        sys_weather       `json:"sys"`
	Timezone   int               `json:"timezone"`
	Id         int               `json:"id"`
	Name       string            `json:"name"`
	Cod        int               `json:"cod"`
}
