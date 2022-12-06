// This identity provider represents non logged-in users.

package authorization

import (
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

type AnonymousAuth struct {
}

func BuildAnonymous() *AnonymousAuth {
	return &AnonymousAuth{}
}

func (u *AnonymousAuth) GetType() string {
	return "Anonymous"
}

func (u *AnonymousAuth) Authorized(db database.AppDatabase) (reqcontext.AuthStatus, error) {
	return reqcontext.UNAUTHORIZED, nil
}

func (u *AnonymousAuth) UserAuthorized(db database.AppDatabase, uid string) (reqcontext.AuthStatus, error) {
	return reqcontext.UNAUTHORIZED, nil
}

func (u *AnonymousAuth) GetUserID() string {
	return ""
}
