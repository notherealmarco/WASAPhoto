package database

import (
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

// Get user stream
// todo implement stidx + offset
func (db *appdbimpl) GetUserStream(uid string, start_index int, limit int) (*[]structures.Photo, error) {

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
								WHERE "p"."user" IN (
									SELECT "followed" FROM "follows" WHERE "follower" = ?
								)
								AND "p"."user" NOT IN (
									SELECT "user" FROM "bans" WHERE "ban" = ?
								)
								ORDER BY "p"."date" DESC
								OFFSET ?
								LIMIT ?`, uid, uid, start_index, limit)
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
