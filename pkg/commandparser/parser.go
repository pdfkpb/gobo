package commandparser

import (
	"errors"
	"strings"
)

var (
	ErrNotACommand = errors.New("failed to parse first parameter as a command")
)

type ParsedCommand struct {
	Command Command
	Params  []Parameter
}

func ParseCommand(raw string) (*ParsedCommand, error) {
	splitCmd := strings.Split(raw, " ")

	var cmd Command
	if !strings.HasPrefix(cmd, "!") && !strings.HasPrefix(cmd, "\\!") {
		return nil, ErrNotACommand
	}

	if strings.HasPrefix(cmd, "\\!") {
		cmd = cmd[1:]
	}

	pCmd := &ParsedCommand{}
	for _, value := range splitCmd[1:] {
		pCmd.Params = append(pCmd.Params, *ParseParameter(value))
	}
	return pCmd, nil
}
