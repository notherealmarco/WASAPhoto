package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) PostPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// defer r.Body.Close()

	uid := ps.ByName("user_id")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	transaction, photo_id, err := rt.db.PostPhoto(uid)

	if err != nil {
		helpers.SendInternalError(err, "Database error: PostPhoto", w, rt.baseLogger)
		return
	}

	path := rt.dataPath + "/photos/" + uid + "/" + strconv.FormatInt(photo_id, 10) + ".jpg"

	if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil { // perms = 511
		helpers.SendInternalError(err, "Error creating directory", w, rt.baseLogger)
		return
	}

	/*file, err := os.Create(path)
	if err != nil {
		helpers.SendInternalError(err, "Error creating file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
		return
	}*/

	bytes, err := io.ReadAll(r.Body)

	if err != nil {
		helpers.SendInternalError(err, "Error checking the file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
		return
	}

	mimeType := http.DetectContentType(bytes)

	if !strings.HasPrefix(mimeType, "image/") {
		helpers.SendStatus(http.StatusBadRequest, w, mimeType+" file is not a valid image", rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
		return
	}

	if err = os.WriteFile(path, bytes, 0644); err != nil {
		helpers.SendInternalError(err, "Error writing the file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
		return
	}

	/*if err = file.Close(); err != nil {
		helpers.SendInternalError(err, "Error closing file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
	}*/

	err = transaction.Commit()

	if err != nil {
		helpers.SendInternalError(err, "Error committing transaction", w, rt.baseLogger)
		return
	}

	helpers.SendStatus(http.StatusCreated, w, "Photo uploaded", rt.baseLogger)
}

func (rt *_router) GetPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) {
		// We want the user to be authenticated
		return
	}

	uid := ps.ByName("user_id")

	photo_id_str := ps.ByName("photo_id")
	photo_id, err := strconv.ParseInt(photo_id_str, 10, 64)

	if err != nil {
		helpers.SendBadRequest(w, "Invalid photo id", rt.baseLogger)
		return
	}

	// This is also checking if the requesting user is banned by the author of the photo
	exists, err := rt.db.PhotoExists(uid, photo_id, ctx.Auth.GetUserID())

	if err != nil {
		helpers.SendInternalError(err, "Database error: PhotoExists", w, rt.baseLogger)
		return
	}

	if !exists {
		helpers.SendNotFound(w, "Resource not found", rt.baseLogger)
		return
	}

	path := rt.dataPath + "/photos/" + uid + "/" + photo_id_str + ".jpg"

	file, err := os.Open(path)

	if err != nil {
		helpers.SendNotFound(w, "Photo not found", rt.baseLogger)
		return
	}

	defer file.Close()

	_, err = io.Copy(w, file)

	if err != nil {
		helpers.SendInternalError(err, "Error writing response", w, rt.baseLogger)
		return
	}
}

func (rt *_router) DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	photo_id := ps.ByName("photo_id")

	// photo id to int64
	photo_id_int, err := strconv.ParseInt(photo_id, 10, 64)

	if err != nil {
		helpers.SendBadRequestError(err, "Bad photo id", w, rt.baseLogger)
		return
	}

	// send error if the user has no permission to perform this action (only the author can delete a photo)
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, rt.baseLogger, http.StatusNotFound) {
		return
	}

	deleted, err := rt.db.DeletePhoto(uid, photo_id_int)
	if err != nil {
		helpers.SendInternalError(err, "Error deleting photo from database", w, rt.baseLogger)
		return
	}

	if !deleted {
		helpers.SendNotFound(w, "Photo not found", rt.baseLogger)
		return
	}

	path := rt.dataPath + "/photos/" + uid + "/" + photo_id + ".jpg"

	if err := os.Remove(path); err != nil {
		helpers.SendInternalError(err, "Error deleting file", w, rt.baseLogger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
