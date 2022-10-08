package lottery

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
	"github.com/pdfkpb/gobo/pkg/userid"
)

var _ commands.Exec = LotteryRoll

const HelpPlay = "Lottery:\n\t!roll"
const payout = 720

func LotteryRoll(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
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

	userID := userid.UserID(m.Author.ID)
	err = patronDB.SetLotteryRoll(string(userID), realRoll)
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

	currentWinners, currWinnerRoll, err := patronDB.GetLotteryWinners()
	if err != nil {
		switch err {
		case patron.ErrNoRoll:
			fmt.Println(err)
		default:
			fmt.Printf("lottery failed to get lottery winners %v\n", err)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Some backend error occured <@384902507383619594> fix it"))
			return
		}
	}

	var currentWinnerIDs []string
	for _, cWinner := range currentWinners {
		tID := userid.UserID(cWinner)
		currentWinnerIDs = append(currentWinnerIDs, tID.Mention())
	}

	switch len(currentWinners) {
	case 1:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s rolled a %d, Current Winner is %s with a %d", userID.Mention(), realRoll, currentWinnerIDs[0], currWinnerRoll))
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s rolled a %d, %s all winning with a %d", userID.Mention(), realRoll, strings.Join(currentWinnerIDs, " "), currWinnerRoll))
	}
}

func ItsLotteryTime(s *discordgo.Session) {
	const channelID = "1023752170442608720"

	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load PatronDB: %v\n", err)
		s.ChannelMessageSend(channelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	winners, roll, err := patronDB.GetLotteryWinners()
	if err != nil {
		switch err {
		case patron.ErrNoRoll:
			fmt.Println(err)
		default:
			fmt.Printf("failed to GetLotteryWinner: %v\n", err)
			s.ChannelMessageSend(channelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
			return
		}
	}

	var winnerFunds int
	var share int
	if len(winners) > 0 {
		share = 720 / len(winners)
	}

	for _, winner := range winners {
		winnerFunds, err = patronDB.AddFunds(winner, share)
		if err != nil {
			fmt.Printf("failed to AddFunds: %v\n", err)
			s.ChannelMessageSend(channelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
			return
		}
	}

	err = patronDB.ClearLottery()
	if err != nil {
		fmt.Printf("failed to ClearLottery: %v\n", err)
		s.ChannelMessageSend(channelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	switch len(winners) {
	case 0:
		break
	case 1:
		s.ChannelMessageSend(channelID, fmt.Sprintf("<@%s> rolled a %d you win! You now have %d", winners[0], roll, winnerFunds))
	default:
		s.ChannelMessageSend(channelID, fmt.Sprintf("<@%s> all won with a %d you win! You each get %d", strings.Join(winners, "><@"), roll, share))
	}
}
