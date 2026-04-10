package billing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/33-token/model-api-platform/services/controlplane/internal/app"
)

func TestBillingHTTPCreateOrderAndQueryStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(app.NewServer().Handler)
	defer server.Close()

	createResp, err := http.Post(server.URL+"/api/v1/billing/orders", "application/json", bytes.NewBufferString(`{
		"accountId":"acct_http_1",
		"planPriceSnapshotId":"snapshot_basic",
		"amount":1900,
		"currency":"USD"
	}`))
	if err != nil {
		t.Fatalf("create order request failed: %v", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResp.StatusCode)
	}

	var created struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	if created.ID == "" || created.Status != "pending_payment" {
		t.Fatalf("unexpected created order: %+v", created)
	}

	getResp, err := http.Get(server.URL + "/api/v1/billing/orders/" + created.ID)
	if err != nil {
		t.Fatalf("get order request failed: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}
}

func TestBillingHTTPListsSellablePlans(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(app.NewServer().Handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/billing/plans")
	if err != nil {
		t.Fatalf("get plans request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Plans []struct {
			Code       string `json:"code"`
			QuotaTotal int64  `json:"quotaTotal"`
		} `json:"plans"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode plans response: %v", err)
	}

	if len(payload.Plans) != 1 || payload.Plans[0].Code != "basic" || payload.Plans[0].QuotaTotal != 1000000 {
		t.Fatalf("unexpected plans payload: %+v", payload)
	}
}

func TestBillingHTTPPaymentActivationAndQuotaQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(app.NewServer().Handler)
	defer server.Close()

	orderID := createTestOrder(t, server.URL, "acct_http_2")

	activateResp, err := http.Post(server.URL+"/api/v1/billing/payment-events", "application/json", bytes.NewBufferString(`{
		"providerEventId":"paypal_evt_1",
		"orderId":"`+orderID+`",
		"eventType":"payment_succeeded",
		"eventStatus":"paid"
	}`))
	if err != nil {
		t.Fatalf("apply payment request failed: %v", err)
	}
	defer activateResp.Body.Close()

	if activateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", activateResp.StatusCode)
	}

	quotaResp, err := http.Get(server.URL + "/api/v1/billing/account-quota?account_id=acct_http_2")
	if err != nil {
		t.Fatalf("get quota request failed: %v", err)
	}
	defer quotaResp.Body.Close()

	var quota struct {
		Status         string `json:"status"`
		QuotaTotal     int64  `json:"quotaTotal"`
		QuotaRemaining int64  `json:"quotaRemaining"`
	}
	if err := json.NewDecoder(quotaResp.Body).Decode(&quota); err != nil {
		t.Fatalf("decode quota response: %v", err)
	}

	if quota.Status != "active" || quota.QuotaTotal != 1000000 || quota.QuotaRemaining != 1000000 {
		t.Fatalf("unexpected quota view: %+v", quota)
	}
}

func TestBillingHTTPPaymentEventIsIdempotent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(app.NewServer().Handler)
	defer server.Close()

	orderID := createTestOrder(t, server.URL, "acct_http_4")

	for i := 0; i < 2; i++ {
		resp, err := http.Post(server.URL+"/api/v1/billing/payment-events", "application/json", bytes.NewBufferString(`{
			"providerEventId":"paypal_evt_same",
			"orderId":"`+orderID+`",
			"eventType":"payment_succeeded",
			"eventStatus":"paid"
		}`))
		if err != nil {
			t.Fatalf("payment event request failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected idempotent payment event to return 200, got %d", resp.StatusCode)
		}
	}

	orderResp, err := http.Get(server.URL + "/api/v1/billing/orders/" + orderID)
	if err != nil {
		t.Fatalf("get order after idempotent payment failed: %v", err)
	}
	defer orderResp.Body.Close()

	var order struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(orderResp.Body).Decode(&order); err != nil {
		t.Fatalf("decode order response: %v", err)
	}

	if order.Status != "completed" {
		t.Fatalf("expected completed order after repeated payment event, got %+v", order)
	}
}

func TestBillingHTTPSettlementPausesAccount(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(app.NewServer().Handler)
	defer server.Close()

	orderID := createTestOrder(t, server.URL, "acct_http_3")
	activateTestOrder(t, server.URL, orderID)

	settleResp, err := http.Post(server.URL+"/api/v1/billing/settlements", "application/json", bytes.NewBufferString(`{
		"accountId":"acct_http_3",
		"inputTokens":100,
		"outputTokens":200,
		"totalTokens":1000000,
		"quotaDelta":1000000
	}`))
	if err != nil {
		t.Fatalf("settlement request failed: %v", err)
	}
	defer settleResp.Body.Close()

	if settleResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", settleResp.StatusCode)
	}

	var quota struct {
		Status         string `json:"status"`
		QuotaRemaining int64  `json:"quotaRemaining"`
	}
	if err := json.NewDecoder(settleResp.Body).Decode(&quota); err != nil {
		t.Fatalf("decode settlement response: %v", err)
	}

	if quota.Status != "paused" || quota.QuotaRemaining != 0 {
		t.Fatalf("expected paused quota after settlement, got %+v", quota)
	}
}

func createTestOrder(t *testing.T, baseURL string, accountID string) string {
	t.Helper()

	resp, err := http.Post(baseURL+"/api/v1/billing/orders", "application/json", bytes.NewBufferString(`{
		"accountId":"`+accountID+`",
		"planPriceSnapshotId":"snapshot_basic",
		"amount":1900,
		"currency":"USD"
	}`))
	if err != nil {
		t.Fatalf("create test order request failed: %v", err)
	}
	defer resp.Body.Close()

	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create test order response: %v", err)
	}

	return created.ID
}

func activateTestOrder(t *testing.T, baseURL string, orderID string) {
	t.Helper()

	resp, err := http.Post(baseURL+"/api/v1/billing/payment-events", "application/json", bytes.NewBufferString(`{
		"providerEventId":"paypal_evt_activate",
		"orderId":"`+orderID+`",
		"eventType":"payment_succeeded",
		"eventStatus":"paid"
	}`))
	if err != nil {
		t.Fatalf("activate test order request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected payment activation 200, got %d", resp.StatusCode)
	}
}
