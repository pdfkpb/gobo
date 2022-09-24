package games

import "github.com/bwmarrin/discordgo"

// HouseCut - funds * HouseCut to take it
const HouseCut float32 = 0.03 // 3%

// TakeHouseCut takes a small cut of the total
func TakeHouseCut(total int) int {
	return int(float32(total) * (1.0 - HouseCut))
}

// Play all games should define this function
type Play func(params []string, s *discordgo.Session, m *discordgo.MessageCreate)
