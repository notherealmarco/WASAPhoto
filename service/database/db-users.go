package database

import (
	"github.com/gofrs/uuid"
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
)

//Check if user exists and if exists return the user id by username
//todo

// Check if user exists
func (db *appdbimpl) UserExists(uid string) (bool, error) {
	var name string
	err := db.c.QueryRow(`SELECT "name" FROM "users" WHERE "uid" = ?`, name).Scan(&name)

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
