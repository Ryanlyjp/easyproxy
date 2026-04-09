package subscription

import (
    "context"
    "path/filepath"
    "testing"

    "easy_proxies/internal/config"
    "easy_proxies/internal/store"
)

func TestFetchSubscriptionURLsUsesEnabledGroupedEntries(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()
    group, err := st.CreateSubscriptionGroup(ctx, store.SubscriptionGroup{Name: "G1", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, store.SubscriptionEntry{GroupID: group.ID, URL: "https://enabled.example/sub", Enabled: true, SortOrder: 1}); err != nil {
        t.Fatalf("create enabled entry: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, store.SubscriptionEntry{GroupID: group.ID, URL: "https://disabled.example/sub", Enabled: false, SortOrder: 2}); err != nil {
        t.Fatalf("create disabled entry: %v", err)
    }

    m := New(&config.Config{Subscriptions: []string{"https://legacy.example/sub"}}, nil, WithStore(st))
    defer m.Stop()

    urls := m.subscriptionURLs()
    if len(urls) != 1 {
        t.Fatalf("expected 1 enabled grouped subscription url, got %d", len(urls))
    }
    if urls[0] != "https://enabled.example/sub" {
        t.Fatalf("unexpected url: %s", urls[0])
    }
}

func TestFetchSubscriptionURLsFallsBackToLegacyConfigWhenGroupedEntriesMissing(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    m := New(&config.Config{Subscriptions: []string{"https://legacy.example/sub"}}, nil, WithStore(st))
    defer m.Stop()

    urls := m.subscriptionURLs()
    if len(urls) != 1 {
        t.Fatalf("expected 1 legacy subscription url, got %d", len(urls))
    }
    if urls[0] != "https://legacy.example/sub" {
        t.Fatalf("unexpected fallback url: %s", urls[0])
    }
}
