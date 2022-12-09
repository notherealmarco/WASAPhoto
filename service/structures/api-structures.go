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

type Photo struct {
	UID      string `json:"user_id"`
	Username string `json:"name"`
	ID       int64  `json:"photo_id"`
	Likes    int64  `json:"likes"`
	Comments int64  `json:"comments"`
	Date     string `json:"date"`
	Liked    bool   `json:"liked"`
}

type UserPhoto struct {
	ID       int64  `json:"photo_id"`
	Likes    int64  `json:"likes"`
	Comments int64  `json:"comments"`
	Date     string `json:"date"`
	Liked    bool   `json:"liked"`
}

type UserProfile struct {
	UID       string `json:"user_id"`
	Name      string `json:"name"`
	Following int64  `json:"following"`
	Followers int64  `json:"followers"`
	Followed  bool   `json:"followed"`
	Photos    int64  `json:"photos"`
}
