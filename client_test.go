package isbclient

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateLease(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/leases" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   Lease{UUID: uuid, LeaseId: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: "tpl", OriginalLeaseTemplateName: "tplname"},
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.CreateLease(context.Background(), &CreateLeaseRequest{LeaseTemplateUUID: "tpl", Comments: "test comment"})
	if err != nil {
		t.Fatalf("CreateLease error: %v", err)
	}
	if resp.Lease.LeaseId != leaseID {
		t.Errorf("expected lease LeaseId %s, got %s", leaseID, resp.Lease.LeaseId)
	}
	if resp.Lease.UUID != uuid {
		t.Errorf("expected lease UUID %s, got %s", uuid, resp.Lease.UUID)
	}
}

func TestUpdateLease(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/leases/"+leaseID {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   Lease{UUID: uuid, LeaseId: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: "tpl", OriginalLeaseTemplateName: "tplname"},
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.UpdateLease(context.Background(), &UpdateLeaseRequest{LeaseID: leaseID})
	if err != nil {
		t.Fatalf("UpdateLease error: %v", err)
	}
	if resp.Lease.LeaseId != leaseID {
		t.Errorf("expected lease LeaseId %s, got %s", leaseID, resp.Lease.LeaseId)
	}
	if resp.Lease.UUID != uuid {
		t.Errorf("expected lease UUID %s, got %s", uuid, resp.Lease.UUID)
	}
}

func TestReviewLease(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/leases/"+leaseID+"/review" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	err := client.ReviewLease(context.Background(), &ReviewLeaseRequest{LeaseID: leaseID, Action: ReviewApprove})
	if err != nil {
		t.Fatalf("ReviewLease error: %v", err)
	}
}

func TestFreezeLease(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/leases/"+leaseID+"/freeze" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	err := client.FreezeLease(context.Background(), &FreezeLeaseRequest{LeaseID: leaseID})
	if err != nil {
		t.Fatalf("FreezeLease error: %v", err)
	}
}

func TestTerminateLease(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/leases/"+leaseID+"/terminate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	err := client.TerminateLease(context.Background(), &TerminateLeaseRequest{LeaseID: leaseID})
	if err != nil {
		t.Fatalf("TerminateLease error: %v", err)
	}
}

func TestUpdateLeaseTemplate(t *testing.T) {
	tplID := "tpl123"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/leaseTemplates/"+tplID {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   LeaseTemplate{UUID: tplID, Name: "tpl", RequiresApproval: true, CreatedBy: "admin@example.com"},
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.UpdateLeaseTemplate(context.Background(), &UpdateLeaseTemplateRequest{LeaseTemplateID: tplID, Name: "tpl", Description: "desc", RequiresApproval: true, MaxSpend: 100, LeaseDurationInHours: 24, BudgetThresholds: nil, DurationThresholds: nil, CreatedBy: "admin@example.com"})
	if err != nil {
		t.Fatalf("UpdateLeaseTemplate error: %v", err)
	}
	if resp.LeaseTemplate.UUID != tplID {
		t.Errorf("expected template UUID %s, got %s", tplID, resp.LeaseTemplate.UUID)
	}
}

func TestDeleteLeaseTemplate(t *testing.T) {
	tplID := "tpl123"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/leaseTemplates/"+tplID {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	err := client.DeleteLeaseTemplate(context.Background(), &DeleteLeaseTemplateRequest{LeaseTemplateID: tplID})
	if err != nil {
		t.Fatalf("DeleteLeaseTemplate error: %v", err)
	}
}

func TestRegisterAccount(t *testing.T) {
	acctID := "123456789012"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/accounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   Account{AwsAccountId: acctID, Status: "Active"},
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.RegisterAccount(context.Background(), &RegisterAccountRequest{AwsAccountId: acctID})
	if err != nil {
		t.Fatalf("RegisterAccount error: %v", err)
	}
	if resp.Account.AwsAccountId != acctID {
		t.Errorf("expected account ID %s, got %s", acctID, resp.Account.AwsAccountId)
	}
}

