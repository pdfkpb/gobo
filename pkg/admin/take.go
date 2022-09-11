package admin

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/gobo/pkg/bank"
)

const (
	HelpTake = "!take @SomeUser <some_amount>"
)

func Take(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if !ChatPox.IsAdmin(m.Author.ID) {
		fmt.Printf("user %s tried to take funds: %v\n", m.Author.Username, ErrNotAdmin)
		s.ChannelMessageSend(m.ChannelID, "Hey, knock it off")
		return
	}

	if len(params) != 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usage: %s", HelpTake))
		return
	}

	match, err := regexp.Match("<@[0-9]{18}>", []byte(params[0]))
	if !match || err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("  Usage: %s", HelpTake))
		return
	}

	uid := params[0][2:20]
	giveTo, err := s.User(uid)
	if err != nil || giveTo == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", params[0]))
		return
	}

	amount, err := strconv.Atoi(params[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a valid monies amount"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("  Usage: %s", HelpTake))
		return
	}

	bankDB, err := bank.LoadBankDB()
	if err != nil {
		fmt.Printf("failed to load BankDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	currentFunds, err := bankDB.AddFunds(uid, amount)
	if err != nil {
		fmt.Printf("failed to add funds: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %d monies", params[0], currentFunds))
}
