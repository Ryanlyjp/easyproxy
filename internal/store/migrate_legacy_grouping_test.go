package store

import (
    "context"
    "path/filepath"
    "testing"
)

func TestMigrateLegacyDataImportsSubscriptionsIntoDefaultGroup(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    err = MigrateLegacyData(context.Background(), st, t.TempDir(), nil, "", []string{
        "https://example.com/sub-a",
        "https://example.com/sub-b",
    })
    if err != nil {
        t.Fatalf("migrate legacy data: %v", err)
    }

    groups, err := st.ListSubscriptionGroups(context.Background())
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 {
        t.Fatalf("expected 1 group after migration, got %d", len(groups))
    }
    if groups[0].Name != "Default" {
        t.Fatalf("expected default group name, got %q", groups[0].Name)
    }
    if len(groups[0].Entries) != 2 {
        t.Fatalf("expected 2 migrated subscriptions, got %d", len(groups[0].Entries))
    }
}
