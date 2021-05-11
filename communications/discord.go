package communications

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/njgreb/stup1d-b0t/cache"
	"github.com/njgreb/stup1d-b0t/commands"
	"github.com/njgreb/stup1d-b0t/version"
	"github.com/njgreb/stup1d-b0t/weather"
)

var discord *discordgo.Session

var (
	BotToken      string
	CommandPrefix string
	BotUserId     string
)

type discordCommand interface {
	parse() string
	run() string
}

func getSession() {
	BotToken = os.Getenv("botToken")
	CommandPrefix = os.Getenv("commandPrefix")

	fmt.Println("Command Prefix is " + CommandPrefix)

	var err error
	discord, err = discordgo.New("Bot " + BotToken)

	botUser, err := discord.User("@me")
	BotUserId = botUser.ID

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
	var response string
	processed := false

	// Show help message
	if m.Content == CommandPrefix+"help" || m.Content == CommandPrefix+"h" {
		response = `
_w set ##### (US Zip code) to set your weather location
_w to see your weather
_w ##### (US Zip code) to see weather somewhere in the US
		`
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"nameme") {
		commandParts := strings.Split(m.Content, " ")
		newNick := strings.Join(commandParts[1:], " ")
		response = commands.SetNick(m, s, m.Author.ID, newNick)
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == CommandPrefix+"ping" {
		response = commands.Pong(m.Content)
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == CommandPrefix+"pong" {
		response = commands.Ping(m.Content)
	}

	if m.Content == "weather in gods country" {
		response = commands.WeatherGodsCountry(m.Content)
	}

	if m.Content == "should jesse ride his bike today" {
		response = "Yes...but he won't"
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"version") || strings.HasPrefix(m.Content, CommandPrefix+"v") {
		response = "Currently running stup1d version " + version.Version
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"weather set") || strings.HasPrefix(m.Content, CommandPrefix+"w set") {
		commandParts := strings.Split(m.Content, " ")
		fmt.Println("Weather set for " + commandParts[2])
		message, err := weather.SetUserWeather(m.Author.ID, commandParts[2])

		if err != nil {
			response = "failed to set weather"
		} else {
			response = message
		}
		processed = true
	}

	if processed == false && (strings.HasPrefix(m.Content, CommandPrefix+"weather") || strings.HasPrefix(m.Content, CommandPrefix+"w")) {
		commandParts := strings.Split(m.Content, " ")

		weatherLocation := ""

		if len(commandParts) == 1 {
			fmt.Println("Getting users prefered weather:" + m.Author.Username)
			// get the users preferred zip
			val := cache.Get(m.Author.ID)
			if val == "" {
				response = "Failed to load weather without a zip code, use the command right or set a weather zip dork :( " + val
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
			response = "Failed to load weather :("
		}
		response = weather
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
