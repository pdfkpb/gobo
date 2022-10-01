package dicechallenge

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/commands"
	"github.com/pdfkpb/gobo/pkg/patron"
)

var _ commands.Exec = DiceChallenge

const (
	helpChallenge = "!dc @SomePlayer <amount>"
	helpAccept    = "!dc @SomePlayer"
	helpCancel    = "!dc cancel"
)

var (
	HelpPlay = fmt.Sprintf("Dice Challenge:\n\tChallenge: %s\n\tAccept: %s\n\tCancel: %s", helpChallenge, helpAccept, helpCancel)
)

func DiceChallenge(params []commands.Parameter, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	switch len(params) {
	case 1:
		if params[0] == "cancel" {
			cancel(patronDB, params, s, m)
		} else {
			accept(patronDB, params, s, m)
		}
	case 2:
		challenge(patronDB, params, s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid number of params")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}
}
