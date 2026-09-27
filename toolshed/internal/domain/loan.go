package domain

import "time"

// Loan records one checkout of one Resource. A Resource with no open Loan
// (ReturnedAt set, or no Loan at all) is available; that's the entire
// availability rule — no separate "status" field to fall out of sync.
type Loan struct {
	ID           string     `json:"id"`
	ResourceID   string     `json:"resource_id"`
	BorrowerID   string     `json:"borrower_id"`
	CheckedOutAt time.Time  `json:"checked_out_at"`
	DueAt        time.Time  `json:"due_at"`
	ReturnedAt   *time.Time `json:"returned_at,omitempty"`
}

func (l Loan) Open() bool {
	return l.ReturnedAt == nil
}

func (l Loan) Overdue(now time.Time) bool {
	return l.Open() && now.After(l.DueAt)
}
