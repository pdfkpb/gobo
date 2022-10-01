package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfkpb/gobo/pkg/admin"
	"github.com/pdfkpb/gobo/pkg/games/dice"
	"github.com/pdfkpb/gobo/pkg/games/dicechallenge"
	"github.com/pdfkpb/gobo/pkg/games/lottery"
)

const (
	errNotMe = "message not for bot"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	defer dg.Close()

	lotteryTicker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-lotteryTicker.C:
			lottery.ItsLotteryTime(dg)
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd, params, err := commandParse(m.Content)
	if err != nil {
		if err.Error() == errNotMe {
			return
		}
		s.ChannelMessageSend(m.ChannelID, "The fuck are you on about?")
		return
	}

	switch cmd {
	case "give":
		admin.Give(params, s, m)
	case "take":
		admin.Take(params, s, m)
	case "check":
		admin.Check(params, s, m)
	case "register":
		admin.RegisterUser(params, s, m)
	case "br":
		admin.BulkRegister(params, s, m)
	case "dice":
		dice.Play(params, s, m)
	case "dc":
		dicechallenge.Play(params, s, m)
	case "roll":
		lottery.Play(params, s, m)
	case "help":
		gamesHelp := fmt.Sprintf("Games:\n```%s\n%s\n%s```", dice.HelpPlay, lottery.HelpPlay, dicechallenge.HelpPlay)
		s.ChannelMessageSend(m.ChannelID, gamesHelp)

		adminHelp := fmt.Sprintf("Admin: params in brackets require admin priveleges\n```\n%s\n%s```", admin.HelpCheck, admin.HelpRegister)
		s.ChannelMessageSend(m.ChannelID, adminHelp)
	default:
		s.ChannelMessageSend(m.ChannelID, "Gobo here, type `!help` to see a list of commands")
	}
}