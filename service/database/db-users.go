package database

import (
	"database/sql"

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
func (db *appdbimpl) GetUserFollowers(uid string) (QueryResult, *[]structures.UIDName, error) {

	// user may exist but have no followers
	exists, err := db.UserExists(uid)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	if !exists {
		return ERR_NOT_FOUND, nil, nil
	}

	rows, err := db.c.Query(`SELECT "follower", "user"."name" FROM "follows", "users"
							WHERE "follows"."follower" = "users"."uid"
							AND "followed" = ?`, uid)

	followers, err := db.uidNameQuery(rows, err)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	return SUCCESS, followers, nil
}

// Get user following
func (db *appdbimpl) GetUserFollowing(uid string) (QueryResult, *[]structures.UIDName, error) {

	// user may exist but have no followers
	exists, err := db.UserExists(uid)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	if !exists {
		return ERR_NOT_FOUND, nil, nil
	}

	rows, err := db.c.Query(`SELECT "followed", "user"."name" FROM "follows", "users"
							WHERE "follows"."followed" = "users"."uid"
							AND "follower" = ?`, uid)

	following, err := db.uidNameQuery(rows, err)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	return SUCCESS, following, nil
}

// Evaluates a query that returns two columns: uid and name
func (db *appdbimpl) uidNameQuery(rows *sql.Rows, err error) (*[]structures.UIDName, error) {

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
	return &followers, nil
}

// Follow a user
func (db *appdbimpl) FollowUser(uid string, follow string) (QueryResult, error) {
	_, err := db.c.Exec(`PRAGMA foreign_keys = ON;
						INSERT INTO "follows" ("follower", "followed") VALUES (?, ?)`, uid, follow)

	if db_errors.UniqueViolation(err) {
		return ERR_EXISTS, nil
	}

	if db_errors.ForeignKeyViolation(err) {
		return ERR_NOT_FOUND, nil
	}

	if err != nil {
		return ERR_INTERNAL, err
	}

	return SUCCESS, nil
}

// Unfollow a user
func (db *appdbimpl) UnfollowUser(uid string, unfollow string) (QueryResult, error) {
	res, err := db.c.Exec(`DELETE FROM "follows" WHERE "follower" = ? AND "followed" = ?`, uid, unfollow)

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

// Ban a user
func (db *appdbimpl) BanUser(uid string, ban string) (QueryResult, error) {
	_, err := db.c.Exec(`PRAGMA foreign_keys = ON;
						INSERT INTO "bans" ("user", "ban") VALUES (?, ?)`, uid, ban)

	// The user is already banned by this user
	if db_errors.UniqueViolation(err) {
		return ERR_EXISTS, nil
	}

	// One of the users does not exist
	if db_errors.ForeignKeyViolation(err) {
		return ERR_NOT_FOUND, nil
	}

	// Other error
	if err != nil {
		return ERR_INTERNAL, err
	}
	return SUCCESS, nil
}

// Unban a user
func (db *appdbimpl) UnbanUser(uid string, unban string) (QueryResult, error) {
	res, err := db.c.Exec(`DELETE FROM "bans" WHERE "user" = ? AND "ban" = ?`, uid, unban)

	if err != nil {
		return ERR_INTERNAL, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return ERR_INTERNAL, err
	}

	// The user was not banned by this user
	if rows == 0 {
		return ERR_NOT_FOUND, nil
	}
	return SUCCESS, nil
}
