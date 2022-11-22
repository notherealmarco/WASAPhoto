package database

import (
	"time"
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

// Check if a given photo owned by a given user exists
func (db *appdbimpl) photoExists(uid string, photo int64) (bool, error) {

	var cnt int64
	err := db.c.QueryRow(`SELECT COUNT(*) FROM "photos" WHERE "id" = ? AND "user" = ?`, photo, uid).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
