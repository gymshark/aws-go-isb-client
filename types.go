package isbclient

import (
	"encoding/json"
	"net/url"
)

const (
	ReviewApprove = "Approve"
	ReviewDeny    = "Deny"

	StatusActive             = "Active"
	StatusDenied             = "ApprovalDenied"
	StatusManuallyTerminated = "ManuallyTerminated"
	StatusFrozen             = "Frozen"
	StatusExpired            = "Expired"
)

// FailResponseBody represents a failed API response.
type FailResponseBody struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// ErrorResponseBody represents an error response from the API.
type ErrorResponseBody struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// BudgetThreshold is an array of objects with specific fields
type BudgetThreshold struct {
	DollarsSpent float64 `json:"dollarsSpent"`
	Action       string  `json:"action"`
}

// DurationThreshold is an array of objects with specific fields
type DurationThreshold struct {
	HoursRemaining float64 `json:"hoursRemaining"`
	Action         string  `json:"action"`
}

type MetaData struct {
	CreatedTime   string      `json:"createdTime"`
	LastEditTime  string      `json:"lastEditTime"`
	SchemaVersion json.Number `json:"schemaVersion"`
}

// Lease represents a lease object (fully defined)
type Lease struct {
	UserEmail                 string              `json:"userEmail"`
	UUID                      string              `json:"uuid"`
	Status                    string              `json:"status"`
	OriginalLeaseTemplateUuid string              `json:"originalLeaseTemplateUuid"`
	OriginalLeaseTemplateName string              `json:"originalLeaseTemplateName"`
	LeaseDurationInHours      int                 `json:"leaseDurationInHours"`
	MaxSpend                  float64             `json:"maxSpend"`
	BudgetThresholds          []BudgetThreshold   `json:"budgetThresholds"`
	DurationThresholds        []DurationThreshold `json:"durationThresholds"`
	Comments                  string              `json:"comments"`
	AwsAccountId              string              `json:"awsAccountId"`
	LeaseId                   string              `json:"leaseId"`
	StartDate                 string              `json:"startDate"`
	ExpirationDate            string              `json:"expirationDate"`
	EndDate                   string              `json:"endDate"`
	TotalCostAccrued          float64             `json:"totalCostAccrued"`
	Meta                      MetaData            `json:"meta"`
}

// LeaseTemplate represents a lease template (fully defined)
type LeaseTemplate struct {
	UUID                 string              `json:"uuid"`
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	RequiresApproval     bool                `json:"requiresApproval"`
	CreatedBy            string              `json:"createdBy"`
	MaxSpend             float64             `json:"maxSpend"`
	BudgetThresholds     []BudgetThreshold   `json:"budgetThresholds"`
	LeaseDurationInHours int                 `json:"leaseDurationInHours"`
	DurationThresholds   []DurationThreshold `json:"durationThresholds"`
	Meta                 MetaData            `json:"meta"`
}

// Account represents an account (fully defined)
type Account struct {
	AwsAccountId    string   `json:"awsAccountId"`
	Status          string   `json:"status"`
	DriftAtLastScan bool     `json:"driftAtLastScan"`
	Meta            MetaData `json:"meta"`
}

// UnregisteredAccount represents an unregistered account
type UnregisteredAccount struct {
	Id              string `json:"Id"`
	Arn             string `json:"Arn"`
	Email           string `json:"Email"`
	Name            string `json:"Name"`
	Status          string `json:"Status"`
	JoinedMethod    string `json:"JoinedMethod"`
	JoinedTimestamp string `json:"JoinedTimestamp"`
}

// GlobalConfiguration represents the global config (fully defined)
type GlobalConfiguration struct {
	TermsOfService  string                   `json:"termsOfService"`
	MaintenanceMode bool                     `json:"maintenanceMode"`
	Leases          GlobalLeasesConfig       `json:"leases"`
	Cleanup         GlobalCleanupConfig      `json:"cleanup"`
	Auth            map[string]interface{}   `json:"auth"`
	Notification    GlobalNotificationConfig `json:"notification"`
}

type GlobalLeasesConfig struct {
	MaxBudget                     float64 `json:"maxBudget"`
	DefaultBudgetThresholds       []int   `json:"defaultBudgetThresholds"`
	DefaultDurationThresholds     []int   `json:"defaultDurationThresholds"`
	MaxBudgetReclamationThreshold int     `json:"maxBudgetReclamationThreshold"`
	MaxDurationHours              float64 `json:"maxDurationHours"`
	MaxLeasesPerUser              int     `json:"maxLeasesPerUser"`
}

