package helpers

import (
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

func MatchRegexOrBadRequest(str string, regex string, error_description string, w http.ResponseWriter, l logrus.FieldLogger) bool {

	stat, err := regexp.Match(regex, []byte(str))

	if err != nil {
		SendInternalError(err, "Error while matching username regex", w, l)
		return false
	}

	if !stat {
		// string didn't match the regex, so it's invalid, let's send a bad request error
		SendBadRequest(w, error_description, l)
		return false
	}
	// string matched the regex, so it's valid
	return true
}

func MatchUsernameOrBadRequest(username string, w http.ResponseWriter, l logrus.FieldLogger) bool {
	return MatchRegexOrBadRequest(username,
		`^[a-zA-Z0-9_]{3,16}$`, "Username must be between 3 and 16 characters long and can only contain letters, numbers and underscores",
		w,
		l)
}
