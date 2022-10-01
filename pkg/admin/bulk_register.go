package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
)

var _ commands.Exec = (BulkRegister)(nil)

func BulkRegister(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !ChatPox.IsAdmin(m.Author.ID) {
		fmt.Printf("user %s tried to register e'eryone: %v\n", m.Author.Username, ErrNotAdmin)
		s.ChannelMessageSend(m.ChannelID, "Hey, knock it off")
		return
	}

	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	for _, member := range ChatPox.Members {
		err = patronDB.AddUser(member, 5000)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to Register <@%s>: %s", member, err.Error()))
		}
	}
}
