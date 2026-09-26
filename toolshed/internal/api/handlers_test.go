package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"toolshed/internal/api"
	"toolshed/internal/domain"
	"toolshed/internal/lending"
	"toolshed/internal/storage/memory"
)

func newTestServer() *httptest.Server {
	svc := lending.NewService(memory.NewResourceStore(), memory.NewLoanStore())
	return httptest.NewServer(api.NewHandler(svc).Routes())
}

func postJSON(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	return resp
}

func TestFullLendingFlow(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	// Create a resource.
	createResp := postJSON(t, srv.URL+"/resources", map[string]any{
		"kind": "tool",
		"name": "Cordless Drill",
	})
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create resource: expected 201, got %d", createResp.StatusCode)
	}
	var created domain.Resource
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created resource to have an ID")
	}

	// List resources: should show it available.
	listResp, err := http.Get(srv.URL + "/resources")
	if err != nil {
		t.Fatalf("list resources: %v", err)
	}
	var listed []struct {
		domain.Resource
		Available bool `json:"available"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&listed); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(listed) != 1 || !listed[0].Available {
		t.Fatalf("expected one available resource, got %+v", listed)
	}

	// Check it out.
	checkoutResp := postJSON(t, srv.URL+"/resources/"+created.ID+"/checkout", map[string]any{
		"borrower_id":    "alice",
		"duration_hours": 48,
	})
	if checkoutResp.StatusCode != http.StatusCreated {
		t.Fatalf("checkout: expected 201, got %d", checkoutResp.StatusCode)
	}
	var loan domain.Loan
	if err := json.NewDecoder(checkoutResp.Body).Decode(&loan); err != nil {
		t.Fatalf("decode checkout response: %v", err)
	}

	// A second checkout attempt should conflict.
	secondCheckout := postJSON(t, srv.URL+"/resources/"+created.ID+"/checkout", map[string]any{
		"borrower_id": "bob",
	})
	if secondCheckout.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 on double checkout, got %d", secondCheckout.StatusCode)
	}

	// Return it.
	returnResp := postJSON(t, srv.URL+"/loans/"+loan.ID+"/return", nil)
	if returnResp.StatusCode != http.StatusOK {
		t.Fatalf("return: expected 200, got %d", returnResp.StatusCode)
	}

	// Now it should be checkout-able again.
	thirdCheckout := postJSON(t, srv.URL+"/resources/"+created.ID+"/checkout", map[string]any{
		"borrower_id": "carol",
	})
	if thirdCheckout.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 on re-checkout after return, got %d", thirdCheckout.StatusCode)
	}
}

func TestCreateResource_RejectsMissingName(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv.URL+"/resources", map[string]any{"kind": "tool"})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing name, got %d", resp.StatusCode)
	}
}
