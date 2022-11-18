package authorization

import (
	"errors"
	"net/http"

	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
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

func SendAuthorizationError(f func(db database.AppDatabase, uid string) (reqcontext.AuthStatus, error), uid string, db database.AppDatabase, w http.ResponseWriter) bool {
	auth, err := f(db, uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// todo: log error and write it to the response
		return false
	}
	if auth == reqcontext.UNAUTHORIZED {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if auth == reqcontext.FORBIDDEN {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}
