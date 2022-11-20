package authorization

import (
	"errors"
	"strings"

	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

type BearerAuth struct {
	token string
}

func (b *BearerAuth) GetType() string {
	return "Bearer"
}

func BuildBearer(header string) (*BearerAuth, error) {
	if header == "" {
		return nil, errors.New("missing authorization header")
	}
	if header == "Bearer" {
		return nil, errors.New("missing token")
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, errors.New("invalid authorization header")
	}
	return &BearerAuth{token: header[7:]}, nil
}

func (b *BearerAuth) GetUserID() string {
	return b.token
}

func (b *BearerAuth) Authorized(db database.AppDatabase) (reqcontext.AuthStatus, error) {
	// this is the way we manage authorization, the bearer token is the user id
	state, err := db.UserExists(b.token)

	if err != nil {
		return reqcontext.UNAUTHORIZED, err
	}

	if state {
		return reqcontext.AUTHORIZED, nil
	}
	return reqcontext.UNAUTHORIZED, nil
}

func (b *BearerAuth) UserAuthorized(db database.AppDatabase, uid string) (reqcontext.AuthStatus, error) {

	// If uid is not a valid user, return USER_NOT_FOUND
	user_exists, err := db.UserExists(uid)

	if err != nil {
		return reqcontext.UNAUTHORIZED, err
	}
	if !user_exists {
		return reqcontext.USER_NOT_FOUND, nil
	}

	if b.token == uid {
		auth, err := b.Authorized(db)

		if err != nil {
			return -1, err
		}

		return auth, nil
	}
	return reqcontext.FORBIDDEN, nil
}
