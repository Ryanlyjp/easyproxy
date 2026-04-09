package monitor

import (
    "path/filepath"
    "testing"

    "easy_proxies/internal/config"
    "easy_proxies/internal/store"
)

func TestGetAllSettingsIncludesSubscriptionGroups(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    group, err := st.CreateSubscriptionGroup(t.Context(), store.SubscriptionGroup{Name: "G1", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(t.Context(), store.SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/sub", Enabled: true, SortOrder: 1}); err != nil {
        t.Fatalf("create entry: %v", err)
    }

    srv := &Server{cfgSrc: &config.Config{Subscriptions: []string{"https://legacy.example/sub"}}, store: st}
    resp := srv.getAllSettings()
    if len(resp.SubscriptionGroups) != 1 {
        t.Fatalf("expected 1 subscription group, got %d", len(resp.SubscriptionGroups))
    }
    if len(resp.SubscriptionGroups[0].Entries) != 1 {
        t.Fatalf("expected 1 entry, got %d", len(resp.SubscriptionGroups[0].Entries))
    }
}
