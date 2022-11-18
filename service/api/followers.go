package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/helpers"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) GetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")

	if !helpers.VerifyUserOrNotFound(rt.db, uid, w) {
		return
	}

	followers, err := rt.db.GetUserFollowers(uid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // todo: is not ok, maybe let's use a helper
		return
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(&followers)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // todo: is not ok
		return
	}
}

func (rt *_router) PutFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	uid := ps.ByName("user_id")
	followed := ps.ByName("follower_uid")

	err := rt.db.FollowUser(uid, followed)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // todo: is not ok, maybe let's use a helper
		return
	}

	w.WriteHeader(http.StatusNoContent) // todo: change to 204 also in API spec
}
