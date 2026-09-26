package lending_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"toolshed/internal/domain"
	"toolshed/internal/lending"
	"toolshed/internal/storage/memory"
)

func newTestService() *lending.Service {
	return lending.NewService(memory.NewResourceStore(), memory.NewLoanStore())
}

func TestAddResource_RejectsInvalidKind(t *testing.T) {
	svc := newTestService()
	_, err := svc.AddResource(context.Background(), domain.Resource{Name: "Drill", Kind: "spaceship"})
	var invalid domain.ErrInvalidInput
	if !errors.As(err, &invalid) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCheckOutAndReturn_HappyPath(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	res, err := svc.AddResource(ctx, domain.Resource{Name: "Cordless Drill", Kind: domain.KindTool})
	if err != nil {
		t.Fatalf("AddResource: %v", err)
	}

	available, err := svc.Available(ctx, res.ID)
	if err != nil || !available {
		t.Fatalf("expected new resource to be available, got available=%v err=%v", available, err)
	}

	loan, err := svc.CheckOut(ctx, res.ID, "alice", 24*time.Hour)
	if err != nil {
		t.Fatalf("CheckOut: %v", err)
	}
	if loan.BorrowerID != "alice" || loan.ResourceID != res.ID {
		t.Fatalf("loan fields wrong: %+v", loan)
	}

	available, err = svc.Available(ctx, res.ID)
	if err != nil || available {
		t.Fatalf("expected checked-out resource to be unavailable, got available=%v err=%v", available, err)
	}

	returned, err := svc.Return(ctx, loan.ID)
	if err != nil {
		t.Fatalf("Return: %v", err)
	}
	if returned.ReturnedAt == nil {
		t.Fatal("expected ReturnedAt to be set")
	}

	available, err = svc.Available(ctx, res.ID)
	if err != nil || !available {
		t.Fatalf("expected returned resource to be available again, got available=%v err=%v", available, err)
	}
}

func TestCheckOut_RejectsDoubleBooking(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	res, _ := svc.AddResource(ctx, domain.Resource{Name: "3D Printer", Kind: domain.KindHardware})
	if _, err := svc.CheckOut(ctx, res.ID, "alice", time.Hour); err != nil {
		t.Fatalf("first CheckOut: %v", err)
	}

	_, err := svc.CheckOut(ctx, res.ID, "bob", time.Hour)
	var conflict domain.ErrConflict
	if !errors.As(err, &conflict) {
		t.Fatalf("expected ErrConflict on double booking, got %v", err)
	}
}

func TestCheckOut_UnknownResource(t *testing.T) {
	svc := newTestService()
	_, err := svc.CheckOut(context.Background(), "does-not-exist", "alice", time.Hour)
	var notFound domain.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestReturn_RejectsDoubleReturn(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	res, _ := svc.AddResource(ctx, domain.Resource{Name: "Repo access: farmOS fork", Kind: domain.KindSoftware})
	loan, _ := svc.CheckOut(ctx, res.ID, "alice", time.Hour)
	if _, err := svc.Return(ctx, loan.ID); err != nil {
		t.Fatalf("first Return: %v", err)
	}

	_, err := svc.Return(ctx, loan.ID)
	var conflict domain.ErrConflict
	if !errors.As(err, &conflict) {
		t.Fatalf("expected ErrConflict on double return, got %v", err)
	}
}