func TestRetryCleanup(t *testing.T) {
	acctID := "123456789012"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/"+acctID+"/retryCleanup" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	err := client.RetryCleanup(context.Background(), &RetryCleanupRequest{AwsAccountId: acctID})
	if err != nil {
		t.Fatalf("RetryCleanup error: %v", err)
	}
}

func TestEjectAccount(t *testing.T) {
	acctID := "123456789012"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/"+acctID+"/eject" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	err := client.EjectAccount(context.Background(), &EjectAccountRequest{AwsAccountId: acctID})
	if err != nil {
		t.Fatalf("EjectAccount error: %v", err)
	}
}

func TestCreateLeaseAsUser(t *testing.T) {
	leaseID := "lease456"
	userEmail := "otheruser@example.com"
	jwtSecret := "testsecret"
	leaseTemplateUUID := "tpl"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/leases" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Errorf("missing Authorization header")
		} else {
			// Validate JWT
			tokenStr := auth[len("Bearer "):]
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				t.Errorf("invalid JWT: %v", err)
			}
			if claims.User.Email != userEmail {
				t.Errorf("expected user email %s in JWT, got %s", userEmail, claims.User.Email)
			}
			if len(claims.User.Roles) != 1 || claims.User.Roles[0] != RoleUser {
				t.Errorf("expected role %s in JWT, got %v", RoleUser, claims.User.Roles)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   Lease{UUID: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: leaseTemplateUUID, OriginalLeaseTemplateName: "tplname"},
		}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.CreateLeaseAsUser(context.Background(), &CreateLeaseRequest{LeaseTemplateUUID: leaseTemplateUUID, Comments: "as user"}, userEmail, jwtSecret)
	if err != nil {
		t.Fatalf("CreateLeaseAsUser error: %v", err)
	}
	if resp.Lease.UUID != leaseID {
		t.Errorf("expected lease UUID %s, got %s", leaseID, resp.Lease.UUID)
	}
	if resp.Lease.UserEmail != userEmail {
		t.Errorf("expected user email %s, got %s", userEmail, resp.Lease.UserEmail)
	}
}

func TestGetLeaseByID(t *testing.T) {
	leaseID := "lease-abc-123"
	leaseObj := Lease{
		UUID:      leaseID,
		UserEmail: "user@example.com",
		Status:    "Active",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/leases/"+leaseID {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   leaseObj,
		})
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.GetLeaseByID(context.Background(), &GetLeaseByIDRequest{LeaseID: leaseID})
	if err != nil {
		t.Fatalf("GetLeaseByID error: %v", err)
	}
	if resp.Lease.UUID != leaseID {
		t.Errorf("expected lease UUID %s, got %s", leaseID, resp.Lease.UUID)
	}
	if resp.Lease.UserEmail != "user@example.com" {
		t.Errorf("expected user email user@example.com, got %s", resp.Lease.UserEmail)
	}
	if resp.Lease.Status != "Active" {
		t.Errorf("expected status Active, got %s", resp.Lease.Status)
	}
}

func TestGetLeases(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/leases" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"result":             []Lease{{UUID: uuid, LeaseId: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: "tpl", OriginalLeaseTemplateName: "tplname"}},
				"nextPageIdentifier": "",
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.GetLeases(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetLeases error: %v", err)
	}
	if len(resp.Leases) != 1 {
		t.Fatalf("expected 1 lease, got %d", len(resp.Leases))
	}
	lease := resp.Leases[0]
	if lease.LeaseId != leaseID {
		t.Errorf("expected lease LeaseId %s, got %s", leaseID, lease.LeaseId)
	}
	if lease.UUID != uuid {
		t.Errorf("expected lease UUID %s, got %s", uuid, lease.UUID)
	}
}

