package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/notherealmarco/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) UpdateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	auth, err := ctx.Auth.UserAuthorized(rt.db, r.URL.Path // todo: prendere il coso giusto dal path)

}
