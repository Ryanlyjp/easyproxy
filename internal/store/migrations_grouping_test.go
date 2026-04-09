package store

import (
    "database/sql"
    "path/filepath"
    "testing"

    _ "modernc.org/sqlite"
)

func TestMigrationsCreateSubscriptionGroupingTables(t *testing.T) {
    dbPath := filepath.Join(t.TempDir(), "test.db")
    st, err := Open(dbPath)
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    if c, ok := st.(interface{ Close() error }); ok {
        defer c.Close()
    }

    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        t.Fatalf("open db: %v", err)
    }
    defer db.Close()

    for _, table := range []string{"subscription_groups", "subscription_entries"} {
        var count int
        err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = ? AND name = ?`, "table", table).Scan(&count)
        if err != nil {
            t.Fatalf("query table %s: %v", table, err)
        }
        if count != 1 {
            t.Fatalf("expected table %s to exist", table)
        }
    }
}
