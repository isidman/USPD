// Package lending holds the one business rule this whole slice exists to
// demonstrate: a Resource can have at most one open Loan at a time. Every
// other rule (who's allowed to borrow what, late fees, reservations) is a
// layer on top of this one — see BLUEPRINT.md.
package lending

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"toolshed/internal/domain"
)

// ResourceRepository and LoanRepository are defined here, by the consumer,
// not next to the structs they store — the standard Go way of keeping
// storage swappable (CLAUDE.md rule 4: modular, rule 5: isolate what's
// likely to break). internal/storage/memory and internal/storage/jsonfile
// both implement these same two interfaces; the Service never knows which
// one it's talking to.
type ResourceRepository interface {
	Create(ctx context.Context, r domain.Resource) error
	Get(ctx context.Context, id string) (domain.Resource, error)
	List(ctx context.Context) ([]domain.Resource, error)
}

type LoanRepository interface {
	Create(ctx context.Context, l domain.Loan) error
	Get(ctx context.Context, id string) (domain.Loan, error)
	// OpenLoanForResource reports the current open loan on a resource, if
	// any. This single query is what enforces "one open loan at a time."
	OpenLoanForResource(ctx context.Context, resourceID string) (domain.Loan, bool, error)
	Update(ctx context.Context, l domain.Loan) error
}

type Service struct {
	resources ResourceRepository
	loans     LoanRepository
	now       func() time.Time
	newID     func() string
}

func NewService(resources ResourceRepository, loans LoanRepository) *Service {
	return &Service{
		resources: resources,
		loans:     loans,
		now:       time.Now,
		newID:     randomID,
	}
}

func randomID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Service) AddResource(ctx context.Context, r domain.Resource) (domain.Resource, error) {
	if err := r.Validate(); err != nil {
		return domain.Resource{}, err
	}
	r.ID = s.newID()
	r.CreatedAt = s.now()
	if err := s.resources.Create(ctx, r); err != nil {
		return domain.Resource{}, err
	}
	return r, nil
}

func (s *Service) ListResources(ctx context.Context) ([]domain.Resource, error) {
	return s.resources.List(ctx)
}

// Available reports whether a resource has no open loan against it.
func (s *Service) Available(ctx context.Context, resourceID string) (bool, error) {
	_, open, err := s.loans.OpenLoanForResource(ctx, resourceID)
	if err != nil {
		return false, err
	}
	return !open, nil
}

// CheckOut lends a resource to a borrower for the given duration. It fails
// if the resource doesn't exist or already has an open loan — that check
// is the whole point of this package.
func (s *Service) CheckOut(ctx context.Context, resourceID, borrowerID string, duration time.Duration) (domain.Loan, error) {
	if borrowerID == "" {
		return domain.Loan{}, domain.ErrInvalidInput("borrower_id is required")
	}
	if _, err := s.resources.Get(ctx, resourceID); err != nil {
		return domain.Loan{}, err
	}
	if _, open, err := s.loans.OpenLoanForResource(ctx, resourceID); err != nil {
		return domain.Loan{}, err
	} else if open {
		return domain.Loan{}, domain.ErrConflict("resource is already checked out")
	}

	now := s.now()
	loan := domain.Loan{
		ID:           s.newID(),
		ResourceID:   resourceID,
		BorrowerID:   borrowerID,
		CheckedOutAt: now,
		DueAt:        now.Add(duration),
	}
	if err := s.loans.Create(ctx, loan); err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}

// Return closes an open loan. It fails if the loan doesn't exist or was
// already returned.
func (s *Service) Return(ctx context.Context, loanID string) (domain.Loan, error) {
	loan, err := s.loans.Get(ctx, loanID)
	if err != nil {
		return domain.Loan{}, err
	}
	if !loan.Open() {
		return domain.Loan{}, domain.ErrConflict("loan was already returned")
	}
	now := s.now()
	loan.ReturnedAt = &now
	if err := s.loans.Update(ctx, loan); err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}
