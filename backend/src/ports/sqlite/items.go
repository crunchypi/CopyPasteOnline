package sqlite

import (
	"copypaste-api/config"
	"database/sql"
	"errors"
	"time"
)

// How many seconds before entries are flushed.
var itemTimeoutSeconds int64 = config.ItemTimeoutSeconds

type Item struct {
	ID       int
	Mnemonic string
	Data     []byte
	Time     int64
}

func (s *SQLiteManager) CreateItemTable() error {
	return s.modifier(func() (string, []interface{}) {
		sql := `
			CREATE TABLE IF NOT EXISTS Items(
				"id" 		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"mnemonic"	TEXT NOT NULL,
				"data" 		BLOB NOT NULL,
				"time"		INTEGER NOT NULL
			);
		`
		return sql, make([]interface{}, 0)
	})
}

func (s *SQLiteManager) CreateItem(mnemonic string, data []byte) error {
	return s.modifier(func() (string, []interface{}) {
		sql := `
			INSERT INTO Items(mnemonic, data, time)
				 VALUES (?,?,?)
		`
		bindings := make([]interface{}, 3)
		bindings[0] = mnemonic
		bindings[1] = data
		bindings[2] = time.Now().Unix()
		return sql, bindings
	})
}

func (s *SQLiteManager) timeoutItems() error {
	threshold := time.Now().Unix() - itemTimeoutSeconds
	return s.modifier(func() (string, []interface{}) {
		sql := `
			DELETE FROM Items
			 WHERE time <= ?
		`
		bindings := make([]interface{}, 1)
		bindings[0] = threshold
		return sql, bindings
	})
}

func (s *SQLiteManager) readItemsWhere(sqlStr string, bindings []interface{}) ([]Item, error) {
	// # Clear old items first.
	s.timeoutItems()

	cmd := func() (string, []interface{}) {
		sql := "SELECT * FROM Items WHERE " + sqlStr
		return sql, bindings
	}
	items := make([]Item, 0, 10) // # 10 is arbitrary.
	callback := func(r *sql.Rows) {
		item := Item{}
		r.Scan(&item.ID, &item.Mnemonic, &item.Data, &item.Time)
		items = append(items, item)
	}

	return items, s.retriever(cmd, callback)
}

func (s *SQLiteManager) ReadItemByMnemonic(mnemonic string) (Item, bool, error) {
	sqlStr := "mnemonic = ?"
	bindings := make([]interface{}, 1)
	bindings[0] = mnemonic
	items, err := s.readItemsWhere(sqlStr, bindings)
	if err != nil {
		return Item{}, false, err
	}
	// # preventing security issue.
	if len(items) > 1 {
		return Item{}, true, errors.New("got multiple results, not allowed")
	}
	if len(items) == 0 {
		return Item{}, false, nil
	}
	return items[0], true, nil
}

func (s *SQLiteManager) ReadMnemonicExists(mnemonic string) (bool, error) {
	sqlStr := "mnemonic = ?"
	bindings := make([]interface{}, 1)
	bindings[0] = mnemonic
	items, err := s.readItemsWhere(sqlStr, bindings)
	return len(items) > 0, err
}
