package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) PostPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	defer r.Body.Close()

	uid := ps.ByName("user_id")

	// send error if the user has no permission to perform this action
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, http.StatusNotFound) {
		return
	}

	transaction, photo_id, err := rt.db.PostPhoto(uid)

	if err != nil {
		helpers.SendInternalError(err, "Database error: PostPhoto", w, rt.baseLogger)
		return
	}

	path := rt.dataPath + "/photos/" + uid + "/" + strconv.FormatInt(photo_id, 10) + ".jpg"
	// todo: we should check if the body is a valid jpg image

	if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil { // perms = 511
		helpers.SendInternalError(err, "Error creating directory", w, rt.baseLogger)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		helpers.SendInternalError(err, "Error creating file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
		return
	}

	if _, err = io.Copy(file, r.Body); err != nil {
		helpers.SendInternalError(err, "Error writing the file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
		return
	}

	if err = file.Close(); err != nil {
		helpers.SendInternalError(err, "Error closing file", w, rt.baseLogger)
		helpers.RollbackOrLogError(transaction, rt.baseLogger)
	}

	err = transaction.Commit()

	if err != nil {
		helpers.SendInternalError(err, "Error committing transaction", w, rt.baseLogger)
		//todo: should I roll back?
		return
	}

	helpers.SendStatus(http.StatusCreated, w, "Photo uploaded", rt.baseLogger)
}

func (rt *_router) GetPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	photo_id := ps.ByName("photo_id")

	if !helpers.VerifyUserOrNotFound(rt.db, uid, w, rt.baseLogger) {
		return
	}

	path := rt.dataPath + "/photos/" + uid + "/" + photo_id + ".jpg"

	file, err := os.Open(path)

	if err != nil {
		helpers.SendNotFound(w, "Photo not found", rt.baseLogger)
		return
	}

	defer file.Close()

	io.Copy(w, file)
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
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, http.StatusNotFound) {
		return
	}

	deleted, err := rt.db.DeletePhoto(uid, photo_id_int)
	if err != nil {
		helpers.SendInternalError(err, "Error deleting photo from database", w, rt.baseLogger)
		return
	} //todo: maybe let's use a transaction also here

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
