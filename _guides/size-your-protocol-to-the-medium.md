---
title: "Size your protocol to the medium, not the other way around"
summary: "Before designing a wire format, find the hard physical limit of what you're transmitting over — then let that number drive every other decision."
difficulty: advanced
related_projects:
  - meshfeed
---

## The mistake this guide helps you avoid

It's tempting to design a message format the way you'd design an API
response — pick a friendly serialization (JSON is the default reflex),
add the fields you need, ship it. That works fine over HTTP, where a
payload can be kilobytes without anyone noticing.

It fails completely on a constrained medium. LoRa radio, for example, has
a real payload ceiling of roughly 51–222 bytes depending on spreading
factor and region. JSON's per-field key overhead (`{"borrower_id":"alice"}`
spends 15 bytes on the key alone) can blow that budget before you've sent
any actual data. USPD's [common/bulk principle]({{ '/principles/' | relative_url }})
says default to boring formats — but "boring" has to be relative to your
actual constraint, and on a 200-byte budget, JSON is the exotic choice,
not the safe one.

## The method

1. **Find the hard number first.** Not an estimate — a cited one, from the
   medium's actual specification. For LoRa this is the regional parameters
   documentation (e.g. The Things Network's EU868 page). Write it down as
   a named constant before you design anything else.

2. **Budget the header before the payload.** Decide what every message
   needs regardless of type — a type tag, an identifier, a hop count,
   a dedup ID — and size that first. It's fixed overhead you pay on every
   message, so it should be as small as correctness allows.

3. **Derive every other limit from the first number.** Don't pick a
   separate "max content length" by feel — subtract the header size (and
   any per-type fields) from the hard ceiling. If the ceiling ever changes
   (a different spreading factor, a different radio entirely), one
   constant changes and everything downstream stays consistent.

4. **Write a test that asserts the boundary, not just the happy path.**
   A test that proves content one byte over the limit is rejected, and
   content exactly at the limit fits, catches the actual failure mode —
   silent truncation or a dropped packet in the field — before it ships.

## Worked example

[`meshfeed`]({{ '/projects/meshfeed/' | relative_url }}) does all four
steps in `internal/domain/packet.go`:

```go
// The hard number, cited, as the one thing everything else derives from:
const MaxPacketBytes = 200

// The header every message pays regardless of type:
//   type(1) + feedID(8) + TTL(1) + msgID(4) = 14 bytes
const headerBytes = 1 + 8 + 1 + 4

// Derived, not independently chosen:
const MaxItemContentBytes = MaxPacketBytes - headerBytes - 4 - 4
```

`internal/domain/packet_test.go` asserts the boundary directly:
content at exactly `MaxItemContentBytes` fits in exactly `MaxPacketBytes`;
content one byte over is rejected with `ErrPacketTooLong`. Both are real,
executable tests — not a comment claiming the limit is respected.

## When this applies beyond radio

The same method applies to SMS (160 characters), a QR code's data
capacity, or any message queue with a per-message size limit your
provider bills you for. The medium's ceiling is always non-negotiable;
your format is the only variable you control, so let the ceiling drive
the format instead of hoping your format fits.
