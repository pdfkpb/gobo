package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
)

var _ commands.Exec = Give

var HelpGive = "!give @SomeUser <some_amount>"

func Give(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !ChatPox.IsAdmin(m.Author.ID) {
		fmt.Printf("user %s tried to give funds: %v\n", m.Author.Username, ErrNotAdmin)
		s.ChannelMessageSend(m.ChannelID, "Hey, knock it off")
		return
	}

	if len(params) != 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usage: %s", HelpGive))
		return
	}

	userID := params[0]
	if userID.Type() != commands.ParamTypeUserID {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpGive))
		return
	}

	giveTo, err := s.User(string(userID.UserID()))
	if err != nil || giveTo == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", userID.UserID().Mention()))
		return
	}

	amount := params[1]
	if amount.Type() != commands.ParamTypeInteger {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a valid monies amount"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpGive))
		return
	}

	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	currentFunds, err := patronDB.AddFunds(giveTo.ID, amount.Integer())
	if err != nil {
		switch err {
		case patron.ErrInvalidAmount:
			s.ChannelMessageSend(m.ChannelID, err.Error())
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpGive))
		default:
			fmt.Printf("failed to add funds: %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %d monies", userID.UserID().Mention(), currentFunds))
}
