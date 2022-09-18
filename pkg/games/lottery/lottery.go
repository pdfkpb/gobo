package blackjack

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/patron"
)

const HelpPlay = "!roll"

func Play(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load PatronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if len(params) != 0 {
		s.ChannelMessageSend(m.ChannelID, "Incorrect number of parameters")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage: %s", HelpPlay))
		return
	}

	roll, _ := rand.Int(rand.Reader, big.NewInt(100))

	userID := m.Author.ID
	err = patronDB.SetLotteryRoll(userID, int(roll.Int64()+1))
	if err != nil {
		fmt.Printf("lottery failed to set users roll %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	currentWinner, currWinnerRoll, err := patronDB.GetLotteryWinner()
	if err != nil {

	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d, Current Winner is <@%s> with a %d", userID, roll, currentWinner, currWinnerRoll))
}

func ItsLotteryTime(s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load PatronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	winner, roll, err := patronDB.GetLotteryWinner()
	if err != nil {

	}

	currentFunds, err := patronDB.AddFunds(winner, 720)
	if err != nil {

	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you win! You now have %d", winner, roll, currentFunds))
}
