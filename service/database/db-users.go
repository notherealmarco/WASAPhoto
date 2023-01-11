package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/notherealmarco/WASAPhoto/service/database/db_errors"
	"github.com/notherealmarco/WASAPhoto/service/structures"
)

// Check if user exists
func (db *appdbimpl) UserExists(uid string) (bool, error) {

	var cnt int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM "users" WHERE "uid" = ?`, uid).Scan(&cnt)

	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// User exists and is not banned
func (db *appdbimpl) UserExistsNotBanned(uid string, requesting_uid string) (bool, error) {

	var cnt int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM "users"
							WHERE "uid" = ?
							AND NOT EXISTS (
								SELECT "bans"."user" FROM "bans"
								WHERE "bans"."user" = "users"."uid"
								AND "bans"."ban" = ?
							)`, uid, requesting_uid).Scan(&cnt)

	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// Get user id by username
func (db *appdbimpl) GetUserID(name string) (string, error) {
	var uid string
	err := db.c.QueryRow(`SELECT "uid" FROM "users" WHERE "name" = ?`, name).Scan(&uid)
	return uid, err
}

// Create a new user
func (db *appdbimpl) CreateUser(name string) (string, error) {

	// check if username is taken (case insensitive)
	exists, err := db.nameExists(name)

	if err != nil {
		return "", err
	} else if exists {
		return "", errors.New("username already exists")
	}

	// create new user id
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	// insert the new user into the database
	_, err = db.c.Exec(`INSERT INTO "users" ("uid", "name") VALUES (?, ?)`, uid.String(), name)
	return uid.String(), err
}

// Check if username exists
func (db *appdbimpl) nameExists(name string) (bool, error) {
	var cnt int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM "users" WHERE "name" LIKE ?`, name).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// Update username
func (db *appdbimpl) UpdateUsername(uid string, name string) (QueryResult, error) {

	// check if username is taken (case insensitive)
	exists, err := db.nameExists(name)

	if err != nil {
		return ERR_INTERNAL, err
	} else if exists {
		return ERR_EXISTS, nil
	}

	_, err = db.c.Exec(`UPDATE "users" SET "name" = ? WHERE "uid" = ?`, name, uid)

	if err != nil {
		return ERR_INTERNAL, err
	}

	return SUCCESS, err
}

// Get user followers
func (db *appdbimpl) GetUserFollowers(uid string, requesting_uid string, start_index int, limit int) (QueryResult, *[]structures.UIDName, error) {

	// user may exist but have no followers
	exists, err := db.UserExistsNotBanned(uid, requesting_uid)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	if !exists {
		return ERR_NOT_FOUND, nil, nil
	}

	rows, err := db.c.Query(`SELECT "follower", "users"."name" FROM "follows", "users"
							WHERE "follows"."follower" = "users"."uid"
							
							AND "follows"."follower" NOT IN (
								SELECT "bans"."user" FROM "bans"
								WHERE "bans"."user" = "follows"."follower"
								AND "bans"."ban" = ?
							)

							AND "followed" = ?
							LIMIT ?
							OFFSET ?`, uid, requesting_uid, limit, start_index)

	followers, err := db.uidNameQuery(rows, err)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	return SUCCESS, followers, nil
}

// Get user following
func (db *appdbimpl) GetUserFollowing(uid string, requesting_uid string, start_index int, offset int) (QueryResult, *[]structures.UIDName, error) {

	// user may exist but have no followers
	exists, err := db.UserExistsNotBanned(uid, requesting_uid)

	if err != nil {
		return ERR_INTERNAL, nil, err
	}

	if !exists {
		return ERR_NOT_FOUND, nil, nil
	}

	rows, err := db.c.Query(`SELECT "followed", "user"."name" FROM "follows", "users"
							WHERE "follows"."followed" = "users"."uid"

							AND "follows"."followed" NOT IN (
								SELECT "bans"."user" FROM "bans"
								WHERE "bans"."user" = "follows"."followed"
								AND "bans"."ban" = ?
							)

							AND "follower" = ?
							LIMIT ?
							OFFSET ?`, uid, requesting_uid, offset, start_index)

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

	defer rows.Close()

	for rows.Next() {
		var uid string
		var name string
		err = rows.Scan(&uid, &name)
		if err != nil {
			return nil, err
		}
		followers = append(followers, structures.UIDName{UID: uid, Name: name})
	}
	// We check if the iteration ended prematurely
	if err = rows.Err(); err != nil {
		return nil, err
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

// Is user banned by another user
func (db *appdbimpl) IsBanned(uid string, banner string) (bool, error) {

	var cnt int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM "bans" WHERE "user" = ? AND "ban" = ?`, banner, uid).Scan(&cnt)

	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (db *appdbimpl) GetUserBans(uid string, start_index int, limit int) (*[]structures.UIDName, error) {

	rows, err := db.c.Query(`SELECT "ban", "users"."name" FROM "bans", "users"
							WHERE "bans"."ban" = "users"."uid"
							AND "bans"."user" = ?
							LIMIT ?
							OFFSET ?`, uid, limit, start_index)

	bans, err := db.uidNameQuery(rows, err)

	if err != nil {
		return nil, err
	}

	return bans, nil
}

// Search by name
func (db *appdbimpl) SearchByName(name string, requesting_uid string, start_index int, limit int) (*[]structures.SearchResult, error) {

	rows, err := db.c.Query(`SELECT "uid", "name",
							(
								SELECT EXISTS(
									SELECT * FROM "follows" AS "f"
									WHERE "f"."follower" = ?
									AND "f"."followed" = "users"."uid"
								)
							),
							(
								SELECT EXISTS(
									SELECT * FROM "bans" AS "b"
									WHERE "b"."user" = ?
									AND "b"."ban" = "users"."uid"
								)
							)
	
							FROM "users"
							WHERE "name" LIKE '%' || ? || '%'

							AND "uid" NOT IN (
								SELECT "bans"."user" FROM "bans"
								WHERE "bans"."user" = "users"."uid"
								AND "bans"."ban" = ?
							)
							LIMIT ?
							OFFSET ?`, requesting_uid, requesting_uid, name, requesting_uid, limit, start_index)

	if err != nil {
		return nil, err
	}

	var search_data []structures.SearchResult = make([]structures.SearchResult, 0)

	defer rows.Close()

	for rows.Next() {
		var search_entry structures.SearchResult
		err = rows.Scan(&search_entry.UID, &search_entry.Name, &search_entry.Followed, &search_entry.Banned)
		if err != nil {
			return nil, err
		}
		search_data = append(search_data, search_entry)
	}
	// We check if the iteration ended prematurely
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &search_data, nil
}
