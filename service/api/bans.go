package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

func (rt *_router) PutBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	banned := ps.ByName("ban_uid")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, http.StatusNotFound) {
		return
	}

	if uid == banned {
		helpers.SendBadRequest(w, "You cannot ban yourself", rt.baseLogger)
		return
	}

	status, err := rt.db.BanUser(uid, banned)

	if err != nil {
		helpers.SendInternalError(err, "Database error: BanUser", w, rt.baseLogger)
		return
	}

	if status == database.ERR_NOT_FOUND {
		helpers.SendBadRequest(w, "You are trying to ban a non-existent user", rt.baseLogger)
		return
	}

	if status == database.ERR_EXISTS {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	helpers.SendStatus(http.StatusCreated, w, "Success", rt.baseLogger)
}

func (rt *_router) DeleteBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	uid := ps.ByName("user_id")
	banned := ps.ByName("ban_uid")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, http.StatusNotFound) {
		return
	}

	status, err := rt.db.UnbanUser(uid, banned)

	if err != nil {
		helpers.SendInternalError(err, "Database error: UnbanUser", w, rt.baseLogger)
		return
	}

	if status == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User not banned", rt.baseLogger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
