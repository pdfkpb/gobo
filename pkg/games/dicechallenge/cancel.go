package dicechallenge

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/patron"
)

func cancel(patronDB *patron.PatronDB, params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	escrow, err := patronDB.GetChallenge(m.Author.ID)
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to GetChallenge %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	_, err = patronDB.AddFunds(m.Author.ID, escrow)
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to AddFunds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	err = patronDB.ClearChallenge(m.Author.ID)
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to ClearChallenge %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprint("I've cleared your challenge"))
}
