package api

import (
	"encoding/json"
	"net/http"
	"regexp"

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

	var uid string
	if err == nil { // test if user exists
		uid, err = rt.db.GetUserID(request.Name)
	}

	if db_errors.EmptySet(err) { // user does not exist

		// before creating the user, check if the name is valid
		stat, regex_err := regexp.Match(`^[a-zA-Z0-9_]{3,16}$`, []byte(request.Name))
		if regex_err != nil {
			helpers.SendInternalError(err, "Error while matching username regex", w, rt.baseLogger)
			return
		}
		if !stat {
			// username didn't match the regex, so it's invalid, let's send a bad request error
			helpers.SendBadRequest(w, "Username must be between 3 and 16 characters long and can only contain letters, numbers and underscores", rt.baseLogger)
			return
		}

		uid, err = rt.db.CreateUser(request.Name)
	}
	if err != nil { // handle any other error
		helpers.SendBadRequestError(err, "Bad request body", w, rt.baseLogger)
		return
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(_respbody{UID: uid})

	if err != nil {
		helpers.SendInternalError(err, "Error encoding response", w, rt.baseLogger)
		return
	}
}
