package dicechallenge

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/games"
	"github.com/pdfkpb/gobo/pkg/patron"
)

func accept(patronDB *patron.PatronDB, params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
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
	challengerID := user.ID

	challengeAmount, err := patronDB.GetChallenge(challengerID)
	if err != nil {
		switch err {
		case patron.ErrChallengeNotFound:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":/ They didn't challenge you"))
		default:
			fmt.Printf("dicechallenge failed to CreateChallenge %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		}
		return
	}

	contenderID := m.Author.ID
	funds, err := patronDB.TakeFunds(contenderID, challengeAmount)
	if err != nil {
		switch err {
		case patron.ErrFundsCannotBeNeg:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, <@%s> only has %d monies", challengerID, funds))
			return
		default:
			fmt.Printf("dicechallenge failed to TakeFunds %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
	}

	one, _ := rand.Int(rand.Reader, big.NewInt(100))
	two, _ := rand.Int(rand.Reader, big.NewInt(100))
	challengerRoll := one.Int64() + 1
	contenderRoll := two.Int64() + 1

	var giveFunds string
	var winningRoll int
	var losingRoll int
	if challengerRoll == contenderRoll {
		_, err = patronDB.AddFunds(challengerID, games.TakeHouseCut(challengeAmount))
		if err != nil {
			fmt.Printf("dicechallenge failed to AddFunds %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
		_, err = patronDB.AddFunds(contenderID, games.TakeHouseCut(challengeAmount))
		if err != nil {
			fmt.Printf("dicechallenge failed to AddFunds %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
		err = patronDB.ClearChallenge(challengerID)
		if err != nil {
			fmt.Printf("dicechallenge failed to ClearChallenge %v\n", err)
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

	_, err = patronDB.AddFunds(giveFunds, games.TakeHouseCut(challengeAmount*2))
	if err != nil {
		fmt.Printf("dicechallenge failed to AddFunds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> won %d to %d your take is %d", giveFunds, winningRoll, losingRoll, games.TakeHouseCut(challengeAmount*2)))
}
