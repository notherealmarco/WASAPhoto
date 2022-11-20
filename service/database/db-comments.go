package database

import (
	"time"

	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

func (db *appdbimpl) PostComment(uid string, photo_id int64, comment_user string, comment string) (QueryResult, error) {

	// Check if the photo exists, as API specification requires
	// photos to be identified also by the user who posted them.
	// But our DB implementation only requires the photo id.
	exists, err := db.photoExists(uid, photo_id)
	if err != nil || !exists {
		return ERR_NOT_FOUND, err
	}

	_, err = db.c.Exec(`PRAGMA foreign_keys = ON;
						INSERT INTO "comments" ("user", "photo", "comment", "date") VALUES (?, ?, ?, ?)`, comment_user, photo_id, comment, time.Now().Format(time.RFC3339))

	// todo: we don't actually need it, it's already done before
	if db_errors.ForeignKeyViolation(err) {
		return ERR_NOT_FOUND, nil
	}

	if err != nil {
		return ERR_INTERNAL, err
	}
	return SUCCESS, nil
}

func (db *appdbimpl) GetCommentOwner(uid string, photo_id int64, comment_id int64) (QueryResult, string, error) {

	// Check if the photo exists, as it exist but have no comments
	exists, err := db.photoExists(uid, photo_id)
	if err != nil || !exists {
		return ERR_NOT_FOUND, "", err
	}

	var comment_user string
	err = db.c.QueryRow(`SELECT "user" FROM "comments" WHERE "photo" = ? AND "id" = ?`, photo_id, comment_id).Scan(&comment_user)

	if db_errors.EmptySet(err) {
		return ERR_NOT_FOUND, "", nil
	}

	if err != nil {
		return ERR_INTERNAL, "", err
	}

	return SUCCESS, comment_user, nil
}

func (db *appdbimpl) DeleteComment(uid string, photo_id int64, comment_id int64) (QueryResult, error) {

	// Check if the photo exists, as API specification requires
	// photos to be identified also by the user who posted them.
	// But our DB implementation only requires the photo id.
	exists, err := db.photoExists(uid, photo_id)
	if err != nil || !exists {
		return ERR_NOT_FOUND, err
	}

	res, err := db.c.Exec(`DELETE FROM "comments" WHERE "photo" = ? AND "id" = ?`, photo_id, comment_id)

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

func (db *appdbimpl) GetComments(uid string, photo_id int64) (QueryResult, *[]structures.Comment, error) {

	// Check if the photo exists, as it exist but have no comments
	exists, err := db.photoExists(uid, photo_id)
	if err != nil || !exists {
		return ERR_NOT_FOUND, nil, err
	}

	rows, err := db.c.Query(`SELECT "id", "user", "comment", "date" FROM "comments" WHERE "photo" = ?`, photo_id)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	defer rows.Close()

	comments := make([]structures.Comment, 0)

	for rows.Next() {
		var c structures.Comment
		err = rows.Scan(&c.CommentID, &c.UID, &c.Comment, &c.Date)
		if err != nil {
			return ERR_INTERNAL, nil, err
		}
		comments = append(comments, c)
	}

	return SUCCESS, &comments, nil
}
