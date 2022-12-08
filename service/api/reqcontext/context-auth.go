package reqcontext

import "github.com/notherealmarco/WASAPhoto/service/database"

type AuthStatus int

const (
	AUTHORIZED     = 0
	UNAUTHORIZED   = 1
	FORBIDDEN      = 2
	USER_NOT_FOUND = 3
)

type Authorization interface {
	GetType() string
	GetUserID() string
	Authorized(db database.AppDatabase) (AuthStatus, error)
	UserAuthorized(db database.AppDatabase, uid string) (AuthStatus, error)
}
