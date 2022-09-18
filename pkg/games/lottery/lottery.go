package lottery

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/patron"
)

const HelpPlay = "!roll"
const payout = 720

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
	realRoll := int(roll.Int64() + 1)

	userID := m.Author.ID
	err = patronDB.SetLotteryRoll(userID, realRoll)
	if err != nil {
		switch err {
		case patron.ErrAlreadyLotteryRolled:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Hey, knock it off, you already rolled"))
		default:
			fmt.Printf("lottery failed to set users roll %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
		}
		return
	}

	currentWinners, currWinnerRoll, err := patronDB.GetLotteryWinner()
	if err != nil {

	}

	switch len(currentWinners) {
	case 1:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d, Current Winner is <@%s> with a %d", userID, roll, currentWinners[0], currWinnerRoll))
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d, <@%s> all winning with a %d", userID, roll, strings.Join(currentWinners, "><@"), currWinnerRoll))
	}
}

func ItsLotteryTime(s *discordgo.Session) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load PatronDB: %v\n", err)
		s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	winners, roll, err := patronDB.GetLotteryWinner()
	if err != nil {
		fmt.Printf("failed to GetLotteryWinner: %v\n", err)
		s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	var winnerFunds int
	share := 720 / len(winners)
	for _, winner := range winners {
		winnerFunds, err = patronDB.AddFunds(winner, share)
		if err != nil {
			fmt.Printf("failed to AddFunds: %v\n", err)
			s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
			return
		}
	}

	err = patronDB.ClearLottery()
	if err != nil {
		fmt.Printf("failed to ClearLottery: %v\n", err)
		s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	switch len(winners) {
	case 0:
		s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("No one rolled, no one wins ¯\\_(ツ)_/¯"))
	case 1:
		s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("<@%s> rolled a %d you win! You now have %d", winners[0], roll, winnerFunds))
	default:
		s.ChannelMessageSend("1020895617947537441", fmt.Sprintf("<@%s> all won with a %d you win! You each get %d", strings.Join(winners, "><@"), roll, share))
	}
}
