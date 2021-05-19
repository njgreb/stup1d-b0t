package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adrg/postcode"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/davecgh/go-spew/spew"
	"github.com/njgreb/stup1d-b0t/cache"
)

var weatherApiToken string

func getToken() string {
	if len(strings.TrimSpace(weatherApiToken)) == 0 {
		weatherApiToken = os.Getenv("weatherApiToken")
	}

	return weatherApiToken
}

func getLatLong(location string) (string, string) {

	// see if we have the result in the cache
	val := cache.Get("get_lat_lon_" + location)
	if val != "" {
		parts := strings.Split(val, ",")
		fmt.Println("we have this location, we good")
		return parts[0], parts[1]
	}

	locationUrl := "http://api.openweathermap.org/geo/1.0/zip?zip=" + location + ",US&appid=" + getToken()
	spew.Dump(locationUrl)
	res, err := http.Get(locationUrl)

	if err != nil {
		fmt.Printf("oh crap, we failed to get the location")
	}

	body, err := ioutil.ReadAll(res.Body)

	spew.Dump(body)

	var loc encoded_location
	json.Unmarshal(body, &loc)
	spew.Dump(loc)

	// store the values in the cache
	cache.Set("get_lat_lon_"+location, fmt.Sprintf("%f", loc.Lat)+","+fmt.Sprintf("%f", loc.Lon), 0)

	return fmt.Sprintf("%f", loc.Lat), fmt.Sprintf("%f", loc.Lon)
}

func SetUserWeather(user string, location string) (string, discordgo.MessageEmbed, error) {

	// store the value in redis
	err := cache.Set(user, location, 0)
	if err != nil {
		embedMessage := embed.NewGenericEmbed("Error", "Failed to save preferred weather location")
		return "Failed to save preferred weather location.", *embedMessage, err
	}

	fmt.Printf("Location set, getting weather for %s, %s\n", user, location)
	weatherString, embedWeather, err := GetWeather(location)

	return "Preferred weather location set, latest weather:\n" + weatherString, embedWeather, nil
}

func GetWeather(location string) (string, discordgo.MessageEmbed, error) {
	// Verify a va)lid postal code
	if err := postcode.Validate(location); err != nil {
		embedMessage := embed.NewGenericEmbed("Error", "This command requires a valid US postal code at this time.")
		return "This command requires a valid US postal code at this time.", *embedMessage, nil
	}

	// get lat/lon of the location provided
	lat, lon := getLatLong(location)

	fmt.Printf("Loading weather for %s,%s\n", lat, lon)

	// see if we have the result in the cache
	var weather_instance one_call_weather
	weatherJson := cache.Get(lat + "," + lon)
	if weatherJson != "" {
		fmt.Printf("Weather found in cache\n")
		json.Unmarshal([]byte(weatherJson), &weather_instance)
	} else {
		// zip code based search (USA)
		//weatherUrl := "http://api.openweathermap.org/data/2.5/weather?zip=" + location + ",US&appid=" + getToken() + "&units=imperial"
		weatherUrl := "https://api.openweathermap.org/data/2.5/onecall?lat=" + lat + "&lon=" + lon + "&appid=" + getToken() + "&units=imperial&exclude=hourly,minutely"
		//spew.Dump(weatherUrl)
		fmt.Printf("Weather URL is: %s\n", weatherUrl)
		res, err := http.Get(weatherUrl)
		weatherJson, err := ioutil.ReadAll(res.Body)

		if err != nil {
		}

		json.Unmarshal(weatherJson, &weather_instance)
		//spew.Dump(body)
	}

	if len(weather_instance.Daily) == 0 {
		embedMessage := embed.NewGenericEmbed("Error", "Failed to load weather :(")
		return "Failed to load weather :(", *embedMessage, nil
	}

	windDirectionText := "West"

	switch {
	case weather_instance.Current.WindDeg > 270 || (weather_instance.Current.WindDeg >= 0 && weather_instance.Current.WindDeg < 45):
		windDirectionText = "North"
		break
	case weather_instance.Current.WindDeg >= 45 && weather_instance.Current.WindDeg < 135:
		windDirectionText = "East"
		break
	case weather_instance.Current.WindDeg >= 135 && weather_instance.Current.WindDeg < 215:
		windDirectionText = "South"
		break
	}

	return_string := fmt.Sprintf("%s, %.1fF | High: %.1fF | Low: %.1fF | Humidity: %d%% | Wind: %.1fmph @ %s (%d deg) | %s",
		weather_instance.Current.Weather[0].Description,
		weather_instance.Current.Temp,
		weather_instance.Daily[0].Temp.Max,
		weather_instance.Daily[0].Temp.Min,
		weather_instance.Current.Humidity,
		weather_instance.Current.WindSpeed,
		windDirectionText,
		weather_instance.Current.WindDeg,
		location)

	fmt.Println("here comes the weather baby!")
	//spew.Dump(weather_instance)
	fmt.Println("ok, thats all the weather")

	// store the value in redis
	err := cache.Set(location, string(weatherJson), 1*time.Minute)
	if err != nil {
		// should we do something here?
	}

	if err != nil {
		fmt.Println(err)
		return_string = "I derped yo"
	}

	embedOut := embed.NewEmbed().
		SetTitle(fmt.Sprintf("Weather for %s", location)).
		AddField("Current", fmt.Sprintf("%.1fF", weather_instance.Current.Temp)).
		AddField("High/Low", fmt.Sprintf("%.1fF/%1.fF", weather_instance.Daily[0].Temp.Max, weather_instance.Daily[0].Temp.Min)).
		AddField("Humidity", fmt.Sprintf("%d%%", weather_instance.Current.Humidity)).
		AddField("Wind", fmt.Sprintf("%.1fmph @ %s", weather_instance.Current.WindSpeed, windDirectionText)).
		SetFooter(weather_instance.Current.Weather[0].Description, fmt.Sprintf("http://openweathermap.org/img/wn/%s@2x.png", weather_instance.Current.Weather[0].Icon)).
		InlineAllFields().
		MessageEmbed

	return return_string, *embedOut, nil
}

