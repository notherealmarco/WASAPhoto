package authorization

import (
	"errors"
	"strings"

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

func (b *BearerAuth) GetToken() string {
	return b.token
}

func (b *BearerAuth) Authorized(db database.AppDatabase) (bool, error) {
	// this is the way we manage authorization, the bearer token is the user id
	state, err := db.UserExists(b.token)

	if err != nil {
		return false, err
	}
	return state, nil
}

func (b *BearerAuth) UserAuthorized(db database.AppDatabase, uid string) (bool, error) {
	if b.token == uid {
		return b.Authorized(db)
	}
	return false, nil
}
