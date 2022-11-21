package authorization

import (
	"errors"
	"net/http"

	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/sirupsen/logrus"
)

func BuildAuth(header string) (reqcontext.Authorization, error) {
	auth, err := BuildBearer(header)
	if err != nil {
		if err.Error() == "invalid authorization header" {
			return nil, errors.New("method not supported") // todo: better error description
		}
		return nil, err
	}
	return auth, nil
}

func SendAuthorizationError(f func(db database.AppDatabase, uid string) (reqcontext.AuthStatus, error), uid string, db database.AppDatabase, w http.ResponseWriter, l logrus.FieldLogger, notFoundStatus int) bool {
	auth, err := f(db, uid)
	if err != nil {
		helpers.SendInternalError(err, "Authorization error", w, l)
		return false
	}
	if auth == reqcontext.UNAUTHORIZED {
		helpers.SendStatus(http.StatusUnauthorized, w, "Unauthorized", l)
		return false
	}
	if auth == reqcontext.FORBIDDEN {
		helpers.SendStatus(http.StatusForbidden, w, "Forbidden", l)
		return false
	}
	// requested user is not found -> 404 as the resource is not found
	if auth == reqcontext.USER_NOT_FOUND {
		helpers.SendStatus(notFoundStatus, w, "Resource not found", l)
		return false
	}
	return true
}

func SendErrorIfNotLoggedIn(f func(db database.AppDatabase) (reqcontext.AuthStatus, error), db database.AppDatabase, w http.ResponseWriter, l logrus.FieldLogger) bool {

	auth, err := f(db)

	if err != nil {
		helpers.SendInternalError(err, "Authorization error", w, l)
		return false
	}

	if auth == reqcontext.UNAUTHORIZED {
		helpers.SendStatus(http.StatusUnauthorized, w, "Unauthorized", l)
		return false
	}

	return true
}
