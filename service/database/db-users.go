package database

import (
	"github.com/gofrs/uuid"
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

//Check if user exists and if exists return the user id by username
//todo

// Check if user exists
func (db *appdbimpl) UserExists(uid string) (bool, error) {
	var name string
	err := db.c.QueryRow(`SELECT "name" FROM "users" WHERE "uid" = ?`, uid).Scan(&name)

	if db_errors.EmptySet(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Get user id by username
func (db *appdbimpl) GetUserID(name string) (string, error) {
	var uid string
	err := db.c.QueryRow(`SELECT "uid" FROM "users" WHERE "name" = ?`, name).Scan(&uid)
	return uid, err
}

// Create a new user
func (db *appdbimpl) CreateUser(name string) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	_, err = db.c.Exec(`INSERT INTO "users" ("uid", "name") VALUES (?, ?)`, uid.String(), name)
	return uid.String(), err
}

// Update username
func (db *appdbimpl) UpdateUsername(uid string, name string) error {
	_, err := db.c.Exec(`UPDATE "users" SET "name" = ? WHERE "uid" = ?`, name, uid)
	return err
}

// Get user followers
func (db *appdbimpl) GetUserFollowers(uid string) ([]structures.UIDName, error) {
	rows, err := db.c.Query(`SELECT "follower", "user.name" FROM "follows", "users"
							WHERE "follows.follower" = "users.uid"
							AND "followed" = ?`, uid)
	if err != nil {
		return nil, err
	}

	var followers []structures.UIDName = make([]structures.UIDName, 0)
	for rows.Next() {
		var uid string
		var name string
		err = rows.Scan(&uid, &name)
		if err != nil {
			return nil, err
		}
		followers = append(followers, structures.UIDName{UID: uid, Name: name})
	}
	return followers, nil
}

// Follow a user
func (db *appdbimpl) FollowUser(uid string, follow string) error {
	_, err := db.c.Exec(`INSERT INTO "follows" ("follower", "followed") VALUES (?, ?)`, uid, follow)
	return err
}

// Unfollow a user
func (db *appdbimpl) UnfollowUser(uid string, unfollow string) error {
	_, err := db.c.Exec(`DELETE FROM "follows" WHERE "follower" = ? AND "followed" = ?`, uid, unfollow)
	return err
} //todo: should return boolean or something similar

// Ban a user
func (db *appdbimpl) BanUser(uid string, ban string) error {
	_, err := db.c.Exec(`INSERT INTO "bans" ("user", "ban") VALUES (?, ?)`, uid, ban)
	return err
}

// Unban a user
func (db *appdbimpl) UnbanUser(uid string, unban string) error {
	_, err := db.c.Exec(`DELETE FROM "bans" WHERE "user" = ? AND "ban" = ?`, uid, unban)
	return err
}