func TestFetchAllLeases(t *testing.T) {
	uuid := "lease123"
	userEmail := "user@example.com"
	leaseIDObj := map[string]string{"userEmail": userEmail, "uuid": uuid}
	leaseIDJson, _ := json.Marshal(leaseIDObj)
	leaseID := b64.StdEncoding.EncodeToString(leaseIDJson)
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		var resp map[string]interface{}
		if callCount == 1 {
			resp = map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"result":             []Lease{{UUID: uuid, LeaseId: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: "tpl", OriginalLeaseTemplateName: "tplname"}},
					"nextPageIdentifier": "page2",
				},
			}
		} else if callCount == 2 {
			resp = map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"result":             []Lease{{UUID: uuid + "2", LeaseId: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: "tpl", OriginalLeaseTemplateName: "tplname"}},
					"nextPageIdentifier": "page3",
				},
			}
		} else {
			resp = map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"result":             []Lease{{UUID: uuid + "3", LeaseId: leaseID, UserEmail: userEmail, Status: "Active", OriginalLeaseTemplateUuid: "tpl", OriginalLeaseTemplateName: "tplname"}},
					"nextPageIdentifier": "",
				},
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	client := NewClient(server.URL, "token")
	resp, err := client.FetchAllLeases(context.Background(), &GetLeasesRequest{})
	if err != nil {
		t.Fatalf("FetchAllLeases error: %v", err)
	}
	if len(resp.Leases) != 3 {
		t.Fatalf("expected 3 leases, got %d", len(resp.Leases))
	}
	if resp.Leases[0].UUID != uuid {
		t.Errorf("expected first lease UUID %s, got %s", uuid, resp.Leases[0].UUID)
	}
	if resp.Leases[1].UUID != uuid+"2" {
		t.Errorf("expected second lease UUID %s, got %s", uuid+"2", resp.Leases[1].UUID)
	}
	if resp.Leases[2].UUID != uuid+"3" {
		t.Errorf("expected third lease UUID %s, got %s", uuid+"3", resp.Leases[2].UUID)
	}
	if callCount != 3 {
		t.Errorf("expected 3 pages to be fetched, got %d", callCount)
	}
}

func TestNonJSONResponses(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow non-JSON response tests in short mode")
	}

	t.Run("doGet returns error for HTML response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("<html><body>Error</body></html>"))
		}))
		defer ts.Close()
		client := NewClient(ts.URL, "token")
		_, err := client.doGet(context.Background(), ts.URL)
		if err == nil || err.Error() == "" || !contains(err.Error(), "non-JSON response") {
			t.Errorf("expected non-JSON response error, got %v", err)
		}
	})

	t.Run("doPost returns error for plain text response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("plain error"))
		}))
		defer ts.Close()
		client := NewClient(ts.URL, "token")
		_, err := client.doPost(context.Background(), ts.URL, []byte(`{}`))
		if err == nil || err.Error() == "" || !contains(err.Error(), "non-JSON response") {
			t.Errorf("expected non-JSON response error, got %v", err)
		}
	})

	t.Run("doPatch returns error for XML response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<error>forbidden</error>"))
		}))
		defer ts.Close()
		client := NewClient(ts.URL, "token")
		_, err := client.doPatch(context.Background(), ts.URL, []byte(`{}`))
		if err == nil || err.Error() == "" || !contains(err.Error(), "non-JSON response") {
			t.Errorf("expected non-JSON response error, got %v", err)
		}
	})

	t.Run("doPut returns error for octet-stream response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("binary error"))
		}))
		defer ts.Close()
		client := NewClient(ts.URL, "token")
		_, err := client.doPut(context.Background(), ts.URL, []byte(`{}`))
		if err == nil || err.Error() == "" || !contains(err.Error(), "non-JSON response") {
			t.Errorf("expected non-JSON response error, got %v", err)
		}
	})

	t.Run("doDelete returns error for XHTML response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<error>forbidden</error>"))
		}))
		defer ts.Close()
		client := NewClient(ts.URL, "token")
		_, err := client.doDelete(context.Background(), ts.URL)
		if err == nil || err.Error() == "" || !contains(err.Error(), "non-JSON response") {
			t.Errorf("expected non-JSON response error, got %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
}
