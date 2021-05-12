package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SetNick(command *discordgo.MessageCreate, s *discordgo.Session, userId string, newNick string) string {
	err := s.GuildMemberNickname(command.GuildID, "@me", newNick)

	if err != nil {
		fmt.Println(err)
	}

	return "Nick updated (I hope)!"
}
