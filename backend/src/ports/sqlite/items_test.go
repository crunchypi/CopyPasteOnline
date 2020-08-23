package sqlite

import (
	"os"
	"testing"
	"time"
)

// # tasker is a util func which creates a db, table and executes code before cleanup.
func tasker(t *testing.T, task func(*SQLiteManager)) {
	path := "item.sqlite"
	db, err := New(path)
	if err != nil {
		t.Error("DB creation failure:", err)
		return
	}
	defer os.Remove(path)
	if err := db.CreateItemTable(); err != nil {
		t.Error(err)
	}
	task(db)
}

// func TestCreateItemTable(t *testing.T) { // # has to be manually checked.
// 	db, err := New("./items.sqlite")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = db.CreateItemTable()
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestCreateReadItem(t *testing.T) {
	tasker(t, func(db *SQLiteManager) {
		mnemonic, data := "testMnemonic", "testData"
		if err := db.CreateItem(mnemonic, []byte(data)); err != nil {
			t.Error("failed while creating item:", err)
		}

		res, ok, err := db.ReadItemByMnemonic(mnemonic)
		if err != nil {
			t.Error("failed while reading data:", err)
		}
		if !ok {
			t.Error("got no result while reading by mnemonic")
		}
		if res.Mnemonic != mnemonic || string(res.Data) != data {
			t.Error("recieved incorrect data:", res)
		}
	})
}

func TestTimeoutItems(t *testing.T) {
	tasker(t, func(db *SQLiteManager) {
		timeoutBackup := itemTimeoutSeconds // # Allow temporary modification of this var
		{
			itemTimeoutSeconds = 1
			mnemonic, data := "testmnemonic", "testdata"
			db.CreateItem(mnemonic, []byte(data))
			time.Sleep(time.Second * 2)
			_, ok, _ := db.ReadItemByMnemonic(mnemonic)
			if ok {
				t.Error("item not removed")
			}
		}
		{
			itemTimeoutSeconds = 3
			mnemonic, data := "testmnemonic", "testdata"
			db.CreateItem(mnemonic, []byte(data))
			time.Sleep(time.Second * 2)
			_, ok, _ := db.ReadItemByMnemonic(mnemonic)
			if !ok {
				t.Error("item removed too soon")
			}
		}
		itemTimeoutSeconds = timeoutBackup

	})
}
