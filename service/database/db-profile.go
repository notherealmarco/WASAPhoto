package database

import (
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

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

	// Convert []Photo to []UserPhoto
	user_photos := make([]structures.UserPhoto, 0)
	for _, photo := range *photos {
		user_photos = append(user_photos, structures.UserPhoto{
			ID:       photo.ID,
			Likes:    photo.Likes,
			Comments: photo.Comments,
			Date:     photo.Date,
		})
	}

	return SUCCESS, &structures.UserProfile{
		UID:       uid,
		Name:      name,
		Following: following,
		Followers: followers,
		Photos:    &user_photos}, nil
}
