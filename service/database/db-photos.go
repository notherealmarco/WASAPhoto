package database

import (
	"time"

	"github.com/notherealmarco/WASAPhoto/service/structures"
)

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

func (db *appdbimpl) getUserPhotos(uid string) (*[]structures.Photo, error) {

	// Get photos
	rows, err := db.c.Query(`SELECT "p"."user", "p"."id", "p"."date",
								(
									SELECT COUNT(*) AS "likes" FROM "likes" AS "l"
									WHERE "l"."photo_id" = "p"."id"
								),
								(
									SELECT COUNT(*) AS "comments" FROM "comments" AS "c"
									WHERE "c"."photo" = "p"."id"
								)
 								FROM "photos" AS "p"
								WHERE "p"."user" = ?`, uid)
	if err != nil {
		// Return the error
		return nil, err
	}

	photos := make([]structures.Photo, 0)

	for rows.Next() {
		// If there is a next row, we create an instance of Photo and add it to the slice
		var photo structures.Photo
		err = rows.Scan(&photo.UID, &photo.ID, &photo.Date, &photo.Likes, &photo.Comments)
		if err != nil {
			// Return the error
			return nil, err
		}
		photos = append(photos, photo)
	}

	return &photos, nil
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
