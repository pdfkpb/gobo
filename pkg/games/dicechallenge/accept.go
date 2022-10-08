package dicechallenge

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/games"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

func accept(patronDB *patron.PatronDB, params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := params[0]
	if userID.Type() != commands.ParamTypeUserID {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Not a user id"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}

	challengerID := userID.UserID()
	user, err := s.User(string(challengerID))
	if err != nil || user == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s not found in this channel", challengerID.Mention()))
		return
	}

	challengeAmount, err := patronDB.GetChallenge(string(challengerID))
	if err != nil {
		switch err {
		case patron.ErrChallengeNotFound:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":/ They didn't challenge you"))
		default:
			fmt.Printf("dicechallenge:accept failed to CreateChallenge %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		}
		return
	}

	contenderID := userid.UserID(m.Author.ID)
	funds, err := patronDB.TakeFunds(string(contenderID), challengeAmount)
	if err != nil {
		switch err {
		case patron.ErrFundsCannotBeNeg:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, %s only has %d monies", challengerID.Mention(), funds))
			return
		default:
			fmt.Printf("dicechallenge:accept failed to TakeFunds %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
	}

	one, _ := rand.Int(rand.Reader, big.NewInt(100))
	two, _ := rand.Int(rand.Reader, big.NewInt(100))
	challengerRoll := one.Int64() + 1
	contenderRoll := two.Int64() + 1

	var giveFunds userid.UserID
	var winningRoll int
	var losingRoll int
	if challengerRoll == contenderRoll {
		_, err = patronDB.AddFunds(string(challengerID), games.TakeHouseCut(challengeAmount))
		if err != nil {
			fmt.Printf("dicechallenge:accept failed to AddFunds %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
		_, err = patronDB.AddFunds(string(contenderID), games.TakeHouseCut(challengeAmount))
		if err != nil {
			fmt.Printf("dicechallenge:accept failed to AddFunds %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
		err = patronDB.ClearChallenge(string(challengerID))
		if err != nil {
			fmt.Printf("dicechallenge:accept failed to ClearChallenge %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Y'all Tied, Congrats!"))
	} else if challengerRoll > contenderRoll {
		giveFunds = challengerID
		winningRoll = int(challengerRoll)
		losingRoll = int(contenderRoll)
	} else if challengerRoll < contenderRoll {
		giveFunds = contenderID
		winningRoll = int(contenderRoll)
		losingRoll = int(challengerRoll)
	}

	_, err = patronDB.AddFunds(string(giveFunds), games.TakeHouseCut(challengeAmount*2))
	if err != nil {
		fmt.Printf("dicechallenge:accept failed to AddFunds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s won %d to %d your take is %d", giveFunds.Mention(), winningRoll, losingRoll, games.TakeHouseCut(challengeAmount*2)))
}
