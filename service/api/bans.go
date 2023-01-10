package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

func (rt *_router) GetUserBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get user id
	uid := ps.ByName("user_id")

	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		// A user should not be able to see other users' bans
		return
	}

	// Get limits, or use default values
	start_index, limit, err := helpers.GetLimits(r.URL.Query())

	if err != nil {
		// Send error if the limits are specified but invalid
		helpers.SendBadRequest(w, "Invalid start_index or limit value", rt.baseLogger)
		return
	}

	// Get bans
	// We don't need to check if the user exists, because the authorization middleware already did that
	bans, err := rt.db.GetUserBans(uid, start_index, limit)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetUserBans", w, rt.baseLogger)
		return
	}

	// Return ban list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // is it needed?

	err = json.NewEncoder(w).Encode(bans) // write the response

	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}

func (rt *_router) PutBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	banned := ps.ByName("ban_uid")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	if uid == banned {
		helpers.SendBadRequest(w, "You cannot ban yourself", rt.baseLogger)
		return
	}

	// Execute the query
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
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	// Execute the query
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
