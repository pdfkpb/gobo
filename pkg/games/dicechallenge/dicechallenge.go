package dicechallenge

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/games"
	"github.com/pdfkpb/gobo/pkg/patron"
)

var _ games.Play = Play

const (
	helpChallenge = "!dc @SomePlayer <amount>"
	helpAccept    = "!dc @SomePlayer"
	helpCancel    = "!dc cancel"
)

var (
	HelpPlay = fmt.Sprintf("Dice Challenge:\n\tChallenge: %s\n\tAccept: %s\n\tCancel: %s", helpChallenge, helpAccept, helpCancel)
)

func Play(params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	patronDB, err := patron.LoadPatronDB()
	if err != nil {
		fmt.Printf("failed to load patronDB: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Some backend error occured <@384902507383619594> fix it"))
		return
	}

	switch len(params) {
	case 1:
		accept(patronDB, params, s, m)
	case 2:
		if strings.Compare(params[0], "cancel") == 0 {
			cancel(patronDB, params, s, m)
		} else {
			challenge(patronDB, params, s, m)
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid number of params")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("  Usage:\n```%s```", HelpPlay))
		return
	}
}
