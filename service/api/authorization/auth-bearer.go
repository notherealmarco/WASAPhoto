package authorization

import (
	"errors"
	"strings"

	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

// BearerAuth is the authentication provider that authorizes users by Bearer tokens
// In this case, a token is the unique identifier for a user.
type BearerAuth struct {
	token string
}

func (b *BearerAuth) GetType() string {
	return "Bearer"
}

// Given the content of the Authorization header, returns a BearerAuth instance for the user
// Returns an error if the header is not valid
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

// Returns the user ID of the user that is currently logged in
func (b *BearerAuth) GetUserID() string {
	return b.token
}

// Checks if the token is valid
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

// Checks if the given user and the currently logged in user are the same user
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
		// If the user is the same as the one in the token, check if the user does actually exist in the database
		auth, err := b.Authorized(db)

		if err != nil {
			return -1, err
		}

		return auth, nil
	}
	// If the user is not the same as the one in the token, return FORBIDDEN
	return reqcontext.FORBIDDEN, nil
}
