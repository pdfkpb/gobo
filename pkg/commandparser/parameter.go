package commandparser

import (
	"regexp"
	"strconv"
	"strings"
)

type ParameterType int64

const (
	ParamTypeUnknown ParameterType = iota
	ParamTypeString
	ParamTypeInteger
	ParamTypeFloat
	ParamTypeUserID
)

var userIDReg = regexp.MustCompile("<@[0-9]{18}>")

type Parameter struct {
	Raw string

	pType ParameterType
	value any
}

func ParseParameter(raw string) *Parameter {
	trimmed := strings.TrimSpace(raw)
	lower := strings.ToLower(trimmed)

	param := &Parameter{
		Raw: raw,
	}

	// Generate UserID Param
	isUserID := userIDReg.Match([]byte(lower))
	if isUserID {
		param.pType = ParamTypeUserID
		param.value = lower
		return param
	}

	// Generate Float Param
	float, err := strconv.ParseFloat(trimmed, 32)
	if err == nil {
		param.pType = ParamTypeFloat
		param.value = float
		return param
	}

	// Generate Integer Param
	integer, err := strconv.Atoi(trimmed)
	if err == nil {
		param.pType = ParamTypeInteger
		param.value = integer
		return param
	}

	// Return a String type Param
	param.pType = ParamTypeString
	param.value = lower
	return param
}

func (p *Parameter) Type() ParameterType {
	return p.pType
}