// Open weather location struct
type encoded_location struct {
	Zip     string  `json:"zip,omitempty"`
	Name    string  `json:"name,omitempty"`
	Lat     float32 `json:"lat,omitempty"`
	Lon     float32 `json:"lon,omitempty"`
	Country string  `json:"country,omitempty"`
}

// Open Weather one call struct
type one_call_weather struct {
	Lat            float32         `json:"lat,omitempty"`
	Lon            float32         `json:"lon,omitempty"`
	Timezone       string          `json:"timezone,omitempty"`
	TimezoneOffset int             `json:"timezone_offset,omitempty"`
	Current        current_weather `json:"current,omitmepty"`
	Daily          []daily_weather
}

type current_weather struct {
	Dt         int            `json:"dt,omitempty"`
	Sunrise    int            `json:"sunrise,omitempty"`
	Sunset     int            `json:"sunset,omitempty"`
	Temp       float32        `json:"temp,omitempty"`
	FeelsLike  float32        `json:"feels_like,omitempty"`
	Pressure   int            `json:"pressure,omitempty"`
	Humidity   int            `json:"humidity,omitempty"`
	DewPoint   float32        `json:"dew_point,omitempty"`
	Uvi        float32        `json:"uvi,omitempty"`
	Clouds     int            `json:"clouds,omitempty"`
	Visibility int            `json:"visibility,omitempty"`
	WindSpeed  float32        `json:"wind_speed,omitempty"`
	WindDeg    int            `json:"wind_deg,omitempty"`
	Weather    []weather_desc `json:"weather,omitempty"`
}

type daily_weather struct {
	Dt        int            `json:"dt,omitempty"`
	Sunrise   int            `json:"sunrise,omitempty"`
	Sunset    int            `json:"sunset,omitempty"`
	Moonrise  int            `json:"moonrise,omitempty"`
	Moonset   int            `json:"moonset,omitempty"`
	MoonPhase float32        `json:"moonphase,omitempty"`
	Temp      temps          `json:"temp,omitempty"`
	FeelsLike temps          `json:"feels_like,omitempty"`
	Pressure  int            `json:"pressure,omitempty"`
	Humidity  int            `json:"humidity,omitempty"`
	DewPoint  float32        `json:"dew_poitn,omitempty"`
	WindSpeed float32        `json:"wind_speed,omitempty"`
	WindDeg   int            `json:"wind_deg,omitempty"`
	WindGust  float32        `json:"wind_gust,omitempty"`
	Weather   []weather_desc `json:"weather,omitempty"`
	Clouds    int            `json:"clouds,omitempty"`
	Pop       int            `json:"pop,omitempty"`
	Rain      float32        `json:"rain,omitempty"`
	Uvi       float32        `json:"uvi,omitempty"`
}

type temps struct {
	Day   float32 `json:"day,omitempty"`
	Min   float32 `json:"min,omitempty"`
	Max   float32 `json:"max,omitempty"`
	Night float32 `json:"night,omitempty"`
	Eve   float32 `json:"eve,omitempty"`
	Morn  float32 `json:"morn,omitempty"`
}

type weather_desc struct {
	Id          int    `json:"id,omitempty"`
	MainW       string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// old weather call struct
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
