package dice

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gobo/pkg/bank"
)

var (
	session   *discordgo.Session
	msgCreate *discordgo.MessageCreate
)

const HelpPlay = "USAGE: !bet dice <amount> over | under"

func Play(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	session = s
	msgCreate = m

	bankDB, err := bank.LoadBankDB()
	if err != nil {
		fmt.Printf("failed to load bankDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	if len(params) == 3 {
		if strings.ToLower(params[0]) == "dice" {
			amt, err := strconv.Atoi(params[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Not a valid amount to bet")
				return
			}

			if amt < 0 {
				s.ChannelMessageSend(m.ChannelID, "You may only bet positive monies")
				return
			}

			userID := m.Author.ID
			funds, err := bankDB.CheckFunds(userID)
			if err != nil {

			}

			if funds < amt {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, <@%s> only has %d monies", userID, funds))
				return
			}

			one, _ := rand.Int(rand.Reader, big.NewInt(6))
			two, _ := rand.Int(rand.Reader, big.NewInt(6))
			n := one.Int64() + two.Int64() + 2

			overOrUnder := strings.ToLower(params[2])
			if overOrUnder != "over" && overOrUnder != "under" {
				printHelp()
				return
			}

			if overOrUnder == "over" && n > 7 || overOrUnder == "under" && n < 7 {
				currentFunds, err := bankDB.AddFunds(userID, amt)
				if err != nil {

				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you win! You now have %d", userID, n, currentFunds))
			} else {
				currentFunds, err := bankDB.TakeFunds(userID, amt)
				if err != nil {

				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you lose ¯\\_(ツ)_/¯ you still have %d", userID, n, currentFunds))
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Dice is the only game currently")
		}
	}
}
