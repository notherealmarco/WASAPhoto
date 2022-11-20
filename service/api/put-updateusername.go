package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/authorization"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

func (rt *_router) UpdateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	if !authorization.SendAuthorizationError(ctx.Auth.UserAuthorized, uid, rt.db, w, http.StatusNotFound) {
		return
	}
	var req structures.UserDetails
	if !helpers.DecodeJsonOrBadRequest(r.Body, w, &req, rt.baseLogger) {
		return
	}

	//err := json.NewDecoder(r.Body).Decode(&req) //todo: capire se serve close

	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest) // todo: move to DecodeOrBadRequest helper
	//	return
	//}

	err := rt.db.UpdateUsername(uid, req.Name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // todo: is not ok, maybe let's use a helper
		return
	}

	w.WriteHeader(http.StatusNoContent) // todo: change to 204 also in API spec
}
