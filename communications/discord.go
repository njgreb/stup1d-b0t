package communications

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/njgreb/stup1d-b0t/cache"
	"github.com/njgreb/stup1d-b0t/version"
	"github.com/njgreb/stup1d-b0t/weather"
)

var discord *discordgo.Session

var BotToken string

func getSession() {
	flag.StringVar(&BotToken, "token", "", "Bot access token")

	if len(strings.TrimSpace(BotToken)) == 0 {
		BotToken = os.Getenv("botToken")
	}

	fmt.Println("Token is now:" + BotToken)

	var err error
	discord, err = discordgo.New("Bot " + BotToken)

	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		panic("Failed to connect to discord :(")
	}
}

func StartDiscord() {
	getSession()

	discord.AddHandler(MessageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}

func CloseDiscord() {
	discord.Close()
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		weather, err := weather.GetPlainvilleWeather()

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

	if strings.HasPrefix(m.Content, "_version") || strings.HasPrefix(m.Content, "_v") {
		s.ChannelMessageSend(m.ChannelID, "Currently running stup1d version "+version.Version)
		return
	}

	if strings.HasPrefix(m.Content, "_weather set") || strings.HasPrefix(m.Content, "_w set") {
		commandParts := strings.Split(m.Content, " ")
		fmt.Println("Weather set for " + commandParts[2])
		message, err := weather.SetUserWeather(m.Author.Username, commandParts[2])

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
			val := cache.Get(m.Author.Username)
			if val == "" {
				s.ChannelMessageSend(m.ChannelID, "Failed to load weather without a zip code, use the command right or set a weather zip dork :( "+val)
				return
			}
			weatherLocation = val
			fmt.Println("weather loc is " + weatherLocation)
		} else {
			weatherLocation = commandParts[1]
		}

		fmt.Println("Weather for " + weatherLocation)

		weather, err := weather.GetWeather(weatherLocation)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to load weather :(")
		}
		s.ChannelMessageSend(m.ChannelID, weather)
		return
	}
}
