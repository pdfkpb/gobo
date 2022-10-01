package dicechallenge

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/games"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

func challenge(patronDB *patron.PatronDB, params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := params[0]
	if userID.Type() != commands.ParamTypeUserID {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	user, err := s.User(string(userID.UserID()))
	if err != nil || user == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", userID.UserID().Mention()))
		return
	}
	contenderID := userID.UserID()

	amount := params[1]
	if amount.Type() != commands.ParamTypeInteger {
		s.ChannelMessageSend(m.ChannelID, "Not a valid amount to bet")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	amnt := amount.Integer()
	if amnt < 0 {
		s.ChannelMessageSend(m.ChannelID, "You may only bet positive monies")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	challengerID, err := userid.GetUserID(m.Author.Mention())
	if err != nil {
		fmt.Printf("dicechallenge:challenge failed to GetUserID %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	funds, err := patronDB.CheckFunds(string(challengerID))
	if err != nil {
		fmt.Printf("dicechallenge:challenge failed to get user funds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if funds < amnt {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, %s only has %d monies", challengerID.Mention(), funds))
		return
	}

	err = patronDB.CreateChallenge(string(challengerID), string(contenderID), amnt)
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
