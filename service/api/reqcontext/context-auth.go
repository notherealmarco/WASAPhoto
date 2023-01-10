package reqcontext

import "github.com/notherealmarco/WASAPhoto/service/database"

type AuthStatus int

const (
	AUTHORIZED     = 0
	UNAUTHORIZED   = 1
	FORBIDDEN      = 2
	USER_NOT_FOUND = 3
)

// Authorization is the interface for an authorization provider
type Authorization interface {
	// Returns the type of the authorization provider
	GetType() string

	// Returns the ID of the currently logged in user
	GetUserID() string

	// Checks if the token is valid
	Authorized(db database.AppDatabase) (AuthStatus, error)

	// Checks if the given user and the currently logged in user are the same user
	UserAuthorized(db database.AppDatabase, uid string) (AuthStatus, error)
}
