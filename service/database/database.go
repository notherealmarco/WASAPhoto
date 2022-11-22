/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/notherealmarco/WASAPhoto/service/structures"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(name string) (string, error)
	UserExists(uid string) (bool, error)
	GetUserID(name string) (string, error)

	SearchByName(name string, requesting_uid string, start_index int, limit int) (*[]structures.UIDName, error)

	UpdateUsername(uid, name string) error

	GetUserFollowers(uid string, requesting_uid string, start_index int, limit int) (QueryResult, *[]structures.UIDName, error)
	GetUserFollowing(uid string, requesting_uid string, start_index int, offset int) (QueryResult, *[]structures.UIDName, error)
	FollowUser(uid string, follow string) (QueryResult, error)
	UnfollowUser(uid string, unfollow string) (QueryResult, error)

	BanUser(uid string, ban string) (QueryResult, error)
	UnbanUser(uid string, unban string) (QueryResult, error)
	IsBanned(uid string, banner string) (bool, error)
	GetUserBans(uid string, start_index int, limit int) (*[]structures.UIDName, error)

	PostPhoto(uid string) (DBTransaction, int64, error)
	DeletePhoto(uid string, photo int64) (bool, error)

	GetPhotoLikes(uid string, photo int64, requesting_uid string, start_index int, offset int) (QueryResult, *[]structures.UIDName, error)
	LikePhoto(uid string, photo int64, liker_uid string) (QueryResult, error)
	UnlikePhoto(uid string, photo int64, liker_uid string) (QueryResult, error)

	GetUserProfile(uid string, requesting_uid string) (QueryResult, *structures.UserProfile, error)
	GetUserPhotos(uid string, start_index int, limit int) (*[]structures.UserPhoto, error)
	GetUserStream(uid string, start_index int, limit int) (*[]structures.Photo, error)

	GetComments(uid string, photo_id int64, requesting_uid string, start_index int, offset int) (QueryResult, *[]structures.Comment, error)
	PostComment(uid string, photo_id int64, comment_user string, comment string) (QueryResult, error)
	DeleteComment(uid string, photo_id int64, comment_id int64) (QueryResult, error)
	GetCommentOwner(uid string, photo_id int64, comment_id int64) (QueryResult, string, error)

	Ping() error
}

type DBTransaction interface {
	Commit() error
	Rollback() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if tables exist. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName) //todo: check for all the tables
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE "users" (
			"uid"	TEXT NOT NULL,
			"name"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("uid")
		)` //todo: one query is enough! Why do I need to do this?
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `CREATE TABLE "follows" (
			"follower"	TEXT NOT NULL,
			"followed"	TEXT NOT NULL,
			FOREIGN KEY("follower") REFERENCES "users"("uid") ON UPDATE CASCADE,
			FOREIGN KEY("followed") REFERENCES "users"("uid") ON UPDATE CASCADE,
			PRIMARY KEY("follower","followed")
		)`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `CREATE TABLE "bans" (
			"user"	TEXT NOT NULL,
			"ban"	TEXT NOT NULL,
			FOREIGN KEY("user") REFERENCES "users"("uid") ON UPDATE CASCADE,
			FOREIGN KEY("ban") REFERENCES "users"("uid") ON UPDATE CASCADE,
			PRIMARY KEY("user","ban")
		)`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `CREATE TABLE "photos" (
			"user"	TEXT NOT NULL,
			"id"	INTEGER NOT NULL,
			"date"	TEXT NOT NULL,
			FOREIGN KEY("user") REFERENCES "users"("uid") ON UPDATE CASCADE,
			PRIMARY KEY("id" AUTOINCREMENT)
		)`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `CREATE TABLE "comments" (
			"user"	TEXT NOT NULL,
			"photo"	INTEGER NOT NULL,
			"id"	INTEGER NOT NULL,
			"comment"	TEXT NOT NULL,
			"date"	TEXT NOT NULL,
			FOREIGN KEY("user") REFERENCES "users"("uid"),
			PRIMARY KEY("id" AUTOINCREMENT),
			FOREIGN KEY("photo") REFERENCES "photos"("id")
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `CREATE TABLE "likes" (
			"user"	TEXT NOT NULL,
			"photo_id"	INTEGER NOT NULL,
			FOREIGN KEY("user") REFERENCES "users"("uid") ON UPDATE CASCADE,
			FOREIGN KEY("photo_id") REFERENCES "photos"("id") ON UPDATE CASCADE,
			PRIMARY KEY("user","photo_id")
		)`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		sqlStmt = `PRAGMA foreign_keys = ON`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
