package database

import "time"

type Photo struct {
	ID       int64
	Likes    int64
	Comments int64
}

type UserProfile struct {
	UID       string
	Name      string
	Following int64
	Followers int64
	Photos    []Photo
}

// Post a new photo
func (db *appdbimpl) PostPhoto(uid string) (int64, error) {
	res, err := db.c.Exec(`INSERT INTO "photos" ("user", "date") VALUES (?, ?)`, uid, time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Get user profile, including username, followers, following, and photos
func (db *appdbimpl) GetUserProfile(uid string) (*UserProfile, error) {
	// Get user info
	var name string
	err := db.c.QueryRow(`SELECT "name" FROM "users" WHERE "uid" = ?`, uid).Scan(&name)
	if err != nil {
		return nil, err
	}

	// Get followers
	var followers int64
	err = db.c.QueryRow(`SELECT COUNT(*) FROM "follows" WHERE "followed" = ?`, uid).Scan(&followers)

	// Get following users
	var following int64
	err = db.c.QueryRow(`SELECT COUNT(*) FROM "follows" WHERE "follower" = ?`, uid).Scan(&following)

	// Get photos
	rows, err := db.c.Query(`SELECT "photos.id", "photos.date",
								COUNT("likes.user") AS "likes",
								COUNT("comments.user") AS "comments"
								FROM "photos", "likes", "comments"
								WHERE "likes.photo_id" = "photos.id"
								AND "comments.photo" = "photos.id"
								AND "user" = ?`, uid)
	if err != nil {
		return nil, err
	}

	var photos []Photo
	for rows.Next() {
		var id int64
		var date string
		var likes int64
		var comments int64
		err = rows.Scan(&id, &date, &likes, &comments)
		if err != nil {
			return nil, err
		}
		photo_data := Photo{id, likes, comments}
		photos = append(photos, photo_data)
	}

	return &UserProfile{uid, name, followers, following, photos}, nil
}

// Like a photo
func (db *appdbimpl) LikePhoto(uid string, photo int64) error {
	_, err := db.c.Exec(`INSERT INTO "likes" ("user", "photo_id") VALUES (?, ?)`, uid, photo)
	return err
}

// Unlike a photo
func (db *appdbimpl) UnlikePhoto(uid string, photo int64) error {
	_, err := db.c.Exec(`DELETE FROM "likes" WHERE "user" = ? AND "photo_id" = ?`, uid, photo)
	return err
}
