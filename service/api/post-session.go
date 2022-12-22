package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
)

type _reqbody struct {
	Name string `json:"name"`
}

type _respbody struct {
	UID string `json:"user_id"`
}

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) PostSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var request _reqbody
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		helpers.SendBadRequestError(err, "Bad request body", w, rt.baseLogger)
		return
	}

	// test if user exists
	var uid string
	uid, err = rt.db.GetUserID(request.Name)

	// check if the database returned an empty set error, if so, create the new user
	if db_errors.EmptySet(err) {

		// before creating the user, check if the name is valid, otherwise send a bad request error
		if !helpers.MatchUsernameOrBadRequest(request.Name, w, rt.baseLogger) {
			return
		}

		uid, err = rt.db.CreateUser(request.Name)
	}

	// handle database errors
	if err != nil {
		helpers.SendInternalError(err, "Database error", w, rt.baseLogger)
		return
	}

	// set the response header
	w.Header().Set("content-type", "application/json")

	// encode the response body
	err = json.NewEncoder(w).Encode(_respbody{UID: uid})

	if err != nil {
		helpers.SendInternalError(err, "Error encoding response", w, rt.baseLogger)
		return
	}
}