type GlobalCleanupConfig struct {
	NumberOfFailedAttemptsToCancelCleanup     int `json:"numberOfFailedAttemptsToCancelCleanup"`
	WaitBeforeRetryFailedAttemptSeconds       int `json:"waitBeforeRetryFailedAttemptSeconds"`
	NumberOfSuccessfulAttemptsToFinishCleanup int `json:"numberOfSuccessfulAttemptsToFinishCleanup"`
	WaitBeforeRerunSuccessfulAttemptSeconds   int `json:"waitBeforeRerunSuccessfulAttemptSeconds"`
}

type GlobalNotificationConfig struct {
	EmailFrom string `json:"emailFrom"`
}

// PageRequestBase provides a base for paginated requests.
type PageRequestBase struct {
	PageIdentifier string
}

func (b *PageRequestBase) SetPageIdentifier(next string) {
	b.PageIdentifier = next
}

type PageIdentifiable interface {
	SetPageIdentifier(string)
}

// Request parameter structs

type QueryBuilder interface {
	BuildQuery() url.Values
}

type GetLeasesRequest struct {
	PageIdentifier string
	PageSize       string
	UserEmail      string
}

func (r *GetLeasesRequest) SetPageIdentifier(next string) {
	r.PageIdentifier = next
}

func (r *GetLeasesRequest) BuildQuery() url.Values {
	if r == nil {
		return url.Values{}
	}
	q := url.Values{}
	if r.PageIdentifier != "" {
		q.Set("pageIdentifier", r.PageIdentifier)
	}
	if r.PageSize != "" {
		q.Set("pageSize", r.PageSize)
	}
	if r.UserEmail != "" {
		q.Set("userEmail", r.UserEmail)
	}
	return q
}

type GetLeaseByIDRequest struct {
	LeaseID string
}

func (r *GetLeaseByIDRequest) BuildQuery() url.Values {
	if r == nil {
		return url.Values{}
	}
	q := url.Values{}
	if r.LeaseID != "" {
		q.Set("leaseId", r.LeaseID)
	}
	return q
}

type CreateLeaseRequest struct {
	LeaseTemplateUUID string
	Comments          string
}

func (r *CreateLeaseRequest) BuildQuery() url.Values {
	if r == nil {
		return url.Values{}
	}
	q := url.Values{}
	if r.LeaseTemplateUUID != "" {
		q.Set("leaseTemplateUuid", r.LeaseTemplateUUID)
	}
	if r.Comments != "" {
		q.Set("comments", r.Comments)
	}
	return q
}

type GetLeaseTemplatesRequest struct {
	PageIdentifier string
	PageSize       string
}

func (r *GetLeaseTemplatesRequest) SetPageIdentifier(next string) {
	r.PageIdentifier = next
}

func (r *GetLeaseTemplatesRequest) BuildQuery() url.Values {
	if r == nil {
		return url.Values{}
	}
	q := url.Values{}
	if r.PageIdentifier != "" {
		q.Set("pageIdentifier", r.PageIdentifier)
	}
	if r.PageSize != "" {
		q.Set("pageSize", r.PageSize)
	}
	return q
}

type GetAccountsRequest struct {
	PageIdentifier string
	PageSize       string
}

func (r *GetAccountsRequest) SetPageIdentifier(next string) {
	r.PageIdentifier = next
}

func (r *GetAccountsRequest) BuildQuery() url.Values {
	if r == nil {
		return url.Values{}
	}
	q := url.Values{}
	if r.PageIdentifier != "" {
		q.Set("pageIdentifier", r.PageIdentifier)
	}
	if r.PageSize != "" {
		q.Set("pageSize", r.PageSize)
	}
	return q
}

// Paginated and result wrapper structs

type PaginatedResults[T any] struct {
	Items              []T    `json:"items"`
	NextPageIdentifier string `json:"nextPageIdentifier"`
}

type PaginatedLeases struct {
	Items              []Lease `json:"items"`
	NextPageIdentifier string  `json:"nextPageIdentifier"`
	Result             []Lease `json:"result"`
}

type PaginatedLeaseTemplates struct {
	Items              []LeaseTemplate `json:"items"`
	NextPageIdentifier string          `json:"nextPageIdentifier"`
	Result             []LeaseTemplate `json:"result"`
}

type PaginatedAccounts struct {
	Items              []Account `json:"items"`
	NextPageIdentifier string    `json:"nextPageIdentifier"`
	Result             []Account `json:"result"`
}

type PaginatedUnregisteredAccounts struct {
	Items              []UnregisteredAccount `json:"items"`
	NextPageIdentifier string                `json:"nextPageIdentifier"`
	Result             []UnregisteredAccount `json:"result"`
}

// Response wrapper structs for client methods

type GetLeasesResponse struct {
	Leases             []Lease `json:"result"`
	NextPageIdentifier string  `json:"nextPageIdentifier,omitempty"`
}

