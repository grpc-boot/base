package sqlite

import "testing"

func TestNewDb(t *testing.T) {
	opt := DefaultOption()
	opt.DbName = "test.db"

	db, _ := NewDb(opt)
	_, err := db.Exec(`DROP TABLE user; CREATE TABLE user(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL, 
    passwd CHAR(32) NOT NULL DEFAULT '',
    created_at INT DEFAULT 0
)`)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	tableSql, err := db.ShowCreateTable("user")
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("tables: %v", tableSql)
}
