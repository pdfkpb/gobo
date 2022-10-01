package userid

import (
	"errors"
	"fmt"
	"regexp"
)

var userIDReg = regexp.MustCompile("<@[0-9]{18}>")

var (
	ErrNotAUserID = errors.New("could not parse as userid")
)

type UserID string

func GetUserID(userID string) (UserID, error) {
	if userIDReg.Match([]byte(userID)) {
		id := userID[2:20]
		return (UserID)(id), nil
	}
	return "", ErrNotAUserID
}

func (u UserID) Mention() string {
	return fmt.Sprintf("<@%s>", u)
}
