package domain_test

import (
	"bytes"
	"testing"

	"meshfeed/internal/domain"
)

func TestPacketRoundTrip(t *testing.T) {
	feedID := domain.NewFeedID("garden-club-updates")

	cases := []domain.Packet{
		{Type: domain.PacketAdvertise, FeedID: feedID, TTL: 3, MsgID: 42, LatestSeq: 7},
		{Type: domain.PacketWant, FeedID: feedID, TTL: 3, MsgID: 43, FromSeq: 5},
		{Type: domain.PacketItem, FeedID: feedID, TTL: 3, MsgID: 44, Item: domain.Item{
			Seq: 6, Timestamp: 1_700_000_000, Content: []byte("rain barrel is full, come get some"),
		}},
	}

	for _, want := range cases {
		data, err := want.Marshal()
		if err != nil {
			t.Fatalf("Marshal(%v): %v", want.Type, err)
		}
		if len(data) > domain.MaxPacketBytes {
			t.Fatalf("packet exceeds MaxPacketBytes: %d", len(data))
		}
		got, err := domain.Unmarshal(data)
		if err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if got.Type != want.Type || got.FeedID != want.FeedID || got.TTL != want.TTL || got.MsgID != want.MsgID {
			t.Fatalf("header mismatch: got %+v, want %+v", got, want)
		}
		switch want.Type {
		case domain.PacketAdvertise:
			if got.LatestSeq != want.LatestSeq {
				t.Fatalf("LatestSeq: got %d, want %d", got.LatestSeq, want.LatestSeq)
			}
		case domain.PacketWant:
			if got.FromSeq != want.FromSeq {
				t.Fatalf("FromSeq: got %d, want %d", got.FromSeq, want.FromSeq)
			}
		case domain.PacketItem:
			if got.Item.Seq != want.Item.Seq || got.Item.Timestamp != want.Item.Timestamp ||
				!bytes.Equal(got.Item.Content, want.Item.Content) {
				t.Fatalf("Item: got %+v, want %+v", got.Item, want.Item)
			}
		}
	}
}

func TestItem_ContentOverBudgetRejected(t *testing.T) {
	p := domain.Packet{
		Type:   domain.PacketItem,
		FeedID: domain.NewFeedID("f"),
		Item:   domain.Item{Content: bytes.Repeat([]byte("x"), domain.MaxItemContentBytes+1)},
	}
	if _, err := p.Marshal(); err != domain.ErrPacketTooLong {
		t.Fatalf("expected ErrPacketTooLong, got %v", err)
	}
}

func TestItem_ContentAtBudgetFits(t *testing.T) {
	p := domain.Packet{
		Type:   domain.PacketItem,
		FeedID: domain.NewFeedID("f"),
		Item:   domain.Item{Content: bytes.Repeat([]byte("x"), domain.MaxItemContentBytes)},
	}
	data, err := p.Marshal()
	if err != nil {
		t.Fatalf("Marshal at exact budget: %v", err)
	}
	if len(data) != domain.MaxPacketBytes {
		t.Fatalf("expected exactly MaxPacketBytes (%d), got %d", domain.MaxPacketBytes, len(data))
	}
}

func TestUnmarshal_RejectsShortData(t *testing.T) {
	if _, err := domain.Unmarshal([]byte{1, 2, 3}); err != domain.ErrPacketTooShort {
		t.Fatalf("expected ErrPacketTooShort, got %v", err)
	}
}
