// Package jsonfile implements lending.ResourceRepository and
// lending.LoanRepository backed by a JSON file on disk. It's the "real"
// storage the demo server runs on: no database server to install, nothing
// to configure, and the file is human-readable if you want to look at it.
//
// This is deliberately not how a production deployment should store data —
// it re-reads and re-writes the whole file on every mutation, which won't
// scale past a small single-instance deployment. It exists to prove the
// swap: a Postgres or SQLite implementation of the same two interfaces
// drops in without lending.Service or the HTTP layer changing at all
// (CLAUDE.md rule 4: modular; rule 5: isolate what's likely to break —
// storage is exactly the kind of thing a community deployment will want to
// change without touching the business logic).
package jsonfile

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"toolshed/internal/domain"
)

type ResourceStore struct {
	mu   sync.Mutex
	path string
}

func NewResourceStore(path string) *ResourceStore {
	return &ResourceStore{path: path}
}

func (s *ResourceStore) load() (map[string]domain.Resource, error) {
	data := make(map[string]domain.Resource)
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return data, nil
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ResourceStore) save(data map[string]domain.Resource) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

func (s *ResourceStore) Create(_ context.Context, r domain.Resource) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return err
	}
	data[r.ID] = r
	return s.save(data)
}

func (s *ResourceStore) Get(_ context.Context, id string) (domain.Resource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return domain.Resource{}, err
	}
	r, ok := data[id]
	if !ok {
		return domain.Resource{}, domain.ErrNotFound("resource not found: " + id)
	}
	return r, nil
}

func (s *ResourceStore) List(_ context.Context) ([]domain.Resource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return nil, err
	}
	out := make([]domain.Resource, 0, len(data))
	for _, r := range data {
		out = append(out, r)
	}
	return out, nil
}

type LoanStore struct {
	mu   sync.Mutex
	path string
}

func NewLoanStore(path string) *LoanStore {
	return &LoanStore{path: path}
}

func (s *LoanStore) load() (map[string]domain.Loan, error) {
	data := make(map[string]domain.Loan)
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return data, nil
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *LoanStore) save(data map[string]domain.Loan) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

func (s *LoanStore) Create(_ context.Context, l domain.Loan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return err
	}
	data[l.ID] = l
	return s.save(data)
}

func (s *LoanStore) Get(_ context.Context, id string) (domain.Loan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return domain.Loan{}, err
	}
	l, ok := data[id]
	if !ok {
		return domain.Loan{}, domain.ErrNotFound("loan not found: " + id)
	}
	return l, nil
}

func (s *LoanStore) OpenLoanForResource(_ context.Context, resourceID string) (domain.Loan, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return domain.Loan{}, false, err
	}
	for _, l := range data {
		if l.ResourceID == resourceID && l.Open() {
			return l, true, nil
		}
	}
	return domain.Loan{}, false, nil
}

func (s *LoanStore) Update(_ context.Context, l domain.Loan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := s.load()
	if err != nil {
		return err
	}
	if _, ok := data[l.ID]; !ok {
		return domain.ErrNotFound("loan not found: " + l.ID)
	}
	data[l.ID] = l
	return s.save(data)
}
