package monitor

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "path/filepath"
	"strconv"
    "testing"

    "easy_proxies/internal/store"
)

func TestHandleSubscriptionGroupsCreate(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    srv := &Server{store: st}
    body, _ := json.Marshal(map[string]any{"name": "G1"})
    req := httptest.NewRequest(http.MethodPost, "/api/subscription/groups", bytes.NewReader(body))
    w := httptest.NewRecorder()

    srv.handleSubscriptionGroups(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    groups, err := st.ListSubscriptionGroups(req.Context())
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 || groups[0].Name != "G1" {
        t.Fatalf("unexpected groups after create: %+v", groups)
    }
}

func TestHandleSubscriptionEntryItemToggle(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    group, err := st.CreateSubscriptionGroup(t.Context(), store.SubscriptionGroup{Name: "G1", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    entry, err := st.CreateSubscriptionEntry(t.Context(), store.SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/sub", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create entry: %v", err)
    }

    srv := &Server{store: st}
    body, _ := json.Marshal(map[string]any{"enabled": false})
    req := httptest.NewRequest(http.MethodPatch, "/api/subscription/entries/"+strconv.FormatInt(entry.ID, 10), bytes.NewReader(body))
    w := httptest.NewRecorder()

    srv.handleSubscriptionEntryItem(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    urls, err := st.ListEnabledSubscriptionEntryURLs(req.Context())
    if err != nil {
        t.Fatalf("list enabled urls: %v", err)
    }
    if len(urls) != 0 {
        t.Fatalf("expected no enabled urls after toggle, got %d", len(urls))
    }
}

func TestHandleSubscriptionGroupItemToggle(t *testing.T) {
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

    srv := &Server{store: st}
    body, _ := json.Marshal(map[string]any{"enabled": false})
    req := httptest.NewRequest(http.MethodPatch, "/api/subscription/groups/"+strconv.FormatInt(group.ID, 10), bytes.NewReader(body))
    w := httptest.NewRecorder()

    srv.handleSubscriptionGroupItem(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    urls, err := st.ListEnabledSubscriptionEntryURLs(req.Context())
    if err != nil {
        t.Fatalf("list enabled urls: %v", err)
    }
    if len(urls) != 0 {
        t.Fatalf("expected no enabled urls after group disable, got %d", len(urls))
    }
}


func TestHandleSubscriptionEntryItemDelete(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    group, err := st.CreateSubscriptionGroup(t.Context(), store.SubscriptionGroup{Name: "G1", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    entry, err := st.CreateSubscriptionEntry(t.Context(), store.SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/sub", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create entry: %v", err)
    }

    srv := &Server{store: st}
    req := httptest.NewRequest(http.MethodDelete, "/api/subscription/entries/"+strconv.FormatInt(entry.ID, 10), nil)
    w := httptest.NewRecorder()

    srv.handleSubscriptionEntryItem(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    groups, err := st.ListSubscriptionGroups(req.Context())
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 {
        t.Fatalf("expected 1 group, got %d", len(groups))
    }
    if len(groups[0].Entries) != 0 {
        t.Fatalf("expected 0 entries after delete, got %d", len(groups[0].Entries))
    }
}


func TestHandleSubscriptionQuickAddCreatesDefaultGroup(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    srv := &Server{store: st}
    body, _ := json.Marshal(map[string]any{"url": "https://example.com/sub-a"})
    req := httptest.NewRequest(http.MethodPost, "/api/subscription/quick-add", bytes.NewReader(body))
    w := httptest.NewRecorder()

    srv.handleSubscriptionQuickAdd(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    groups, err := st.ListSubscriptionGroups(req.Context())
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 || groups[0].Name != "Default" {
        t.Fatalf("unexpected groups after quick add: %+v", groups)
    }
    if len(groups[0].Entries) != 1 {
        t.Fatalf("expected 1 entry in Default, got %d", len(groups[0].Entries))
    }
    if groups[0].Entries[0].Alias != "SUB-1" {
        t.Fatalf("expected alias SUB-1, got %s", groups[0].Entries[0].Alias)
    }
}

func TestQuickAddSubscriptionAliasReusesSmallestHole(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    srv := &Server{store: st}
    for _, url := range []string{"https://example.com/sub-a", "https://example.com/sub-b"} {
        body, _ := json.Marshal(map[string]any{"url": url})
        req := httptest.NewRequest(http.MethodPost, "/api/subscription/quick-add", bytes.NewReader(body))
        w := httptest.NewRecorder()
        srv.handleSubscriptionQuickAdd(w, req)
        if w.Code != http.StatusOK {
            t.Fatalf("quick add failed: %d %s", w.Code, w.Body.String())
        }
    }

    groups, err := st.ListSubscriptionGroups(t.Context())
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 1 || len(groups[0].Entries) != 2 {
        t.Fatalf("unexpected groups after initial quick add: %+v", groups)
    }

    firstID := groups[0].Entries[0].ID
    if err := st.DeleteSubscriptionEntry(t.Context(), firstID); err != nil {
        t.Fatalf("delete first entry: %v", err)
    }

    body, _ := json.Marshal(map[string]any{"url": "https://example.com/sub-c"})
    req := httptest.NewRequest(http.MethodPost, "/api/subscription/quick-add", bytes.NewReader(body))
    w := httptest.NewRecorder()
    srv.handleSubscriptionQuickAdd(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    groups, err = st.ListSubscriptionGroups(t.Context())
    if err != nil {
        t.Fatalf("list groups after reuse: %v", err)
    }
    aliases := map[string]bool{}
    for _, entry := range groups[0].Entries {
        aliases[entry.Alias] = true
    }
    if !aliases["SUB-1"] || !aliases["SUB-2"] {
        t.Fatalf("expected alias hole reuse to yield SUB-1 and SUB-2, got %+v", aliases)
    }
}

func TestHandleSubscriptionGroupItemDeleteCascadesEntries(t *testing.T) {
    st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer st.Close()

    group, err := st.CreateSubscriptionGroup(t.Context(), store.SubscriptionGroup{Name: "G1", Enabled: true, SortOrder: 1})
    if err != nil {
        t.Fatalf("create group: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(t.Context(), store.SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/a", Alias: "SUB-1", Enabled: true, SortOrder: 1}); err != nil {
        t.Fatalf("create first entry: %v", err)
    }
    if _, err := st.CreateSubscriptionEntry(t.Context(), store.SubscriptionEntry{GroupID: group.ID, URL: "https://example.com/b", Alias: "SUB-2", Enabled: true, SortOrder: 2}); err != nil {
        t.Fatalf("create second entry: %v", err)
    }

    srv := &Server{store: st}
    req := httptest.NewRequest(http.MethodDelete, "/api/subscription/groups/"+strconv.FormatInt(group.ID, 10), nil)
    w := httptest.NewRecorder()

    srv.handleSubscriptionGroupItem(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    groups, err := st.ListSubscriptionGroups(t.Context())
    if err != nil {
        t.Fatalf("list groups: %v", err)
    }
    if len(groups) != 0 {
        t.Fatalf("expected 0 groups after delete, got %d", len(groups))
    }
}
