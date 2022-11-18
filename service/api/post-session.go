package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
)

type _reqbody struct {
	Name string `json:"name"`
}

type _respbody struct {
	UID string `json:"uid"`
}

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) PostSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var request _reqbody
	err := json.NewDecoder(r.Body).Decode(&request) //todo: capire se serve close

	var uid string
	if err == nil { // test if user exists
		uid, err = rt.db.GetUserID(request.Name)
	}
	if db_errors.EmptySet(err) { // user does not exist
		err = nil
		uid, err = rt.db.CreateUser(request.Name)
	}
	if err != nil { // handle any other error
		w.WriteHeader(http.StatusInternalServerError) // todo: is not ok
		return
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(_respbody{UID: uid})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // todo: is not ok
		return
	}
}