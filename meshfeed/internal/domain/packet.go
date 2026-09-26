package domain

import (
	"encoding/binary"
	"errors"
)

// MaxPacketBytes is a conservative real-world LoRa payload budget. Actual
// limits range roughly 51-222 bytes depending on spreading factor and
// region (see BLUEPRINT.md); 200 leaves headroom across common
// configurations rather than assuming the best case. This single constant
// is the thing every other size limit in this file is derived from —
// change it here if you're tuning for a specific radio config, and every
// packet type stays consistent.
const MaxPacketBytes = 200

// header layout, common to every packet type:
//
//	[0]      type (1 byte)
//	[1:9]    feed ID (8 bytes)
//	[9]      TTL (1 byte) — hop budget, decremented by relays, dropped at 0
//	[10:14]  message ID (4 bytes) — random, used for relay dedup
const headerBytes = 1 + 8 + 1 + 4

type PacketType uint8

const (
	PacketAdvertise PacketType = iota // "I have this feed up to seq N"
	PacketWant                        // "send me this feed from seq N"
	PacketItem                        // one feed entry
)

// Packet is the union of everything that goes over the air. Only the
// fields relevant to Type are meaningful — same tagged-union approach as
// protobuf oneof, done by hand because a protobuf runtime is a lot of
// dependency weight for a 200-byte message (CLAUDE.md rule 6).
type Packet struct {
	Type   PacketType
	FeedID FeedID
	TTL    uint8
	MsgID  uint32

	LatestSeq uint32 // PacketAdvertise
	FromSeq   uint32 // PacketWant
	Item      Item   // PacketItem
}

var (
	ErrPacketTooShort = errors.New("packet: shorter than header")
	ErrPacketTooLong  = errors.New("packet: exceeds MaxPacketBytes")
	ErrUnknownType    = errors.New("packet: unknown type byte")
)

// MaxItemContentBytes is how much of an Item's Content actually fits once
// the item packet's own fields (seq + timestamp, 8 bytes) are accounted
// for on top of the shared header.
const MaxItemContentBytes = MaxPacketBytes - headerBytes - 4 - 4

func (p Packet) Marshal() ([]byte, error) {
	buf := make([]byte, headerBytes)
	buf[0] = byte(p.Type)
	copy(buf[1:9], p.FeedID[:])
	buf[9] = p.TTL
	binary.BigEndian.PutUint32(buf[10:14], p.MsgID)

	switch p.Type {
	case PacketAdvertise:
		tail := make([]byte, 4)
		binary.BigEndian.PutUint32(tail, p.LatestSeq)
		buf = append(buf, tail...)
	case PacketWant:
		tail := make([]byte, 4)
		binary.BigEndian.PutUint32(tail, p.FromSeq)
		buf = append(buf, tail...)
	case PacketItem:
		if len(p.Item.Content) > MaxItemContentBytes {
			return nil, ErrPacketTooLong
		}
		tail := make([]byte, 8+len(p.Item.Content))
		binary.BigEndian.PutUint32(tail[0:4], p.Item.Seq)
		binary.BigEndian.PutUint32(tail[4:8], p.Item.Timestamp)
		copy(tail[8:], p.Item.Content)
		buf = append(buf, tail...)
	default:
		return nil, ErrUnknownType
	}

	if len(buf) > MaxPacketBytes {
		return nil, ErrPacketTooLong
	}
	return buf, nil
}

func Unmarshal(data []byte) (Packet, error) {
	if len(data) < headerBytes {
		return Packet{}, ErrPacketTooShort
	}
	if len(data) > MaxPacketBytes {
		return Packet{}, ErrPacketTooLong
	}

	p := Packet{
		Type:  PacketType(data[0]),
		TTL:   data[9],
		MsgID: binary.BigEndian.Uint32(data[10:14]),
	}
	copy(p.FeedID[:], data[1:9])
	rest := data[headerBytes:]

	switch p.Type {
	case PacketAdvertise:
		if len(rest) != 4 {
			return Packet{}, ErrPacketTooShort
		}
		p.LatestSeq = binary.BigEndian.Uint32(rest)
	case PacketWant:
		if len(rest) != 4 {
			return Packet{}, ErrPacketTooShort
		}
		p.FromSeq = binary.BigEndian.Uint32(rest)
	case PacketItem:
		if len(rest) < 8 {
			return Packet{}, ErrPacketTooShort
		}
		p.Item.Seq = binary.BigEndian.Uint32(rest[0:4])
		p.Item.Timestamp = binary.BigEndian.Uint32(rest[4:8])
		content := make([]byte, len(rest)-8)
		copy(content, rest[8:])
		p.Item.Content = content
	default:
		return Packet{}, ErrUnknownType
	}

	return p, nil
}
