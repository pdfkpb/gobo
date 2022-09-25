package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/patron"
)

func BulkRegister(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !ChatPox.IsAdmin(m.Author.ID) {

	}

	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Printf("failed to load get members from channel: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	for _, member := range channel.Members {
		err = patronDB.AddUser(member.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to Register <@%s>: %s", member.ID, err.Error()))
		}
	}
}
