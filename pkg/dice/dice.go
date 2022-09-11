package dice

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	session   *discordgo.Session
	msgCreate *discordgo.MessageCreate
)

func printHelp() {
	session.ChannelMessageSend(
		msgCreate.ChannelID,
		fmt.Sprintf("USAGE: !bet dice <amount> over | under"),
	)
}

func Go(bank map[string]int, params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	session = s
	msgCreate = m

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

			user := m.Author.ID
			if bank[user] < amt {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Insufficent funds, <@%s> only has %d monies", user, bank[user]))
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
				bank[user] += amt
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you win! You now have %d", user, n, bank[user]))
			} else {
				bank[user] -= amt
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> rolled a %d you lose ¯\\_(ツ)_/¯ you still have %d", user, n, bank[user]))
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Dice is the only game currently")
		}
	}
}
