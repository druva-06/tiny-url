package db

import (
	"database/sql"
	"sync/atomic"
)

type DBCluster struct {
	Primary  *sql.DB
	Replicas []*sql.DB
	counter  uint64
}

func (c *DBCluster) getReplica() *sql.DB {
	n := atomic.AddUint64(&c.counter, 1)
	return c.Replicas[n%uint64(len(c.Replicas))]
}

func (c *DBCluster) Exec(query string, args ...any) (sql.Result, error) {
	return c.Primary.Exec(query, args...)
}

func (c *DBCluster) Query(query string, args ...any) (*sql.Rows, error) {
	replica := c.getReplica()

	rows, err := replica.Query(query, args...)
	if err != nil {
		// fallback to primary
		return c.Primary.Query(query, args...)
	}

	return rows, nil
}

func (c *DBCluster) QueryWithConsistency(usePrimary bool, query string, args ...any) (*sql.Rows, error) {
	if usePrimary {
		return c.Primary.Query(query, args...)
	}
	return c.Query(query, args...)
}
