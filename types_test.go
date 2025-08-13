package isbclient

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetLeasesRequest_BuildQuery(t *testing.T) {
	tests := []struct {
		name  string
		input GetLeasesRequest
		want  url.Values
	}{
		{
			name:  "all fields set",
			input: GetLeasesRequest{PageIdentifier: "abc", PageSize: "10", UserEmail: "user@example.com"},
			want:  url.Values{"pageIdentifier": {"abc"}, "pageSize": {"10"}, "userEmail": {"user@example.com"}},
		},
		{
			name:  "only pageIdentifier",
			input: GetLeasesRequest{PageIdentifier: "abc"},
			want:  url.Values{"pageIdentifier": {"abc"}},
		},
		{
			name:  "empty",
			input: GetLeasesRequest{},
			want:  url.Values{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.BuildQuery()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLeaseByIDRequest_BuildQuery(t *testing.T) {
	r := GetLeaseByIDRequest{LeaseID: "lease-123"}
	want := url.Values{"leaseId": {"lease-123"}}
	got := r.BuildQuery()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildQuery() = %v, want %v", got, want)
	}
}

func TestCreateLeaseRequest_BuildQuery(t *testing.T) {
	r := CreateLeaseRequest{LeaseTemplateUUID: "uuid-1", Comments: "test comment"}
	want := url.Values{"leaseTemplateUuid": {"uuid-1"}, "comments": {"test comment"}}
	got := r.BuildQuery()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildQuery() = %v, want %v", got, want)
	}
}

func TestGetLeaseTemplatesRequest_BuildQuery(t *testing.T) {
	r := GetLeaseTemplatesRequest{PageIdentifier: "next", PageSize: "20"}
	want := url.Values{"pageIdentifier": {"next"}, "pageSize": {"20"}}
	got := r.BuildQuery()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildQuery() = %v, want %v", got, want)
	}
}

func TestGetAccountsRequest_BuildQuery(t *testing.T) {
	r := GetAccountsRequest{PageIdentifier: "next", PageSize: "50"}
	want := url.Values{"pageIdentifier": {"next"}, "pageSize": {"50"}}
	got := r.BuildQuery()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildQuery() = %v, want %v", got, want)
	}
}

func TestGetLeasesResponse_FilterByLeaseTemplateName(t *testing.T) {
	resp := &GetLeasesResponse{
		Leases: []Lease{
			{UUID: "1", OriginalLeaseTemplateName: "foo"},
			{UUID: "2", OriginalLeaseTemplateName: "bar"},
			{UUID: "3", OriginalLeaseTemplateName: "foo"},
		},
	}
	filtered := resp.FilterByLeaseTemplateName("foo")
	if len(filtered) != 2 || filtered[0].UUID != "1" || filtered[1].UUID != "3" {
		t.Errorf("expected leases with UUIDs 1 and 3, got %+v", filtered)
	}
}

func TestGetLeasesResponse_FilterByLeaseTemplateUUID(t *testing.T) {
	resp := &GetLeasesResponse{
		Leases: []Lease{
			{UUID: "1", OriginalLeaseTemplateUuid: "abc"},
			{UUID: "2", OriginalLeaseTemplateUuid: "def"},
			{UUID: "3", OriginalLeaseTemplateUuid: "abc"},
		},
	}
	filtered := resp.FilterByLeaseTemplateUUID("abc")
	if len(filtered) != 2 || filtered[0].UUID != "1" || filtered[1].UUID != "3" {
		t.Errorf("expected leases with UUIDs 1 and 3, got %+v", filtered)
	}
}

func TestCreateLeaseResponse(t *testing.T) {
	lease := Lease{UUID: "lease-1", UserEmail: "user@example.com"}
	resp := CreateLeaseResponse{Lease: lease}
	if resp.Lease.UUID != "lease-1" {
		t.Errorf("expected lease UUID 'lease-1', got %s", resp.Lease.UUID)
	}
}

func TestUpdateLeaseResponse(t *testing.T) {
	lease := Lease{UUID: "lease-2", UserEmail: "user2@example.com"}
	resp := UpdateLeaseResponse{Lease: lease}
	if resp.Lease.UserEmail != "user2@example.com" {
		t.Errorf("expected user 'user2@example.com', got %s", resp.Lease.UserEmail)
	}
}

func TestUpdateLeaseTemplateResponse(t *testing.T) {
	tpl := LeaseTemplate{UUID: "tpl-1", Name: "TestTpl"}
	resp := UpdateLeaseTemplateResponse{LeaseTemplate: tpl}
	if resp.LeaseTemplate.Name != "TestTpl" {
		t.Errorf("expected template name 'TestTpl', got %s", resp.LeaseTemplate.Name)
	}
}

func TestRegisterAccountResponse(t *testing.T) {
	acct := Account{AwsAccountId: "acc-1", Status: "Active"}
	resp := RegisterAccountResponse{Account: acct}
	if resp.Account.Status != "Active" {
		t.Errorf("expected status 'Active', got %s", resp.Account.Status)
	}
}

func TestNilReceivers_BuildQuery(t *testing.T) {
	t.Run("GetLeaseByIDRequest nil receiver", func(t *testing.T) {
		var r *GetLeaseByIDRequest
		got := r.BuildQuery()
		if got == nil || len(got) != 0 {
			t.Errorf("expected empty url.Values, got %v", got)
		}
	})
	t.Run("GetLeasesRequest nil receiver", func(t *testing.T) {
		var r *GetLeasesRequest
		got := r.BuildQuery()
		if got == nil || len(got) != 0 {
			t.Errorf("expected empty url.Values, got %v", got)
		}
	})
	t.Run("CreateLeaseRequest nil receiver", func(t *testing.T) {
		var r *CreateLeaseRequest
		got := r.BuildQuery()
		if got == nil || len(got) != 0 {
			t.Errorf("expected empty url.Values, got %v", got)
		}
	})
	t.Run("GetLeaseTemplatesRequest nil receiver", func(t *testing.T) {
		var r *GetLeaseTemplatesRequest
		got := r.BuildQuery()
		if got == nil || len(got) != 0 {
			t.Errorf("expected empty url.Values, got %v", got)
		}
	})
	t.Run("GetAccountsRequest nil receiver", func(t *testing.T) {
		var r *GetAccountsRequest
		got := r.BuildQuery()
		if got == nil || len(got) != 0 {
			t.Errorf("expected empty url.Values, got %v", got)
		}
	})
}
