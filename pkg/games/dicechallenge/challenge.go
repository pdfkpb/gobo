package dicechallenge

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/games"
	"github.com/pdfkpb/gobo/pkg/patron"
)

func challenge(patronDB *patron.PatronDB, params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	match, err := regexp.Match("<@[0-9]{18}>", []byte(params[0]))
	if !match || err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	uid := params[0][2:20]
	user, err := s.User(uid)
	if err != nil || user == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", params[0]))
		return
	}
	contenderID := user.ID

	amount, err := strconv.Atoi(params[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Not a valid amount to bet")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	if amount < 0 {
		s.ChannelMessageSend(m.ChannelID, "You may only bet positive monies")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	challengerID := m.Author.ID
	funds, err := patronDB.CheckFunds(challengerID)
	if err != nil {
		fmt.Printf("dicechallenge:challenge failed to get user funds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if funds < amount {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, <@%s> only has %d monies", challengerID, funds))
		return
	}

	err = patronDB.CreateChallenge(challengerID, contenderID, amount)
	if err != nil {
		switch err {
		case patron.ErrChallengeAlreadyPosed:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have an outstanding challenge, you can remove it by `%s`", helpCancel))
		default:
			fmt.Printf("dicechallenge:challenge failed to CreateChallenge %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		}
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Thank you for you challenge, we charge a %.1f%% house cut you may cancel with `%s`", games.HouseCut*100, helpCancel))
}
