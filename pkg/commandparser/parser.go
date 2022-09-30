package commandparser

import "strings"

type ParsedCommand struct {
	Command Command
	Params  []Parameter
}

func ParseCommand(raw string) *ParsedCommand {
	pCmd := &ParsedCommand{}
	for _, value := range strings.Split(raw, " ") {
		pCmd.Params = append(pCmd.Params, *ParseParameter(value))
	}
	return pCmd
}
