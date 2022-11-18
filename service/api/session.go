package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
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
	json.NewDecoder(r.Body).Decode(&request) //todo: capire se serve close

	uid, err := rt.db.GetUserID(request.Name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(_respbody{UID: uid})
}
