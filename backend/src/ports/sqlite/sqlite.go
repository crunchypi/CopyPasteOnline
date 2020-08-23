package sqlite

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteManager manages an sql connection
type SQLiteManager struct {
	db    *sql.DB
	mutex sync.Mutex
}

// New creates a new sql connection and a template table
// (if one does not already exist).
func New(path string) (*SQLiteManager, error) {
	man := SQLiteManager{}
	conn, err := sql.Open("sqlite3", path)
	man.db = conn
	man.integrity()
	return &man, err
}

// modifier is a wrapper for all functions which modify the db.
// argument should be a function which return an sql string and
// any bindings (as slice of interfaces)
func (s *SQLiteManager) modifier(query func() (string, []interface{})) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sql, bindings := query()
	statement, err := s.db.Prepare(sql)
	if err != nil {
		return err
	}
	if _, err := statement.Exec(bindings...); err != nil {
		return err
	}
	if err := statement.Close(); err != nil {
		return err
	}
	return nil
}

// retriever is a wrapper for all functions which retrieve data
// from the db. Args:
// (1) query: function which returns an sql string and bindings
// (2) callback: called on each read, must take a ref to sql.Rows.
//   		   this is used to pull data from each row.
func (s *SQLiteManager) retriever(query func() (string, []interface{}),
	callback func(*sql.Rows)) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	sql, bindings := query()
	row, err := s.db.Query(sql, bindings...)
	if err != nil {
		return err
	}

	for row.Next() {
		callback(row)
	}
	return nil
}

// Enable foreign key constraint.
func (s *SQLiteManager) integrity() error {
	return s.modifier(func() (string, []interface{}) {
		return "PRAGMA foreign_keys = ON", make([]interface{}, 0)
	})
}
