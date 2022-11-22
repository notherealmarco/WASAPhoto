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

type reqbody struct {
	UID     string `json:"user_id"`
	Comment string `json:"comment"`
}

func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// get the user's comments
	success, comments, err := rt.db.GetComments(uid, photo_id, ctx.Auth.GetUserID())

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetComments", w, rt.baseLogger)
		return
	}

	if success == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User or photo not found", rt.baseLogger)
		return
	}

	// send the response
	err = json.NewEncoder(w).Encode(comments)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding comments", w, rt.baseLogger)
		return
	}
}

func (rt *_router) PostComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")

	photo_id, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		helpers.SendBadRequestError(err, "Bad photo_id", w, rt.baseLogger)
		return
	}

	// get the comment from the request
	var request_body reqbody
	if !helpers.DecodeJsonOrBadRequest(r.Body, w, &request_body, rt.baseLogger) {
		return
	}

	// check if the user is authorized to post a comment
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, request_body.UID, rt.db, w, rt.baseLogger, http.StatusBadRequest) {
		// It returns 400 Bad Request if the user_id field in the request body is missing or an invalid user_id
		// It returns 401 if the user is not logged in
		// It returns 403 if the user is not authorized to post a comment as the requested user
		return
	}

	// add the comment to the database
	success, err := rt.db.PostComment(uid, photo_id, request_body.UID, request_body.Comment)

	if err != nil {
		helpers.SendInternalError(err, "Database error: PostComment", w, rt.baseLogger)
		return
	}

	// if user or photo does not exist, send 404
	if success == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "User or photo not found", rt.baseLogger)
		return
	}

	// send the response
	helpers.SendStatus(http.StatusCreated, w, "Comment created with success", rt.baseLogger)
}

func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")

	photo_id, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		helpers.SendBadRequestError(err, "Bad photo_id", w, rt.baseLogger)
		return
	}

	comment_id, err := strconv.ParseInt(ps.ByName("comment_id"), 10, 64)
	if err != nil {
		helpers.SendBadRequestError(err, "Bad comment_id", w, rt.baseLogger)
		return
	}

	// Check if the user is authorized to delete that comment
	// (only the user who posted the comment or the owner of the photo can delete it)
	status, comment_owner, err := rt.db.GetCommentOwner(uid, photo_id, comment_id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetCommentOwner", w, rt.baseLogger)
		return
	}

	if status == database.ERR_NOT_FOUND {
		helpers.SendNotFound(w, "Resource not found", rt.baseLogger)
		return
	}

	// The authorized user must be either comment_owner or uid
	owner_auth, err := ctx.Auth.UserAuthorized(rt.db, comment_owner)

	if err != nil {
		helpers.SendInternalError(err, "Error checking authorization", w, rt.baseLogger)
		return
	}

	// If the status is UNAUTHORIZED, this means that the Authorization header is missing or invalid
	// We don't need to check if user is 'uid' and we can send the error
	if owner_auth == reqcontext.UNAUTHORIZED {
		helpers.SendStatus(http.StatusUnauthorized, w, "Unauthorized", rt.baseLogger)
	}

	if owner_auth != reqcontext.AUTHORIZED {
		// Authorized user is not the owner of the comment
		// let's check if it's the owner of the photo

		if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusForbidden) {
			// The authorized user is not the owner of the photo, so we sent an error
			return
		}
		// If it is authorized, the user can delete the comment.
		//		(else the status must be FORBIDDEN. It can't be UNAUTHORIZED because we already checked it
		//		and it can't be NOT_FOUND because we used it before to get the comment owner)
	}

	// Delete the comment
	_, err = rt.db.DeleteComment(uid, photo_id, comment_id)

	if err != nil {
		helpers.SendInternalError(err, "Database error: DeleteComment", w, rt.baseLogger)
		return
	}

	// We don't need to check the status because if the comment didn't exist
	// we'd have already returned an error when getting the comment owner
	// so we know the comment existed and was deleted, and we can safely send 204

	w.WriteHeader(http.StatusNoContent)
}
