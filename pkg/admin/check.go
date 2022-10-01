package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

var _ commands.Exec = Check

const HelpCheck = "!check [@SomeUser]"

func Check(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	checkID, err := userid.GetUserID(m.Author.ID)
	if err != nil {
		fmt.Printf("failed to GetUserID: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if len(params) > 0 {
		userID := params[0]
		if userID.Type() != commands.ParamTypeUserID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpCheck))
			return
		}

		checkID = userID.UserID()
		checkFor, err := s.User(string(checkID))
		if err != nil || checkFor == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", checkID.Mention()))
			return
		}
	}

	funds, err := patronDB.CheckFunds(string(checkID))
	if err != nil {
		switch err {
		case patron.ErrUserNotRegistered:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not registered, please do so by typing _!register_", checkID.Mention()))
		default:
			fmt.Printf("failed to check funds for unknown reason: %v\n", err)
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %d monies", checkID.Mention(), funds))
}
