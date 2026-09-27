// Package api is a small read-only HTTP view onto a node's local feed
// store — the "RSS reader" half of the metaphor. It's deliberately
// read-only: subscribing to a feed is a node-configuration decision (who
// do I want to hear from), not a web request, so there's no POST here.
package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"meshfeed/internal/domain"
	"meshfeed/internal/store"
)

var errInvalidFeedID = errors.New("invalid feed id")

type Handler struct {
	store *store.FeedStore
}

func NewHandler(s *store.FeedStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /feeds", h.listFeeds)
	mux.HandleFunc("GET /feeds/{id}/items", h.listItems)
	return mux
}

type feedView struct {
	ID        string `json:"id"`
	LatestSeq uint32 `json:"latest_seq"`
}

func (h *Handler) listFeeds(w http.ResponseWriter, _ *http.Request) {
	feeds := h.store.Feeds()
	views := make([]feedView, 0, len(feeds))
	for _, id := range feeds {
		views = append(views, feedView{ID: id.String(), LatestSeq: h.store.LatestSeq(id)})
	}
	writeJSON(w, views)
}

type itemView struct {
	Seq       uint32 `json:"seq"`
	Timestamp uint32 `json:"timestamp"`
	Content   string `json:"content"`
}

func (h *Handler) listItems(w http.ResponseWriter, r *http.Request) {
	idHex := r.PathValue("id")
	feedID, err := parseFeedID(idHex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items := h.store.ItemsSince(feedID, 0)
	views := make([]itemView, 0, len(items))
	for _, item := range items {
		views = append(views, itemView{Seq: item.Seq, Timestamp: item.Timestamp, Content: string(item.Content)})
	}
	writeJSON(w, views)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func parseFeedID(hexStr string) (domain.FeedID, error) {
	var id domain.FeedID
	decoded, err := hex.DecodeString(hexStr)
	if err != nil || len(decoded) != len(id) {
		return id, errInvalidFeedID
	}
	copy(id[:], decoded)
	return id, nil
}
