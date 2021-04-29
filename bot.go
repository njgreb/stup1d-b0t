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
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
	"github.com/njgreb/stup1d-b0t/weatherUtils"
)

var token string = "ODM2NTg3OTY1MzkxMzA2NzUy.YIgLQg.zSdT2ej90-ELtqgXR6usA4vRSNo"

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func setUserWeather(user string, location string) (string, error) {

	// store the value in redis
	err := rdb.Set(ctx, user, location, 0).Err()
	if err != nil {
		return "Failed to save preferred weather location.", err
	}

	weather, err := getWeather(location)

	return "Preferred weather location set, latest weather: " + weather, nil
}

func getWeather(location string) (string, error) {
	// see if we have the result in the cache
	val, err := rdb.Get(ctx, location).Result()
	if err == nil {
		return val, nil
	}

	// city name based search
	//res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + location + "&appid=84a719ec00c69a35d7821a0ae543b545&units=imperial")
	// zip code based search (USA)
	weatherUrl := "http://api.openweathermap.org/data/2.5/weather?zip=" + location + ",US&appid=84a719ec00c69a35d7821a0ae543b545&units=imperial"
	spew.Dump(weatherUrl)
	res, err := http.Get(weatherUrl)
	body, err := ioutil.ReadAll(res.Body)

	var weather_instance weatherUtils.Weather_main
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
	err = rdb.Set(ctx, location, return_string, 1*time.Minute).Err()
	if err != nil {
		return "", err
	}

	if err != nil {
		fmt.Println(err)
		return_string = "I derped yo"
	}

	return return_string, nil
}

func getPlainvilleWeather() (string, error) {
	return getWeather("Plainville,KS USA")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	//spew.Dump(m)

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		return
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
		return
	}

	if m.Content == "weather in gods country" {
		weather, err := getPlainvilleWeather()

		if err != nil {
			weather = "I failed to get Gods weather :("
		}

		s.ChannelMessageSend(m.ChannelID, weather)
		return
	}

	if m.Content == "should jesse ride his bike today" {
		s.ChannelMessageSend(m.ChannelID, "Yes...but he won't")
		return
	}

	if strings.HasPrefix(m.Content, "_weather set") || strings.HasPrefix(m.Content, "_w set") {
		commandParts := strings.Split(m.Content, " ")
		fmt.Println("Weather set for " + commandParts[2])
		message, err := setUserWeather(m.Author.Username, commandParts[2])

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "failed to set weather")
		}
		s.ChannelMessageSend(m.ChannelID, "weather set debug: "+message)
		return
	}

	if strings.HasPrefix(m.Content, "_weather") || strings.HasPrefix(m.Content, "_w") {
		commandParts := strings.Split(m.Content, " ")

		weatherLocation := ""

		if len(commandParts) == 1 {
			fmt.Println("Getting users prefered weather:" + m.Author.Username)
			// get the users preferred zip
			val, err := rdb.Get(ctx, m.Author.Username).Result()
			if err != nil || val == "" {
				s.ChannelMessageSend(m.ChannelID, "Failed to load weather without a zip code, use the command right or set a weather zip dork :( "+val)
				return
			}
			weatherLocation = val
			fmt.Println("weather loc is " + weatherLocation)
		} else {
			weatherLocation = commandParts[1]
		}

		fmt.Println("Weather for " + weatherLocation)

		weather, err := getWeather(weatherLocation)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to load weather :(")
		}
		s.ChannelMessageSend(m.ChannelID, weather)
		return
	}
}

var s *discordgo.Session
var ctx = context.Background()
var rdb *redis.Client

func init() {
	fmt.Println("init")
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.86.250:32768",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Connected to redis")

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
