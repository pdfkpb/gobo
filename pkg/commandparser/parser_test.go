package commandparser_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pdfkpb/gobo/pkg/commandparser"
)

func TestParserSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

var _ = Describe("command_parser", func() {
	Describe("ParseCommand", func() {
		Context("where valid command with one of each types is passed", func() {
			It("should properly return a struct with valid Parameters", func() {
				pCmd := commandparser.ParseCommand("!test abc 123 0.04 <@384902507383619594>")

				Expect(pCmd.Command).To(Equal("test"))
				Expect(len(pCmd.Params)).To(Equal(4))

				// Test the string case
				Expect(pCmd.Params[0].Type()).To(Equal(commandparser.ParamTypeString))
				Expect(pCmd.Params[0].String()).To(Equal("abc"))

				// Test the integer case
				Expect(pCmd.Params[0].Type()).To(Equal(commandparser.ParamTypeInteger))
				Expect(pCmd.Params[0].Integer()).To(Equal(123))

				// Test the float case
				Expect(pCmd.Params[0].Type()).To(Equal(commandparser.ParamTypeFloat))
				Expect(pCmd.Params[0].Float()).To(Equal(0.04))

				// Test the userid case
				Expect(pCmd.Params[0].Type()).To(Equal(commandparser.ParamTypeUserID))
				Expect(pCmd.Params[0].UserID().Mention()).To(Equal("<@384902507383619594>"))
			})
		})
	})
})
