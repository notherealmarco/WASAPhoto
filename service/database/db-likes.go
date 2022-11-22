package database

import (
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

// Get the list of users who liked a photo
func (db *appdbimpl) GetPhotoLikes(uid string, photo int64, requesting_uid string, start_index int, limit int) (QueryResult, *[]structures.UIDName, error) {

	// Check if the photo exists, as it could exist but have no likes
	exists, err := db.photoExists(uid, photo)
	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	if !exists {
		return ERR_NOT_FOUND, nil, nil
	}

	rows, err := db.c.Query(`SELECT "users"."uid", "users"."name" FROM "likes", "users"
								WHERE "likes"."photo_id" = ?
								AND "likes"."user" NOT IN (
									SELECT "bans"."user" FROM "bans"
									WHERE "bans"."user" = "likes"."user"
									AND "bans"."ban" = ?
								)
								AND "likes"."user" = "users"."uid"
								LIMIT ?
								OFFSET ?`, photo, requesting_uid, limit, start_index)
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
