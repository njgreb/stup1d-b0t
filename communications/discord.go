package communications

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/davecgh/go-spew/spew"
	"github.com/njgreb/stup1d-b0t/cache"
	"github.com/njgreb/stup1d-b0t/commands"
	"github.com/njgreb/stup1d-b0t/gif"
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

type CommandResponse struct {
	simpleMessage string
	embedMessage  *discordgo.MessageEmbed
	messageType   int // 0 = simple, 1 = embed
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

	var response string
	var commandResponse CommandResponse
	processed := false

	// Show help message
	if m.Content == CommandPrefix+"help" || m.Content == CommandPrefix+"h" {
		commandResponse.messageType = 0
		commandResponse.simpleMessage = `
_w set ##### (US Zip code) to set your weather location
_w to see your weather
_w ##### (US Zip code) to see weather somewhere in the US
		`
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"config") {
		commandParts := strings.Split(m.Content, " ")

		if commandParts[1] == "gifcontentfilter" {
			commandResponse.messageType = 0
			if gif.SetFilterLevel(m.GuildID, commandParts[2]) {
				commandResponse.simpleMessage = "tenor filtering set to " + commandParts[2]
			} else {
				commandResponse.simpleMessage = "Invalid content filter sent. Please use: off, low, medium, high"
			}
		}
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"nameme") {
		commandParts := strings.Split(m.Content, " ")
		newNick := strings.Join(commandParts[1:], " ")
		commandResponse.messageType = 0
		commandResponse.simpleMessage = commands.SetNick(m, s, m.Author.ID, newNick)
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"gif") {
		commandParts := strings.Split(m.Content, " ")
		gifSearch := strings.Join(commandParts[1:], " ")
		commandResponse.messageType = 0
		commandResponse.simpleMessage = commands.Gif(gifSearch, false, m.GuildID)
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"version") || strings.HasPrefix(m.Content, CommandPrefix+"v") {
		commandResponse.messageType = 1

		embedOut := embed.NewEmbed().
			SetTitle("stup1d-b0t info").
			AddField("Version", version.Version).
			SetImage(commands.Gif("bots", true, m.GuildID)).
			MessageEmbed

		commandResponse.embedMessage = embedOut
		fmt.Print("outputing version")
	}

	if strings.HasPrefix(m.Content, CommandPrefix+"weather set") || strings.HasPrefix(m.Content, CommandPrefix+"w set") {
		commandParts := strings.Split(m.Content, " ")
		fmt.Println("Weather set for " + commandParts[2])
		message, embedWeather, err := weather.SetUserWeather(m.Author.ID, commandParts[2])

		if err != nil {
			commandResponse.messageType = 0
			response = "failed to set weather"
		} else {
			response = message
			commandResponse.messageType = 1
			commandResponse.embedMessage = &embedWeather
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
				response = "Please set a preferred weather locationed with .w set #####"
			} else {
				weatherLocation = val
				fmt.Println("weather loc is " + weatherLocation)
			}
		} else {
			if commandParts[1] == "clear" {
				cache.Set(m.Author.ID, "", -1)
				response = "Preferred weather location cleared."
			} else {
				weatherLocation = commandParts[1]
			}
		}

		if weatherLocation != "" {
			fmt.Println("Weather for " + weatherLocation)

			weather, weatherEmbed, err := weather.GetWeather(weatherLocation)
			commandResponse.embedMessage = &weatherEmbed

			if err != nil {
				response = "Failed to load weather :("
			}
			response = weather
		}

		commandResponse.messageType = 0
		if commandResponse.embedMessage != nil {
			commandResponse.messageType = 1
		}

		commandResponse.simpleMessage = response
	}

	// extended forecast!
	if processed == false && (strings.HasPrefix(m.Content, CommandPrefix+"weather long") || strings.HasPrefix(m.Content, CommandPrefix+"wl")) {
		commandParts := strings.Split(m.Content, " ")

		weatherLocation := ""

		if len(commandParts) == 1 {
			fmt.Println("Getting users prefered weather:" + m.Author.Username)
			// get the users preferred zip
			val := cache.Get(m.Author.ID)
			if val == "" {
				response = "Please set a preferred weather locationed with .w set #####"
			} else {
				weatherLocation = val
				fmt.Println("weather loc is " + weatherLocation)
			}
		} else {
			if commandParts[1] == "clear" {
				cache.Set(m.Author.ID, "", -1)
				response = "Preferred weather location cleared."
			} else {
				weatherLocation = commandParts[1]
			}
		}

		if weatherLocation != "" {
			fmt.Println("Weather for " + weatherLocation)

			weather, weatherEmbed, err := weather.GetWeatherLong(weatherLocation)
			commandResponse.embedMessage = &weatherEmbed

			if err != nil {
				response = "Failed to load weather :("
			}
			response = weather
		}

		commandResponse.messageType = 0
		if commandResponse.embedMessage != nil {
			commandResponse.messageType = 1
		}

		commandResponse.simpleMessage = response
	}

	if commandResponse.messageType == 0 {
		s.ChannelMessageSend(m.ChannelID, commandResponse.simpleMessage)
	} else {
		spew.Dump(commandResponse)
		s.ChannelMessageSendEmbed(m.ChannelID, commandResponse.embedMessage)
	}

}
