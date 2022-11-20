package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database"
)

func (rt *_router) GetLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// get the user id from the url
	uid := ps.ByName("user_id")
	photo_id, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)

	if err != nil {
		helpers.SendBadRequest(w, "Invalid photo id", rt.baseLogger)
		return
	}

	// send 404 if the user does not exist
	if !helpers.VerifyUserOrNotFound(rt.db, uid, w, rt.baseLogger) {
		return
	}

	// get the user's likes
	success, likes, err := rt.db.GetPhotoLikes(uid, photo_id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetLikes", w, rt.baseLogger)
		return
	}

	if success == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User or photo not found", rt.baseLogger)
		return
	}

	// send the response
	err = json.NewEncoder(w).Encode(likes)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding response", w, rt.baseLogger)
		return
	}
}

func (rt *_router) PutDeleteLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")

	photo_id, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		helpers.SendBadRequestError(err, "Bad photo_id", w, rt.baseLogger)
		return
	}

	liker_uid := ps.ByName("liker_uid")

	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, liker_uid, rt.db, w, rt.baseLogger, http.StatusBadRequest) {
		return
	}

	var success database.QueryResult

	// If the request is a PUT, then we like the photo
	if r.Method == "PUT" {
		success, err = rt.db.LikePhoto(uid, photo_id, liker_uid)
	} else { // Request is a DELETE, so we unlike the photo
		success, err = rt.db.UnlikePhoto(uid, photo_id, liker_uid)
	}

	if err != nil {
		helpers.SendInternalError(err, "Database error", w, rt.baseLogger)
		return
	}

	if success == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "Resource not found", rt.baseLogger)
		return
	}

	// User already liked the photo
	if success == database.ERR_EXISTS {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == "PUT" {
		// User liked the photo successfully
		helpers.SendStatus(http.StatusCreated, w, "Success", rt.baseLogger)
	} else {
		// User unliked the photo successfully
		w.WriteHeader(http.StatusNoContent)
	}
}
