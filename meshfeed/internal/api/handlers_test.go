package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"meshfeed/internal/api"
	"meshfeed/internal/domain"
	"meshfeed/internal/store"
)

func TestListFeedsAndItems(t *testing.T) {
	s := store.New()
	feed := domain.NewFeedID("garden-club")
	s.Append(feed, domain.Item{Seq: 1, Timestamp: 100, Content: []byte("hello")})
	s.Append(feed, domain.Item{Seq: 2, Timestamp: 200, Content: []byte("world")})

	srv := httptest.NewServer(api.NewHandler(s).Routes())
	defer srv.Close()

	feedsResp, err := http.Get(srv.URL + "/feeds")
	if err != nil {
		t.Fatalf("GET /feeds: %v", err)
	}
	var feeds []struct {
		ID        string `json:"id"`
		LatestSeq uint32 `json:"latest_seq"`
	}
	if err := json.NewDecoder(feedsResp.Body).Decode(&feeds); err != nil {
		t.Fatalf("decode /feeds: %v", err)
	}
	if len(feeds) != 1 || feeds[0].ID != feed.String() || feeds[0].LatestSeq != 2 {
		t.Fatalf("unexpected /feeds response: %+v", feeds)
	}

	itemsResp, err := http.Get(srv.URL + "/feeds/" + feed.String() + "/items")
	if err != nil {
		t.Fatalf("GET items: %v", err)
	}
	var items []struct {
		Seq     uint32 `json:"seq"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(itemsResp.Body).Decode(&items); err != nil {
		t.Fatalf("decode items: %v", err)
	}
	if len(items) != 2 || items[0].Content != "hello" || items[1].Content != "world" {
		t.Fatalf("unexpected items response: %+v", items)
	}
}

func TestListItems_RejectsInvalidFeedID(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler(store.New()).Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/feeds/not-hex/items")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid feed id, got %d", resp.StatusCode)
	}
}
