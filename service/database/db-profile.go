package database

import (
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

//this should be changed, but we need to change OpenAPI first

// Get user profile, including username, followers, following, and photos
func (db *appdbimpl) GetUserProfile(uid string) (QueryResult, *structures.UserProfile, error) {
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

	photos, err := db.getUserPhotos(uid)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	return SUCCESS, &structures.UserProfile{
		UID:       uid,
		Name:      name,
		Following: following,
		Followers: followers,
		Photos:    photos,
	}, nil
}

func (db *appdbimpl) getUserPhotos(uid string) (*[]structures.UserPhoto, error) {

	// Get photos
	rows, err := db.c.Query(`SELECT "p"."id", "p"."date",
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

	photos := make([]structures.UserPhoto, 0)

	for rows.Next() {
		// If there is a next row, we create an instance of Photo and add it to the slice
		var photo structures.UserPhoto
		err = rows.Scan(&photo.ID, &photo.Date, &photo.Likes, &photo.Comments)
		if err != nil {
			// Return the error
			return nil, err
		}
		photos = append(photos, photo)
	}

	return &photos, nil
}
