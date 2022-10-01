package commands_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pdfkpb/gobo/pkg/commands"
)

func TestParserSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

var _ = Describe("command_parser", func() {
	Describe("ParseCommand", func() {
		Context("where valid command with one of each types is passed", func() {
			It("should properly return a struct with valid Parameters", func() {
				pCmd, err := commands.ParseCommand("!give abc 123 0.04 <@384902507383619594>")

				Expect(err).To(BeNil())

				Expect(pCmd.Command).To(Equal(commands.Give))
				Expect(len(pCmd.Params)).To(Equal(4))

				// Test the string case
				Expect(pCmd.Params[0].Type()).To(Equal(commands.ParamTypeString))
				Expect(pCmd.Params[0].String()).To(Equal("abc"))

				// Test the integer case
				Expect(pCmd.Params[1].Type()).To(Equal(commands.ParamTypeInteger))
				Expect(pCmd.Params[1].Integer()).To(Equal(123))

				// Test the float case
				Expect(pCmd.Params[2].Type()).To(Equal(commands.ParamTypeFloat))
				Expect(pCmd.Params[2].Float()).Should(BeNumerically("~", 0.04, .001))

				// Test the userid case
				Expect(pCmd.Params[3].Type()).To(Equal(commands.ParamTypeUserID))
				Expect(pCmd.Params[3].UserID().Mention()).To(Equal("<@384902507383619594>"))
			})
		})
	})
})
