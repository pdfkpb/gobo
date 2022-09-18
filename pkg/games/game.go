package games

import "github.com/bwmarrin/discordgo"

type Game interface {
	Play(params []string, s *discordgo.Session, m *discordgo.MessageCreate)
}
