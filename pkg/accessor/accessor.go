package accessor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/altairsix/eventsource/mysqlstore"

	// to register mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ErrTypeCast if returned when the expected type does not match
var ErrTypeCast = errors.New("type cast failed")

// New return a MySQL accessor and creates a table if not exists
func New(driver, dsn, tableName string) (mysqlstore.Accessor, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if err := mysqlstore.CreateIfNotExists(db, tableName); err != nil {
		return nil, err
	}

	return &accessor{
		driver: driver,
		dsn:    dsn,
	}, nil
}

type accessor struct {
	driver string
	dsn    string
}

// Open creates a database object and tests the connection
func (a *accessor) Open(ctx context.Context) (mysqlstore.DB, error) {

	db, err := sql.Open(a.driver, a.dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Close the database
func (a *accessor) Close(db mysqlstore.DB) error {
	d, ok := db.(*sql.DB)
	if !ok {
		return ErrTypeCast
	}
	return d.Close()
}
