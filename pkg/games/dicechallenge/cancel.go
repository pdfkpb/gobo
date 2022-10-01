package dicechallenge

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

func cancel(patronDB *patron.PatronDB, params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	userID, err := userid.GetUserID(m.Author.Mention())
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to GetUserID %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	escrow, err := patronDB.GetChallenge(string(userID))
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to GetChallenge %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	_, err = patronDB.AddFunds(string(userID), escrow)
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to AddFunds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	err = patronDB.ClearChallenge(string(userID))
	if err != nil {
		fmt.Printf("dicechallenge:cancel failed to ClearChallenge %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprint("I've cleared your challenge"))
}
