package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.PostSession))

	rt.router.PUT("/users/:user_id/username", rt.wrap(rt.UpdateUsername))

	rt.router.GET("/users", rt.wrap(rt.GetSearchUsers))

	rt.router.GET("/users/:user_id/followers", rt.wrap(rt.GetFollowersFollowing))
	rt.router.PUT("/users/:user_id/followers/:follower_uid", rt.wrap(rt.PutFollow))
	rt.router.DELETE("/users/:user_id/followers/:follower_uid", rt.wrap(rt.DeleteFollow))
	rt.router.GET("/users/:user_id/following", rt.wrap(rt.GetFollowersFollowing))

	rt.router.GET("/users/:user_id/bans", rt.wrap(rt.GetUserBans))
	rt.router.PUT("/users/:user_id/bans/:ban_uid", rt.wrap(rt.PutBan))
	rt.router.DELETE("/users/:user_id/bans/:ban_uid", rt.wrap(rt.DeleteBan))

	rt.router.POST("/users/:user_id/photos", rt.wrap(rt.PostPhoto))
	rt.router.GET("/users/:user_id/photos/:photo_id", rt.wrap(rt.GetPhoto))
	rt.router.DELETE("/users/:user_id/photos/:photo_id", rt.wrap(rt.DeletePhoto))

	rt.router.GET("/users/:user_id/photos/:photo_id/likes", rt.wrap(rt.GetLikes))
	rt.router.PUT("/users/:user_id/photos/:photo_id/likes/:liker_uid", rt.wrap(rt.PutDeleteLike))
	rt.router.DELETE("/users/:user_id/photos/:photo_id/likes/:liker_uid", rt.wrap(rt.PutDeleteLike))

	rt.router.GET("/users/:user_id/photos/:photo_id/comments", rt.wrap(rt.GetComments))
	rt.router.POST("/users/:user_id/photos/:photo_id/comments", rt.wrap(rt.PostComment))
	rt.router.DELETE("/users/:user_id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.DeleteComment))

	rt.router.GET("/users/:user_id", rt.wrap(rt.GetUserProfile))

	rt.router.GET("/stream", rt.wrap(rt.GetUserStream)) //todo: why not "/users/:user_id/stream"?

	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
