package admin

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Check(bank map[string]int, params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has %d monies", m.Author.ID, bank[m.Author.ID]))
}
