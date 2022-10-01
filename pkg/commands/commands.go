package commands

import "github.com/bwmarrin/discordgo"

type Exec func([]Parameter, *discordgo.Session, *discordgo.MessageCreate)

type Command int64

const (
	Unknown Command = iota
	Help
	//Admin
	BulkRegister
	Register
	Check
	Give
	Take
	// Games
	Dice
	DiceChallenge
	Lottery
)

var commandMap = map[Command]string{
	Unknown:       "unknown",
	BulkRegister:  "bulkRegister",
	Register:      "register",
	Check:         "check",
	Give:          "give",
	Take:          "take",
	Dice:          "dice",
	DiceChallenge: "diceChallenge",
	Lottery:       "lottery",
}

func commandFromString(cmdStr string) Command {
	for k, v := range commandMap {
		if v == cmdStr {
			return k
		}
	}
	return Unknown
}

func (c *Command) String() string {
	return commandMap[*c]
}
