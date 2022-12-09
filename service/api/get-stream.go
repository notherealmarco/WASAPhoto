package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) GetUserStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// We must know who is requesting the stream
	if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) {
		return
	}

	// Get user id, probably this should be changed as it's not really REST
	uid := ctx.Auth.GetUserID()

	// Get start index and limit, or their default values
	start_index, limit, err := helpers.GetLimits(r.URL.Query())

	if err != nil {
		helpers.SendBadRequest(w, "Invalid start_index or limit value", rt.baseLogger)
		return
	}

	// Get the stream
	stream, err := rt.db.GetUserStream(uid, start_index, limit)

	if err != nil {
		helpers.SendInternalError(err, "Database error: GetUserStream", w, rt.baseLogger)
		return
	}

	// Return the stream in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(stream)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}
