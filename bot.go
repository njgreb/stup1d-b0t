package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
)

var token string = "ODM2NTg3OTY1MzkxMzA2NzUy.YIgLQg.zSdT2ej90-ELtqgXR6usA4vRSNo"

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

type weather struct {
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

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func getWeather(location string) string {
	// city name based search
	//res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + location + "&appid=84a719ec00c69a35d7821a0ae543b545&units=imperial")
	// zip code based search (USA)
	weatherUrl := "http://api.openweathermap.org/data/2.5/weather?zip=" + location + ",US&appid=84a719ec00c69a35d7821a0ae543b545&units=imperial"
	spew.Dump(weatherUrl)
	res, err := http.Get(weatherUrl)
	body, err := ioutil.ReadAll(res.Body)
	var weather_instance weather
	json.Unmarshal(body, &weather_instance)
	spew.Dump(body)

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

	if err != nil {
		fmt.Println(err)
		return_string = "I derped yo"
	}

	return return_string
}

func getPlainvilleWeather() string {
	return getWeather("Plainville,KS USA")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	spew.Dump(m)

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "weather in gods country" {
		s.ChannelMessageSend(m.ChannelID, getPlainvilleWeather())
	}

	if m.Content == "should jesse ride his bike today" {
		s.ChannelMessageSend(m.ChannelID, "Yes...but he won't")
	}

	if strings.Contains(m.Content, "weather") {
		commandParts := strings.Split(m.Content, " ")
		fmt.Println("Weather for " + commandParts[1])
		s.ChannelMessageSend(m.ChannelID, getWeather(commandParts[1]))
	}
}

var s *discordgo.Session
var ctx = context.Background()

func init() {
	fmt.Println("init")
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.86.250:32768",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}

func main() {
	fmt.Println("Hello, World!")

	discord, err := discordgo.New("Bot " + token)

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()

}
