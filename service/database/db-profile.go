package database

import (
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

//this should be changed, but we need to change OpenAPI first

// Get user profile, including username, followers, following, and photos
func (db *appdbimpl) GetUserProfile(uid string, requesting_uid string) (QueryResult, *structures.UserProfile, error) {
	// Get user info
	var name string
	err := db.c.QueryRow(`SELECT "name" FROM "users" WHERE "uid" = ?`, uid).Scan(&name)

	if db_errors.EmptySet(err) {
		// Query returned no rows, the user does not exist
		return ERR_NOT_FOUND, nil, nil

	} else if err != nil {
		return ERR_INTERNAL, nil, err
	}

	// Get followers
	var followers int64
	err = db.c.QueryRow(`SELECT COUNT(*) FROM "follows" WHERE "followed" = ?`, uid).Scan(&followers)

	if err != nil {
		// Return the error
		return ERR_INTERNAL, nil, err
	}

	// Get following users
	var following int64
	err = db.c.QueryRow(`SELECT COUNT(*) FROM "follows" WHERE "follower" = ?`, uid).Scan(&following)

	if err != nil {
		// Return the error
		return ERR_INTERNAL, nil, err
	}

	var photos int64
	err = db.c.QueryRow(`SELECT COUNT(*) FROM "photos" WHERE "photos"."user" = ?`, uid).Scan(&photos)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	// Get follow status
	var follow_status bool
	err = db.c.QueryRow(`SELECT EXISTS (SELECT * FROM "follows" WHERE "follower" = ? AND "followed" = ?)`, requesting_uid, uid).Scan(&follow_status)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	// Get ban status
	var ban_status bool
	err = db.c.QueryRow(`SELECT EXISTS (SELECT * FROM "bans" WHERE "user" = ? AND "ban" = ?)`, requesting_uid, uid).Scan(&ban_status)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	return SUCCESS, &structures.UserProfile{
		UID:       uid,
		Name:      name,
		Following: following,
		Followers: followers,
		Followed:  follow_status,
		Banned:    ban_status,
		Photos:    photos,
	}, nil
}

func (db *appdbimpl) GetUserPhotos(uid string, requesting_uid string, start_index int, limit int) (*[]structures.UserPhoto, error) {

	// Get photos
	rows, err := db.c.Query(`SELECT "p"."id", "p"."date",
								(
									SELECT COUNT(*) AS "likes" FROM "likes" AS "l"
									WHERE "l"."photo_id" = "p"."id"
								),
								(
									SELECT COUNT(*) AS "comments" FROM "comments" AS "c"
									WHERE "c"."photo" = "p"."id"
								),
								EXISTS (
									SELECT * FROM "likes" AS "l"
									WHERE "l"."photo_id" = "p"."id"
									AND "l"."user" = ?
								)
 								FROM "photos" AS "p"
								WHERE "p"."user" = ?
								ORDER BY "p"."date" DESC
								LIMIT ?
								OFFSET ?`, requesting_uid, uid, limit, start_index)
	if err != nil {
		// Return the error
		return nil, err
	}

	photos := make([]structures.UserPhoto, 0)

	defer rows.Close()

	for rows.Next() {
		// If there is a next row, we create an instance of Photo and add it to the slice
		var photo structures.UserPhoto
		err = rows.Scan(&photo.ID, &photo.Date, &photo.Likes, &photo.Comments, &photo.Liked)
		if err != nil {
			// Return the error
			return nil, err
		}
		photos = append(photos, photo)
	}
	// We check if the iteration ended prematurely
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &photos, nil
}
