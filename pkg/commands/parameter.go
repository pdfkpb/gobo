package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pdfkpb/gobo/pkg/userid"
)

type ParameterType int64

const (
	ParamTypeUnknown ParameterType = iota
	ParamTypeString
	ParamTypeInteger
	ParamTypeFloat
	ParamTypeUserID
)

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
	userID, err := userid.GetUserID(lower)
	if err == nil {
		param.pType = ParamTypeUserID
		param.value = userID
		return param
	}

	// Generate Integer Param
	fmt.Println(trimmed)
	integer, err := strconv.Atoi(trimmed)
	if err == nil {
		param.pType = ParamTypeInteger
		param.value = integer
		return param
	}

	// Generate Float Param
	float, err := strconv.ParseFloat(trimmed, 32)
	if err == nil {
		param.pType = ParamTypeFloat
		param.value = (float32)(float)
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

func (p *Parameter) String() string {
	if str, ok := p.value.(string); ok {
		return str
	}

	return ""
}

func (p *Parameter) Integer() int {
	if integer, ok := p.value.(int); ok {
		return integer
	}

	return 0
}

func (p *Parameter) Float() float32 {
	if f32, ok := p.value.(float32); ok {
		return f32
	}

	return 0.0
}

func (p *Parameter) UserID() userid.UserID {
	if userID, ok := p.value.(userid.UserID); ok {
		return userID
	}

	return ""
}
