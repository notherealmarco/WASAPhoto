package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

func (rt *_router) GetFollowersFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) {
		return
	}

	uid := ps.ByName("user_id")

	if !helpers.VerifyUserOrNotFound(rt.db, uid, w, rt.baseLogger) {
		return
	}

	var users *[]structures.UIDName
	var err error
	var status database.QueryResult

	// Get limits, or use default values
	start_index, limit, err := helpers.GetLimits(r.URL.Query())

	if err != nil {
		helpers.SendBadRequest(w, "Invalid start_index or limit", rt.baseLogger)
		return
	}

	// Check if client is asking for followers or following
	if strings.HasSuffix(r.URL.Path, "/followers") {
		// Get the followers from the database
		status, users, err = rt.db.GetUserFollowers(uid, ctx.Auth.GetUserID(), start_index, limit)
	} else {
		// Get the following users from the database
		status, users, err = rt.db.GetUserFollowing(uid, ctx.Auth.GetUserID(), start_index, limit)
	}

	// Send a 500 response if there was an error
	if err != nil {
		helpers.SendInternalError(err, "Database error: GetUserFollowers", w, rt.baseLogger)
		return
	}

	// Send a 404 response if the user was not found
	if status == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User not found", rt.baseLogger)
		return
	}

	// Send the users to the client
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(users)

	// Send a 500 response if there was an error
	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}

func (rt *_router) PutFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	followed := ps.ByName("follower_uid")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	status, err := rt.db.FollowUser(uid, followed)

	if err != nil {
		helpers.SendInternalError(err, "Database error: FollowUser", w, rt.baseLogger)
		return
	}

	if status == database.ERR_EXISTS {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if status == database.ERR_NOT_FOUND {
		helpers.SendBadRequest(w, "You are trying to follow a non-existent user", rt.baseLogger)
		return
	}

	helpers.SendStatus(http.StatusCreated, w, "Success", rt.baseLogger)
}

func (rt *_router) DeleteFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	followed := ps.ByName("follower_uid")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	status, err := rt.db.UnfollowUser(uid, followed)

	if err != nil {
		helpers.SendInternalError(err, "Database error: UnfollowUser", w, rt.baseLogger)
		return
	}

	if status == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User not found", rt.baseLogger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
