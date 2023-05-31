package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

func (rt *_router) UpdateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")

	// check if the user is changing his own username
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	// decode request body
	var req structures.UserDetails
	if !helpers.DecodeJsonOrBadRequest(r.Body, w, &req, rt.baseLogger) {
		return
	}

	// check if the username is valid, and if it's not, send a bad request error
	if !helpers.MatchUsernameOrBadRequest(req.Name, w, rt.baseLogger) {
		return
	}

	status, err := rt.db.UpdateUsername(uid, req.Name)

	// check if the username already exists
	if status == database.ERR_EXISTS {
		helpers.SendStatus(http.StatusConflict, w, "Username already exists", rt.baseLogger)
		return
	}

	// handle any other database error
	if err != nil {
		helpers.SendInternalError(err, "Database error: UpdateUsername", w, rt.baseLogger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
