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

func (rt *_router) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get user id
	uid := ps.ByName("user_id")

	if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) ||
		!helpers.SendNotFoundIfBanned(rt.db, ctx.Auth.GetUserID(), uid, w, rt.baseLogger) {
		return
	}

	// Get user profile
	status, profile, err := rt.db.GetUserProfile(uid)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetUserProfile", w, rt.baseLogger)
		return
	}

	if status == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User not found", rt.baseLogger)
		return
	}

	// Return user profile
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profile)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}
