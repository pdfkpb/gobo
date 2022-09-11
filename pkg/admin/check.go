package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gobo/pkg/bank"
)

const HelpCheck = "!check"

func Check(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	bankDB, err := bank.LoadBankDB()
	if err != nil {
		fmt.Printf("failed to load BankDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	funds, err := bankDB.CheckFunds(m.Author.ID)
	if err != nil {
		switch err {
		case bank.ErrUserNotRegistered:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> is not registered, please do so by typing _!register_", m.Author.ID))
		default:
			fmt.Printf("failed to check funds for unknown reason: %v\n", err)
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has %d monies", m.Author.ID, funds))
}
