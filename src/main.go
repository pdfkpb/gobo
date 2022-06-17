package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
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

	if m.Author.Username == "GeekMartyr" {
		s.ChannelMessageSend(m.ChannelID, "Nerd")
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

	//fmt.Println(cmd)
	//fmt.Println(params)

	switch cmd {
	case "give":
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
	case "take":
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
					s.ChannelMessageSend(m.ChannelID, "Cannot take negative funds, must give positive")
					return
				}

				bank[uid] -= amt
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has %d monies", params[0], bank[uid]))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Fuck your regex and %s's eyebrows", params[0]))
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Only I Gobo... or Dylan/Kevin may take funds")
		}
	case "check":
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has %d monies", m.Author.ID, bank[m.Author.ID]))
	case "bet":
		switch strings.ToLower(params[0]) {
		case "dice":
			dice.Go(bank, params, s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "Game not implemented")
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Go fuck yoself")
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
