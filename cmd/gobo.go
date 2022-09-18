package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gobo/pkg/admin"
	"github.com/gobo/pkg/dice"
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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

var bank = map[string]int{}

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
	case "dice":
		dice.Play(params, s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Games: "))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```  %s```", dice.HelpPlay))

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Admin: "))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```  %s```", admin.HelpGive))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```  %s```", admin.HelpTake))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```  %s```", admin.HelpCheck))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```  %s```", admin.HelpRegister))
	}
}

func commandParse(cmd string) (string, []string, error) {
	fmt.Println(cmd)

	if !strings.HasPrefix(cmd, "!") && !strings.HasPrefix(cmd, "\\!") {
		return "", []string{}, errors.New(errNotMe)
	}

	if strings.HasPrefix(cmd, "\\!") {
		cmd = cmd[1:]
	}

	pCmd := strings.Split(cmd, " ")
	trueCmd := pCmd[0][1:]
	params := []string{}
	if len(pCmd) > 1 {
		params = pCmd[1:]
	}

	return trueCmd, params, nil
}
