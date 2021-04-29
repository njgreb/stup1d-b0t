package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/njgreb/stup1d-b0t/cache"
	"github.com/njgreb/stup1d-b0t/weather"
)

var token string = "ODM2NTg3OTY1MzkxMzA2NzUy.YIgLQg.zSdT2ej90-ELtqgXR6usA4vRSNo"

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

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

var s *discordgo.Session
var ctx = context.Background()

func init() {
	fmt.Println("init")
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

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
