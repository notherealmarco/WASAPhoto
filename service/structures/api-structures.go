package structures

type UserDetails struct {
	Name string `json:"name"`
}

type UIDName struct {
	UID  string `json:"user_id"`
	Name string `json:"name"`
}

type GenericResponse struct {
	Status string `json:"status"`
}

type Comment struct {
	CommentID string `json:"comment_id"`
	UID       string `json:"user_id"`
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	Date      string `json:"date"`
}
