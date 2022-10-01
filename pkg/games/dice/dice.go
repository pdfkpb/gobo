package dice

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

var _ commands.Exec = Dice

const HelpPlay = "O'er Under:\n\t!dice <amount> over | under"

func Dice(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
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

	amount := params[0]
	if amount.Type() != commands.ParamTypeInteger {
		s.ChannelMessageSend(m.ChannelID, "Not a valid amount to bet")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	amnt := amount.Integer()
	if amnt < 0 {
		s.ChannelMessageSend(m.ChannelID, "You may only bet positive monies")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	userID := userid.UserID(m.Author.ID)
	funds, err := patronDB.CheckFunds(string(userID))
	if err != nil {
		fmt.Printf("dice failed to get user funds %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if funds < amnt {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, %s only has %d monies", userID.Mention(), funds))
		return
	}

	one, _ := rand.Int(rand.Reader, big.NewInt(6))
	two, _ := rand.Int(rand.Reader, big.NewInt(6))
	n := one.Int64() + two.Int64() + 2

	overUnder := params[1]
	if overUnder.Type() != commands.ParamTypeString {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("either over or under"))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	overOrUnder := overUnder.String()
	if overOrUnder == "over" && n > 7 || overOrUnder == "under" && n < 7 {
		currentFunds, err := patronDB.AddFunds(m.Author.ID, amnt)
		if err != nil {

		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s rolled a %d you win! You now have %d", userID.Mention(), n, currentFunds))
	} else {
		currentFunds, err := patronDB.TakeFunds(string(userID), amnt)
		if err != nil {

		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s rolled a %d you lose ¯\\_(ツ)_/¯ you still have %d", userID.Mention(), n, currentFunds))
	}
}
