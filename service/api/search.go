package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) GetSearchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// We require user to be authenticated
	if !authorization.SendErrorIfNotLoggedIn(ctx.Auth.Authorized, rt.db, w, rt.baseLogger) {
		return
	}

	// Get search query
	query := r.URL.Query().Get("query")

	if query == "" {
		helpers.SendBadRequest(w, "Missing query parameter", rt.baseLogger)
		return
	}

	// Get start index and limit, or their default values
	start_index, limit, err := helpers.GetLimits(r.URL.Query())

	if err != nil {
		helpers.SendBadRequest(w, "Invalid start_index or limit value", rt.baseLogger)
		return
	}

	// Get search results
	results, err := rt.db.SearchByName(query, ctx.Auth.GetUserID(), start_index, limit)

	if err != nil {
		helpers.SendInternalError(err, "Database error: SearchByName", w, rt.baseLogger)
		return
	}

	// Return the results in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)

	if err != nil {
		helpers.SendInternalError(err, "Error encoding json", w, rt.baseLogger)
		return
	}
}
