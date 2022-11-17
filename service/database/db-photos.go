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
	Following []string
	Followers []string
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
