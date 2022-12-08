package database

import "database/sql"

type dbtransaction struct {
	c *sql.Tx
}

func (tx *dbtransaction) Commit() error {
	return tx.c.Commit()
}

func (tx *dbtransaction) Rollback() error {
	return tx.c.Rollback()
}
