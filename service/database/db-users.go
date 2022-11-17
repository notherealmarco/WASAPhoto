package database

// Create a new user
func (db *appdbimpl) CreateUser(uid string, name string) error {
	_, err := db.c.Exec(`INSERT INTO "users" ("uid", "name") VALUES (?, ?)`, uid, name)
	return err
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
