package commandparser

type Command int64

const (
	Unknown Command = iota
	//Admin
	BulkRegister
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
	Check:         "check",
	Give:          "give",
	Take:          "take",
	Dice:          "dice",
	DiceChallenge: "diceChallenge",
	Lottery:       "lottery",
}

func (c *Command) String() string {
	return commandMap[*c]
}
