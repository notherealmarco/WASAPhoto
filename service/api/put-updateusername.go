package api

import (
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

func (rt *_router) UpdateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}
	var req structures.UserDetails
	if !helpers.DecodeJsonOrBadRequest(r.Body, w, &req, rt.baseLogger) {
		return
	}

	stat, err := regexp.Match(`^[a-zA-Z0-9_]{3,16}$`, []byte(req.Name))

	if err != nil {
		helpers.SendInternalError(err, "Error while matching username", w, rt.baseLogger)
		return
	}

	if !stat { // todo: sta regex non me piace
		helpers.SendBadRequest(w, "Username must be between 3 and 16 characters long and can only contain letters, numbers and underscores", rt.baseLogger)
		return
	}

	status, err := rt.db.UpdateUsername(uid, req.Name)

	if status == database.ERR_EXISTS {
		helpers.SendStatus(http.StatusConflict, w, "Username already exists", rt.baseLogger)
		return
	}

	if err != nil {
		helpers.SendInternalError(err, "Database error: UpdateUsername", w, rt.baseLogger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