// FilterByLeaseTemplateName is a helper to filter leases by LeaseTemplateName
func (r *GetLeasesResponse) FilterByLeaseTemplateName(name string) []Lease {
	var filtered []Lease
	for _, l := range r.Leases {
		if l.OriginalLeaseTemplateName == name {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// FilterByLeaseTemplateUUID is a helper to filter leases by LeaseTemplateUUID
func (r *GetLeasesResponse) FilterByLeaseTemplateUUID(uuid string) []Lease {
	var filtered []Lease
	for _, l := range r.Leases {
		if l.OriginalLeaseTemplateUuid == uuid {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

type GetLeaseTemplatesResponse struct {
	LeaseTemplates     []LeaseTemplate `json:"leaseTemplates"`
	NextPageIdentifier string          `json:"nextPageIdentifier,omitempty"`
}

type GetAccountsResponse struct {
	Accounts           []Account `json:"accounts"`
	NextPageIdentifier string    `json:"nextPageIdentifier,omitempty"`
}

type GetUnregisteredAccountsResponse struct {
	UnregisteredAccounts []UnregisteredAccount `json:"unregisteredAccounts"`
	NextPageIdentifier   string                `json:"nextPageIdentifier,omitempty"`
}

// CreateLeaseResponse represents the response for creating a lease.
// POST /leases
// Contains a single Lease.
type CreateLeaseResponse struct {
	Lease Lease `json:"data"`
}

// UpdateLeaseResponse represents the response for updating a lease.
// PATCH /leases/{leaseId}
// Contains a single Lease.
type UpdateLeaseResponse struct {
	Lease Lease `json:"data"`
}

// UpdateLeaseTemplateResponse represents the response for updating a lease template.
// PUT /leaseTemplates/{leaseTemplateId}
// Contains a single LeaseTemplate.
type UpdateLeaseTemplateResponse struct {
	LeaseTemplate LeaseTemplate `json:"data"`
}

// RegisterAccountResponse represents the response for registering an account.
// POST /accounts
// Contains a single Account.
type RegisterAccountResponse struct {
	Account Account `json:"data"`
}

// UpdateLeaseRequest represents a request to update a lease.
// PATCH /leases/{leaseId}
type UpdateLeaseRequest struct {
	LeaseID            string
	MaxSpend           *float64             `json:"maxSpend,omitempty"`
	BudgetThresholds   *[]BudgetThreshold   `json:"budgetThresholds,omitempty"`
	ExpirationDate     *string              `json:"expirationDate,omitempty"`
	DurationThresholds *[]DurationThreshold `json:"durationThresholds,omitempty"`
}

// ReviewLeaseRequest represents a request to review (approve/deny) a lease.
// POST /leases/{leaseId}/review
type ReviewLeaseRequest struct {
	LeaseID string
	Action  string `json:"action"` // Approve or Deny
}

// FreezeLeaseRequest represents a request to freeze an active lease (no body).
// POST /leases/{leaseId}/freeze
type FreezeLeaseRequest struct {
	LeaseID string
}

// TerminateLeaseRequest represents a request to terminate a lease (no body).
// POST /leases/{leaseId}/terminate
type TerminateLeaseRequest struct {
	LeaseID string
}

// UpdateLeaseTemplateRequest represents a request to update a lease template.
// PUT /leaseTemplates/{leaseTemplateId}
type UpdateLeaseTemplateRequest struct {
	LeaseTemplateID      string
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	RequiresApproval     bool                `json:"requiresApproval"`
	MaxSpend             float64             `json:"maxSpend"`
	LeaseDurationInHours int                 `json:"leaseDurationInHours"`
	BudgetThresholds     []BudgetThreshold   `json:"budgetThresholds"`
	DurationThresholds   []DurationThreshold `json:"durationThresholds"`
	CreatedBy            string              `json:"createdBy"`
}

// DeleteLeaseTemplateRequest represents a request to delete a lease template (no body).
// DELETE /leaseTemplates/{leaseTemplateId}
type DeleteLeaseTemplateRequest struct {
	LeaseTemplateID string
}

// RegisterAccountRequest represents a request to register an account.
// POST /accounts
type RegisterAccountRequest struct {
	AwsAccountId string `json:"awsAccountId"`
}

// RetryCleanupRequest represents a request to retry cleanup for an account (no body).
// POST /accounts/{awsAccountId}/retryCleanup
type RetryCleanupRequest struct {
	AwsAccountId string
}

// EjectAccountRequest represents a request to eject an account from the sandbox (no body).
// POST /accounts/{awsAccountId}/eject
type EjectAccountRequest struct {
	AwsAccountId string
}

// GetLeaseByIDResponse is the response struct for GetLeaseByID
// (not paginated, always a single lease)
type GetLeaseByIDResponse struct {
	Lease Lease `json:"lease"`
}
