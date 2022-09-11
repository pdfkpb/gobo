package admin

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/gobo/pkg/bank"
)

const HelpRegister = "!register [@SomeUser]"

func RegisterUser(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	bankDB, err := bank.LoadBankDB()
	if err != nil {
		fmt.Printf("failed to load BankDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	registerID := m.Author.ID
	if len(params) > 0 {
		if !ChatPox.IsAdmin(m.Author.ID) {
			fmt.Printf("user <@%s> tried register %s\n", m.Author.ID, params[0])
			s.ChannelMessageSend(m.ChannelID, "Hey, knock it off")
			return
		}

		match, err := regexp.Match("<@[0-9]{18}>", []byte(params[0]))
		if !match || err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpRegister))
			return
		}

		uid := params[0][2:20]
		newUser, err := s.User(uid)
		if err != nil || newUser == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", params[0]))
			return
		}
		registerID = newUser.ID
	}

	err = bankDB.AddUser(registerID)
	if err != nil {
		fmt.Println("Oh Shit Me Boyo")
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has 1000 monies", registerID))
}
