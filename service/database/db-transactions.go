package database

import "database/sql"

// dbtransaction is a struct to represent an SQL transaction, it implements the DBTransaction interface
type dbtransaction struct {
	c *sql.Tx
}

func (tx *dbtransaction) Commit() error {
	// Commit the SQL transaction
	return tx.c.Commit()
}

func (tx *dbtransaction) Rollback() error {
	// Rollback the SQL transaction
	return tx.c.Rollback()
}
