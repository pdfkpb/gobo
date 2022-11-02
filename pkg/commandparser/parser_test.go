package commandparser_test

import (
	"testing"

	"github.com/pdfkpb/gobo/pkg/commandparser"
)

func FuzzParseCommand(f *testing.F) {
	corpus := map[string]commandparser.ParsedCommand{
		"!test abc 123 <@546745767457876>": commandparser.ParsedCommand{},
	}
}
