package database

import (
	"database/sql"
	"time"

	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

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

type QueryResult int // todo: move to a separate file

const (
	SUCCESS       = 0
	ERR_NOT_FOUND = 1
	ERR_EXISTS    = 2
	ERR_INTERNAL  = 3
)

type dbtransaction struct {
	c *sql.Tx
}

func (tx *dbtransaction) Commit() error {
	return tx.c.Commit()
}

func (tx *dbtransaction) Rollback() error {
	return tx.c.Rollback()
}

// Post a new photo
func (db *appdbimpl) PostPhoto(uid string) (DBTransaction, int64, error) {
	tx, err := db.c.Begin()

	if err != nil {
		return nil, 0, err
	}

	res, err := tx.Exec(`INSERT INTO "photos" ("user", "date") VALUES (?, ?)`, uid, time.Now().Format(time.RFC3339))
	if err != nil {
		tx.Rollback() // error ?
		return nil, 0, err
	}
	id, err := res.LastInsertId()

	if err != nil {
		tx.Rollback() // error ?
		return nil, 0, err
	}

	return &dbtransaction{
		c: tx,
	}, id, nil
}

// Delete a photo, returns true if the photo was deleted and false if it did not exist
func (db *appdbimpl) DeletePhoto(uid string, photo int64) (bool, error) {
	res, err := db.c.Exec(`DELETE FROM "photos" WHERE "id" = ? AND "user" = ?`, photo, uid)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
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

	photos := make([]Photo, 0)
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

// Check if a given photo owned by a given user exists
func (db *appdbimpl) photoExists(uid string, photo int64) (bool, error) {

	var cnt int64
	err := db.c.QueryRow(`SELECT COUNT(*) FROM "photos" WHERE "id" = ? AND "user" = ?`, photo, uid).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// Get the list of users who liked a photo
func (db *appdbimpl) GetPhotoLikes(uid string, photo int64) (QueryResult, *[]structures.UIDName, error) {

	// Check if the photo exists, as it could exist but have no likes
	exists, err := db.photoExists(uid, photo)
	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	if !exists {
		return ERR_NOT_FOUND, nil, nil
	}

	rows, err := db.c.Query(`SELECT "users.uid", "users.name" FROM "likes", "users"
								WHERE "likes.photo_id" = ?
								AND "likes.user" = "users.uid"`, photo)
	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	likes := make([]structures.UIDName, 0)
	for rows.Next() {
		var uid string
		var name string
		err = rows.Scan(&uid, &name)
		if err != nil {
			return ERR_INTERNAL, nil, err
		}
		likes = append(likes, structures.UIDName{UID: uid, Name: name})
	}

	return SUCCESS, &likes, nil
}

// Like a photo
func (db *appdbimpl) LikePhoto(uid string, photo int64, liker_uid string) (QueryResult, error) {

	// Check if the photo exists, as API specification requires
	// photos to be identified also by the user who posted them.
	// But our DB implementation only requires the photo id.
	exists, err := db.photoExists(uid, photo)
	if err != nil || !exists {
		return ERR_NOT_FOUND, err
	}

	_, err = db.c.Exec(`PRAGMA foreign_keys = ON;
						INSERT INTO "likes" ("user", "photo_id") VALUES (?, ?)`, liker_uid, photo)

	// The photo exists, but the user already liked it
	if db_errors.UniqueViolation(err) {
		return ERR_EXISTS, nil
	}

	if db_errors.ForeignKeyViolation(err) {
		return ERR_NOT_FOUND, nil
	}

	if err != nil {
		return ERR_INTERNAL, err
	}
	return SUCCESS, nil
}

// Unlike a photo
func (db *appdbimpl) UnlikePhoto(uid string, photo int64, liker_uid string) (QueryResult, error) {

	// Check if the photo exists, as API specification requires
	// photos to be identified also by the user who posted them.
	// But our DB implementation only requires the photo id.
	exists, err := db.photoExists(uid, photo)
	if err != nil || !exists {
		return ERR_NOT_FOUND, err
	}

	res, err := db.c.Exec(`DELETE FROM "likes" WHERE "user" = ? AND "photo_id" = ?`, liker_uid, photo)

	if err != nil {
		return ERR_INTERNAL, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return ERR_INTERNAL, err
	}

	if rows == 0 {
		return ERR_NOT_FOUND, nil
	}
	return SUCCESS, nil
}
