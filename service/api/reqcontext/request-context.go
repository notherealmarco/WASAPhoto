/*
Package reqcontext contains the request context. Each request will have its own instance of RequestContext filled by the
middleware code in the api-context-wrapper.go (parent package).

Each value here should be assumed valid only per request only, with some exceptions like the logger.
*/
package reqcontext

import (
	"github.com/gofrs/uuid"
	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/sirupsen/logrus"
)

// RequestContext is the context of the request, for request-dependent parameters
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID uuid.UUID

	// Logger is a custom field logger for the request
	Logger logrus.FieldLogger

	Auth Authorization
}

type Authorization interface {
	GetType() string
	Authorized(db database.AppDatabase) (bool, error)
	UserAuthorized(db database.AppDatabase, uid string) (bool, error)
}
