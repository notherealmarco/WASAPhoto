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
	//
	// This also checks if the author has banned the user who is posting the comment
	// as he should not be able to post comments on his photos
	exists, err := db.PhotoExists(uid, photo_id, comment_user)
	if err != nil || !exists {
		return ERR_NOT_FOUND, err
	}

	_, err = db.c.Exec(`PRAGMA foreign_keys = ON;
						INSERT INTO "comments" ("user", "photo", "comment", "date") VALUES (?, ?, ?, ?)`, comment_user, photo_id, comment, time.Now().Format(time.RFC3339))

	if db_errors.ForeignKeyViolation(err) {
		// trying to post a comment on a photo that does not exist
		// (actually this should never happen, as we checked if the photo exists before)
		return ERR_NOT_FOUND, nil
	}

	if err != nil {
		return ERR_INTERNAL, err
	}
	return SUCCESS, nil
}

func (db *appdbimpl) GetCommentOwner(uid string, photo_id int64, comment_id int64) (QueryResult, string, error) {

	// Check if the photo exists, as it may exist but have no comments
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

func (db *appdbimpl) GetComments(uid string, photo_id int64, requesting_uid string, start_index int, limit int) (QueryResult, *[]structures.Comment, error) {

	// Check if the photo exists, as it exist but have no comments
	// this also checks if the author has banned the requesting user
	exists, err := db.PhotoExists(uid, photo_id, requesting_uid)
	if err != nil || !exists {
		return ERR_NOT_FOUND, nil, err
	}

	rows, err := db.c.Query(`SELECT "c"."id", "c"."user", "c"."comment", "c"."date", "u"."name"
								FROM "comments" AS "c", "users" AS "u"
								WHERE "c"."photo" = ?
								AND "c"."user" NOT IN (
									SELECT "bans"."user" FROM "bans"
									WHERE "bans"."user" = "c"."user"
									AND "bans"."ban" = ?
								)
								AND "u"."uid" = "c"."user"
								LIMIT ?
								OFFSET ?`, photo_id, requesting_uid, limit, start_index)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	defer rows.Close()

	comments := make([]structures.Comment, 0)

	defer rows.Close()

	for rows.Next() {
		var c structures.Comment
		err = rows.Scan(&c.CommentID, &c.UID, &c.Comment, &c.Date, &c.Name)
		if err != nil {
			return ERR_INTERNAL, nil, err
		}
		comments = append(comments, c)
	}

	return SUCCESS, &comments, nil
}
