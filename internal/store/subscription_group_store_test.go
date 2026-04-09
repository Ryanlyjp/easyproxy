package store

import (
    "context"
    "path/filepath"
    "testing"
)

func TestSubscriptionGroupCRUD(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()

    group, err := st.CreateSubscriptionGroup(ctx, SubscriptionGroup{
        Name:      "Default",
        Enabled:   true,
        SortOrder: 1,
    })
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    if group.ID == 0 {
        t.Fatalf("expected created group id")
    }

    entry, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{
        GroupID:   group.ID,
        URL:       "https://example.com/sub1",
        Alias:     "sub1",
        Enabled:   true,
        SortOrder: 1,
    })
    if err != nil {
        t.Fatalf("create entry: %v", err)
    }
    if entry.ID == 0 {
        t.Fatalf("expected created entry id")
    }

    groups, err := st.ListSubscriptionGroups(ctx)
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 {
        t.Fatalf("expected 1 group, got %d", len(groups))
    }
    if len(groups[0].Entries) != 1 {
        t.Fatalf("expected 1 entry in group, got %d", len(groups[0].Entries))
    }
    if groups[0].Entries[0].URL != "https://example.com/sub1" {
        t.Fatalf("unexpected entry url: %s", groups[0].Entries[0].URL)
    }
}

func TestSubscriptionGroupToggleEntryEnabledAffectsEnabledURLList(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()
    group, err := st.CreateSubscriptionGroup(ctx, SubscriptionGroup{Name: "Default", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    entry, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/sub1", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create entry: %v", err)
    }

    if err := st.UpdateSubscriptionEntryEnabled(ctx, entry.ID, false); err != nil {
        t.Fatalf("disable entry: %v", err)
    }

    urls, err := st.ListEnabledSubscriptionEntryURLs(ctx)
    if err != nil {
        t.Fatalf("list enabled urls: %v", err)
    }
    if len(urls) != 0 {
        t.Fatalf("expected 0 enabled urls after disable, got %d", len(urls))
    }
}


func TestDeleteSubscriptionEntry(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()
    group, err := st.CreateSubscriptionGroup(ctx, SubscriptionGroup{Name: "Default", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    first, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/a", Alias: "A", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create first entry: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/b", Alias: "B", Enabled: true, SortOrder: 2}); err != nil {
        t.Fatalf("create second entry: %v", err)
    }

    if err := st.DeleteSubscriptionEntry(ctx, first.ID); err != nil {
        t.Fatalf("delete entry: %v", err)
    }

    groups, err := st.ListSubscriptionGroups(ctx)
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 {
        t.Fatalf("expected 1 group, got %d", len(groups))
    }
    if len(groups[0].Entries) != 1 {
        t.Fatalf("expected 1 remaining entry, got %d", len(groups[0].Entries))
    }
    if groups[0].Entries[0].URL != "https://example.com/b" {
        t.Fatalf("unexpected remaining entry: %+v", groups[0].Entries[0])
    }
}


func TestFindSubscriptionGroupByName(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()
    if _, err := st.CreateSubscriptionGroup(ctx, SubscriptionGroup{Name: "Default", Enabled: true, SortOrder: 1}); err != nil {
        t.Fatalf("create group: %v", err)
    }

    group, err := st.FindSubscriptionGroupByName(ctx, "Default")
    if err != nil {
        t.Fatalf("find group: %v", err)
    }
    if group == nil {
        t.Fatalf("expected group, got nil")
    }
    if group.Name != "Default" {
        t.Fatalf("unexpected group name: %s", group.Name)
    }

    missing, err := st.FindSubscriptionGroupByName(ctx, "Missing")
    if err != nil {
        t.Fatalf("find missing group: %v", err)
    }
    if missing != nil {
        t.Fatalf("expected nil for missing group, got %+v", missing)
    }
}

func TestListSubscriptionEntryAliases(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()
    group, err := st.CreateSubscriptionGroup(ctx, SubscriptionGroup{Name: "Default", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/a", Alias: "SUB-1", Enabled: true, SortOrder: 1}); err != nil {
        t.Fatalf("create first entry: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/b", Alias: "SUB-3", Enabled: true, SortOrder: 2}); err != nil {
        t.Fatalf("create second entry: %v", err)
    }

    aliases, err := st.ListSubscriptionEntryAliases(ctx)
    if err != nil {
        t.Fatalf("list aliases: %v", err)
    }
    if len(aliases) != 2 {
        t.Fatalf("expected 2 aliases, got %d", len(aliases))
    }
    if aliases[0] != "SUB-1" || aliases[1] != "SUB-3" {
        t.Fatalf("unexpected aliases: %+v", aliases)
    }
}

func TestDeleteSubscriptionGroupCascadesEntries(t *testing.T) {
    st, err := Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    ctx := context.Background()
    group, err := st.CreateSubscriptionGroup(ctx, SubscriptionGroup{Name: "Default", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/a", Alias: "SUB-1", Enabled: true, SortOrder: 1}); err != nil {
        t.Fatalf("create first entry: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(ctx, SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/b", Alias: "SUB-2", Enabled: true, SortOrder: 2}); err != nil {
        t.Fatalf("create second entry: %v", err)
    }

    if err := st.DeleteSubscriptionGroup(ctx, group.ID); err != nil {
        t.Fatalf("delete group: %v", err)
    }

    groups, err := st.ListSubscriptionGroups(ctx)
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 0 {
        t.Fatalf("expected 0 groups, got %d", len(groups))
    }

    aliases, err := st.ListSubscriptionEntryAliases(ctx)
    if err != nil {
        t.Fatalf("list aliases after delete: %v", err)
    }
    if len(aliases) != 0 {
        t.Fatalf("expected 0 aliases after cascade delete, got %d", len(aliases))
    }
}
