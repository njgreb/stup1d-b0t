package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SetNick(command *discordgo.MessageCreate, s *discordgo.Session, userId string, newNick string) string {
	fmt.Printf("Updating nick: %s, %s, %s", command.GuildID, userId, newNick)

	err := s.GuildMemberNickname(command.GuildID, "@me", newNick)

	if err != nil {
		fmt.Println(err)
	}

	return "Nick updated (I hope)!"
}
