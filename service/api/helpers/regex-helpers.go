package helpers

import (
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

// Given a string, a regex and an error description, if the string doesn't match the regex, it sends a bad request error to the client and return false
// Otherwise it returns true without sending anything to the client
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

// Validates a username (must be between 3 and 16 characters long and can only contain letters, numbers and underscores)
func MatchUsernameOrBadRequest(username string, w http.ResponseWriter, l logrus.FieldLogger) bool {
	return MatchRegexOrBadRequest(username,
		`^[a-zA-Z0-9_]{3,16}$`, "Username must be between 3 and 16 characters long and can only contain letters, numbers and underscores",
		w,
		l)
}

// Validates a comment (must be between 1 and 255 characters long)
func MatchCommentOrBadRequest(comment string, w http.ResponseWriter, l logrus.FieldLogger) bool {
	return MatchRegexOrBadRequest(comment,
		`^(.){1,255}$`, "Comment must be between 1 and 255 characters long",
		w,
		l)
}
