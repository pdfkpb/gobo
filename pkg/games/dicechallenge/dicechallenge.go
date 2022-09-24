package dicechallenge

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/games"
	"github.com/pdfkpb/gobo/pkg/patron"
)

var _ games.Play = Play

const (
	helpChallenge = "!dc @SomePlayer <amount>"
	helpAccept    = "!dc @SomePlayer"
)

var (
	HelpPlay = fmt.Sprintf("Dice Challenge:\n\tChallenge: %s\n\tAccept: %s", helpChallenge, helpAccept)
)

func Play(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if len(params) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Incorrect number of parameters")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	amt, err := strconv.Atoi(params[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Not a valid amount to bet")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	if amt < 0 {
		s.ChannelMessageSend(m.ChannelID, "You may only bet positive monies")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	userID := m.Author.ID
	funds, err := patronDB.CheckFunds(userID)
	if err != nil {
		fmt.Printf("dice failed to get user funds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if funds < amt {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, <@%s> only has %d monies", userID, funds))
		return
	}

	one, _ := rand.Int(rand.Reader, big.NewInt(6))
	two, _ := rand.Int(rand.Reader, big.NewInt(6))
	n := one.Int64() + two.Int64() + 2

	overOrUnder := strings.ToLower(params[1])
	if overOrUnder != "over" && overOrUnder != "under" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("either over or under"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	if overOrUnder == "over" && n > 7 || overOrUnder == "under" && n < 7 {
		currentFunds, err := patronDB.AddFunds(userID, amt)
		if err != nil {

		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you win! You now have %d", userID, n, currentFunds))
	} else {
		currentFunds, err := patronDB.TakeFunds(userID, amt)
		if err != nil {

		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you lose ¯\\_(ツ)_/¯ you still have %d", userID, n, currentFunds))
	}
}
