// This identity provider represents non logged-in users.

package authorization

import (
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

// AnonymousAuth is the authentication provider for non logged-in users
type AnonymousAuth struct {
}

// Returns a newly created AnonymousAuth instance
func BuildAnonymous() *AnonymousAuth {
	return &AnonymousAuth{}
}

func (u *AnonymousAuth) GetType() string {
	return "Anonymous"
}

// Returns UNAUTHORIZED, as anonymous users are logged in
func (u *AnonymousAuth) Authorized(db database.AppDatabase) (reqcontext.AuthStatus, error) {
	return reqcontext.UNAUTHORIZED, nil
}

// Returns UNAUTHORIZED, as anonymous users are not logged in
func (u *AnonymousAuth) UserAuthorized(db database.AppDatabase, uid string) (reqcontext.AuthStatus, error) {
	return reqcontext.UNAUTHORIZED, nil
}

// Returns an empty string, as anonymous users have no user ID
func (u *AnonymousAuth) GetUserID() string {
	return ""
}
