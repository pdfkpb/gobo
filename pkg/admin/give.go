package admin

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/bwmarrin/discordgo"
)

func Give(bank map[string]int, params []string, s *discordgo.Session, m *discordgo.MessageCreate) {
    if (m.Author.ID == "303750733700923392" || m.Author.ID == "384902507383619594") && len(params) == 2 {
        if match, err := regexp.Match("<@[0-9]{18}>", []byte(params[0])); match && err == nil {
            uid := params[0][2:20]
            giveTo, err := s.User(uid)
            if err != nil || giveTo == nil {
                s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hmm not sure what fucked up there 'twas probably %s", uid))
                return
            }

            amt, err := strconv.Atoi(params[1])
            if err != nil {
                s.ChannelMessageSend(m.ChannelID, "Somehow failed to parse the monies")
                return
            }

            if amt < 0 {
                s.ChannelMessageSend(m.ChannelID, "Cannot give negative funds, must take positive")
                return
            }

            bank[uid] += amt
            s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %d monies", params[0], bank[uid]))
        } else {
            s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Fuck your regex and %s's eyebrows", params[0]))
        }
    } else {
        s.ChannelMessageSend(m.ChannelID, "Only I Gobo... or Dylan/Kevin may give funds")
    }
}
