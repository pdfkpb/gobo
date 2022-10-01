package admin

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
)

var _ commands.Exec = (Check)(nil)

const HelpCheck = "!check [@SomeUser]"

func Check(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	checkID := m.Author.ID
	if len(params) > 0 {
		if !ChatPox.IsAdmin(m.Author.ID) {
			fmt.Printf("user <@%s> tried check %s\n", m.Author.ID, params[0])
			s.ChannelMessageSend(m.ChannelID, "Hey, knock it off")
			return
		}

		match, err := regexp.Match("<@[0-9]{18}>", []byte(params[0]))
		if !match || err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpCheck))
			return
		}

		uid := params[0][2:20]
		checkFor, err := s.User(uid)
		if err != nil || checkFor == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", params[0]))
			return
		}
		checkID = checkFor.ID
	}

	funds, err := patronDB.CheckFunds(checkID)
	if err != nil {
		switch err {
		case patron.ErrUserNotRegistered:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> is not registered, please do so by typing _!register_", checkID))
		default:
			fmt.Printf("failed to check funds for unknown reason: %v\n", err)
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has %d monies", checkID, funds))
}
