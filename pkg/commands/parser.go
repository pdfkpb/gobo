package commands

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

	cmd := splitCmd[0]
	if !strings.HasPrefix(cmd, "!") && !strings.HasPrefix(cmd, "\\!") {
		return nil, ErrNotACommand
	}

	if strings.HasPrefix(cmd, "\\!") {
		cmd = cmd[1:]
	}

	pCmd := &ParsedCommand{
		Command: commandFromString(cmd[1:]),
	}

	for _, value := range splitCmd[1:] {
		pCmd.Params = append(pCmd.Params, *ParseParameter(value))
	}
	return pCmd, nil
}
