package authorization

import (
	"errors"
	"net/http"

	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/sirupsen/logrus"
)

// BuildAuth returns an Authorization implementation for the currently logged in user
func BuildAuth(header string) (reqcontext.Authorization, error) {
	auth, err := BuildBearer(header)
	if err != nil {
		if err.Error() == "invalid authorization header" {
			return nil, errors.New("authentication method not supported")
		}
		return nil, err
	}
	return auth, nil
}

// Given a user authorization function, if the function returns some error, it sends the error to the client and return false
// Otherwise it returns true without sending anything to the client
func SendAuthorizationError(f func(db database.AppDatabase, uid string) (reqcontext.AuthStatus, error), uid string, db database.AppDatabase, w http.ResponseWriter, l logrus.FieldLogger, notFoundStatus int) bool {
	auth, err := f(db, uid)
	if err != nil {
		helpers.SendInternalError(err, "Authorization error", w, l)
		return false
	}
	if auth == reqcontext.UNAUTHORIZED {
		// The token is not valid
		helpers.SendStatus(http.StatusUnauthorized, w, "Unauthorized", l)
		return false
	}
	if auth == reqcontext.FORBIDDEN {
		// The user is not authorized for this action
		helpers.SendStatus(http.StatusForbidden, w, "Forbidden", l)
		return false
	}
	if auth == reqcontext.USER_NOT_FOUND {
		// Attempting to perform an action on a non-existent user
		helpers.SendStatus(notFoundStatus, w, "User not found", l)
		return false
	}
	return true
}

// Given a function that validates a token, if the function returns some error, it sends the error to the client and return false
// Otherwise it returns true without sending anything to the client
func SendErrorIfNotLoggedIn(f func(db database.AppDatabase) (reqcontext.AuthStatus, error), db database.AppDatabase, w http.ResponseWriter, l logrus.FieldLogger) bool {

	auth, err := f(db)

	if err != nil {
		helpers.SendInternalError(err, "Authorization error", w, l)
		return false
	}

	if auth == reqcontext.UNAUTHORIZED {
		// The token is not valid
		helpers.SendStatus(http.StatusUnauthorized, w, "Unauthorized", l)
		return false
	}

	return true
}
