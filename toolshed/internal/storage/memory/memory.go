// Package memory implements lending.ResourceRepository and
// lending.LoanRepository by keeping everything in a map. It exists for
// tests (internal/lending/service_test.go uses it) and as the simplest
// possible example of "implement these two interfaces" for anyone reading
// this as a blueprint.
package memory

import (
	"context"
	"sync"

	"toolshed/internal/domain"
)

// ResourceStore implements lending.ResourceRepository.
type ResourceStore struct {
	mu   sync.Mutex
	data map[string]domain.Resource
}

func NewResourceStore() *ResourceStore {
	return &ResourceStore{data: make(map[string]domain.Resource)}
}

func (s *ResourceStore) Create(_ context.Context, r domain.Resource) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[r.ID] = r
	return nil
}

func (s *ResourceStore) Get(_ context.Context, id string) (domain.Resource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.data[id]
	if !ok {
		return domain.Resource{}, domain.ErrNotFound("resource not found: " + id)
	}
	return r, nil
}

func (s *ResourceStore) List(_ context.Context) ([]domain.Resource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Resource, 0, len(s.data))
	for _, r := range s.data {
		out = append(out, r)
	}
	return out, nil
}

// LoanStore implements lending.LoanRepository.
type LoanStore struct {
	mu   sync.Mutex
	data map[string]domain.Loan
}

func NewLoanStore() *LoanStore {
	return &LoanStore{data: make(map[string]domain.Loan)}
}

func (s *LoanStore) Create(_ context.Context, l domain.Loan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[l.ID] = l
	return nil
}

func (s *LoanStore) Get(_ context.Context, id string) (domain.Loan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.data[id]
	if !ok {
		return domain.Loan{}, domain.ErrNotFound("loan not found: " + id)
	}
	return l, nil
}

func (s *LoanStore) OpenLoanForResource(_ context.Context, resourceID string) (domain.Loan, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, l := range s.data {
		if l.ResourceID == resourceID && l.Open() {
			return l, true, nil
		}
	}
	return domain.Loan{}, false, nil
}

func (s *LoanStore) Update(_ context.Context, l domain.Loan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[l.ID]; !ok {
		return domain.ErrNotFound("loan not found: " + l.ID)
	}
	s.data[l.ID] = l
	return nil
}
