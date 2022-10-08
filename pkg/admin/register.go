package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

var _ commands.Exec = RegisterUser

const HelpRegister = "!register [@SomeUser]"

func RegisterUser(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	registerID := userid.UserID(m.Author.ID)
	if err != nil {
		fmt.Printf("failed to GetUserID: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if len(params) > 0 {
		if !ChatPox.IsAdmin(m.Author.ID) {
			fmt.Printf("user %s tried register %s\n", registerID, params[0].UserID().Mention())
			s.ChannelMessageSend(m.ChannelID, "Hey, knock it off")
			return
		}

		userID := params[0]
		if userID.Type() != commands.ParamTypeUserID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpRegister))
			return
		}

		registerID = userID.UserID()
		newUser, err := s.User(string(registerID))
		if err != nil || newUser == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", registerID.Mention()))
			return
		}
	}

	err = patronDB.AddUser(string(registerID), 1000)
	if err != nil {
		fmt.Printf("failed to AddUser: %v\n", err)
		switch err {
		case patron.ErrorUserAlreadyRegistered:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is already Registered", registerID.Mention()))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Check your balance with: %s", HelpCheck))
		default:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has 1000 monies", registerID.Mention()))
}
