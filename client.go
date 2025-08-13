package isbclient

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the HTTP client for the Innovation Sandbox API.
// It supports bearer token authentication.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// authTransport is a custom RoundTripper that injects the Authorization header.
type authTransport struct {
	base  http.RoundTripper
	token string
}

// RoundTrip implements the http.RoundTripper interface.
func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if a.token != "" {
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", "Bearer "+a.token)
	}
	return a.base.RoundTrip(req)
}

// NewClient creates a new API client with recommended timeouts and settings.
func NewClient(baseURL, token string) *Client {
	baseTransport := http.DefaultTransport
	httpClient := &http.Client{
		Timeout: 15 * time.Second, // 15 seconds
		Transport: &authTransport{
			base:  baseTransport,
			token: token,
		},
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
		Token:      token,
	}
}

// GetLeases fetches a paginated list of leases and returns typed data
func (c *Client) GetLeases(ctx context.Context, req QueryBuilder) (*GetLeasesResponse, error) {
	u, err := url.Parse(c.BaseURL + "/leases")
	if err != nil {
		return nil, &APIRequestError{Op: "parse", URL: c.BaseURL + "/leases", Err: err}
	}

	if req != nil {
		u.RawQuery = req.BuildQuery().Encode()
	}

	resp, err := c.doGet(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Status string            `json:"status"`
		Data   GetLeasesResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}

	return &wrapper.Data, nil
}

// GetLeaseByID fetches a lease by its ID and returns typed data
func (c *Client) GetLeaseByID(ctx context.Context, req *GetLeaseByIDRequest) (*GetLeaseByIDResponse, error) {
	if req == nil || req.LeaseID == "" {
		return nil, &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseID is required")}
	}
	leaseURL := fmt.Sprintf("%s/leases/%s", c.BaseURL, req.LeaseID)
	resp, err := c.doGet(ctx, leaseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Status string `json:"status"`
		Data   Lease  `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}
	return &GetLeaseByIDResponse{Lease: wrapper.Data}, nil
}

// CreateLease requests a new lease and returns the created Lease in a response struct
func (c *Client) CreateLease(ctx context.Context, req *CreateLeaseRequest) (*CreateLeaseResponse, error) {
	if req == nil || req.LeaseTemplateUUID == "" {
		return nil, &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseTemplateUUID is required")}
	}
	leaseURL := c.BaseURL + "/leases"
	body := map[string]interface{}{
		"leaseTemplateUuid": req.LeaseTemplateUUID,
	}
	if req.Comments != "" {
		body["comments"] = req.Comments
	}
	b, _ := json.Marshal(body)
	resp, err := c.doPost(ctx, leaseURL, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Status string `json:"status"`
		Data   Lease  `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}

	leaseIdComponents := map[string]string{
		"userEmail": wrapper.Data.UserEmail,
		"uuid":      wrapper.Data.UUID,
	}

	leaseId, err := json.Marshal(leaseIdComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal lease ID components: %w", err)
	}

	wrapper.Data.LeaseId = b64.StdEncoding.EncodeToString(leaseId)

	return &CreateLeaseResponse{Lease: wrapper.Data}, nil
}

// CreateLeaseAsUser creates a lease as a different user by generating a JWT for that user and using it for the request only.
func (c *Client) CreateLeaseAsUser(ctx context.Context, req *CreateLeaseRequest, userEmail string, jwtSecret string) (*CreateLeaseResponse, error) {
	if req == nil || req.LeaseTemplateUUID == "" {
		return nil, &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseTemplateUUID is required")}
	}
	leaseURL := c.BaseURL + "/leases"
	body := map[string]interface{}{
		"leaseTemplateUuid": req.LeaseTemplateUUID,
	}
	if req.Comments != "" {
		body["comments"] = req.Comments
	}
	b, _ := json.Marshal(body)

	// Generate JWT using helper
	userClaims := NewUserUserClaims(userEmail)
	jwt, err := GenerateJWT(userClaims, jwtSecret, 15*time.Minute)
	if err != nil {
		return nil, &APIRequestError{Op: "jwt_gen", URL: leaseURL, Err: err}
	}

	// Use a custom request with the user JWT, but otherwise match doPost logic
	httpReq, err := http.NewRequestWithContext(ctx, "POST", leaseURL, bytes.NewReader(b))
	if err != nil {
		return nil, &APIRequestError{Op: "new_request", URL: leaseURL, Err: err}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+jwt)

	// Use a custom HTTP client that does NOT use the default authTransport for this request
	customClient := &http.Client{
		Timeout: c.HTTPClient.Timeout,
	}
	resp, err := customClient.Do(httpReq)
	if err != nil {
		return nil, &APIRequestError{Op: "do", URL: leaseURL, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, DecodeAPIError(b, resp)
	}

	var wrapper struct {
		Status string `json:"status"`
		Data   Lease  `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}

	leaseIdComponents := map[string]string{
		"userEmail": wrapper.Data.UserEmail,
		"uuid":      wrapper.Data.UUID,
	}

	leaseId, err := json.Marshal(leaseIdComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal lease ID components: %w", err)
	}

	wrapper.Data.LeaseId = b64.StdEncoding.EncodeToString(leaseId)

	return &CreateLeaseResponse{Lease: wrapper.Data}, nil
}

// GetLeaseTemplates fetches lease templates and returns typed data
func (c *Client) GetLeaseTemplates(ctx context.Context, req QueryBuilder) (*GetLeaseTemplatesResponse, error) {
	u, err := url.Parse(c.BaseURL + "/leaseTemplates")
	if err != nil {
		return nil, &APIRequestError{Op: "parse", URL: c.BaseURL + "/leaseTemplates", Err: err}
	}

	if req != nil {
		u.RawQuery = req.BuildQuery().Encode()
	}

	resp, err := c.doGet(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Status string                    `json:"status"`
		Data   GetLeaseTemplatesResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}

	return &wrapper.Data, nil
}

// FetchAllLeases fetches all leases using pagination
func (c *Client) FetchAllLeases(ctx context.Context, req *GetLeasesRequest) (*GetLeasesResponse, error) {
	allLeases, err := paginateAll(ctx, req, func(ctx context.Context, r *GetLeasesRequest) ([]Lease, string, error) {
		resp, err := c.GetLeases(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.Leases, resp.NextPageIdentifier, nil
	})
	if err != nil {
		return nil, err
	}
	return &GetLeasesResponse{Leases: allLeases}, nil
}

// FetchAllLeaseTemplates fetches all lease templates using pagination
func (c *Client) FetchAllLeaseTemplates(ctx context.Context, req *GetLeaseTemplatesRequest) (*GetLeaseTemplatesResponse, error) {
	allTemplates, err := paginateAll(ctx, req, func(ctx context.Context, r *GetLeaseTemplatesRequest) ([]LeaseTemplate, string, error) {
		resp, err := c.GetLeaseTemplates(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.LeaseTemplates, resp.NextPageIdentifier, nil
	})
	if err != nil {
		return nil, err
	}
	return &GetLeaseTemplatesResponse{LeaseTemplates: allTemplates}, nil
}

// GetAccounts fetches accounts and returns typed data
func (c *Client) GetAccounts(ctx context.Context, req QueryBuilder) (*GetAccountsResponse, error) {
	u, err := url.Parse(c.BaseURL + "/accounts")
	if err != nil {
		return nil, &APIRequestError{Op: "parse", URL: c.BaseURL + "/accounts", Err: err}
	}

	if req != nil {
		u.RawQuery = req.BuildQuery().Encode()
	}

	resp, err := c.doGet(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Status string              `json:"status"`
		Data   GetAccountsResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}

	return &wrapper.Data, nil
}

// FetchAllAccounts fetches all accounts using pagination
func (c *Client) FetchAllAccounts(ctx context.Context, req *GetAccountsRequest) (*GetAccountsResponse, error) {
	allAccounts, err := paginateAll(ctx, req, func(ctx context.Context, r *GetAccountsRequest) ([]Account, string, error) {
		resp, err := c.GetAccounts(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.Accounts, resp.NextPageIdentifier, nil
	})
	if err != nil {
		return nil, err
	}
	return &GetAccountsResponse{Accounts: allAccounts}, nil
}

// GetConfigurations fetches the global configuration
func (c *Client) GetConfigurations(ctx context.Context) (*GlobalConfiguration, error) {
	configURL := c.BaseURL + "/configurations"
	resp, err := c.doGet(ctx, configURL)
	err = nil
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var wrapper struct {
		Status string              `json:"status"`
		Data   GlobalConfiguration `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}

	return &wrapper.Data, nil
}

// UpdateLease updates a lease by leaseId (PATCH /leases/{leaseId})
func (c *Client) UpdateLease(ctx context.Context, req *UpdateLeaseRequest) (*UpdateLeaseResponse, error) {
	if req == nil || req.LeaseID == "" {
		return nil, &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseID is required")}
	}
	urlStr := c.BaseURL + "/leases/" + req.LeaseID
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &APIRequestError{Op: "marshal", URL: urlStr, Err: err}
	}
	resp, err := c.doPatch(ctx, urlStr, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var wrapper struct {
		Status string `json:"status"`
		Data   Lease  `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}
	return &UpdateLeaseResponse{Lease: wrapper.Data}, nil
}

// ReviewLease reviews (approve/deny) a lease (POST /leases/{leaseId}/review)
func (c *Client) ReviewLease(ctx context.Context, req *ReviewLeaseRequest) error {
	if req == nil || req.LeaseID == "" || (req.Action != ReviewApprove && req.Action != ReviewDeny) {
		return &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseID and Action are required")}
	}
	urlStr := c.BaseURL + "/leases/" + req.LeaseID + "/review"
	body, err := json.Marshal(map[string]string{"action": req.Action})
	if err != nil {
		return &APIRequestError{Op: "marshal", URL: urlStr, Err: err}
	}
	resp, err := c.doPost(ctx, urlStr, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// FreezeLease freezes an active lease (POST /leases/{leaseId}/freeze)
func (c *Client) FreezeLease(ctx context.Context, req *FreezeLeaseRequest) error {
	if req == nil || req.LeaseID == "" {
		return &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseID is required")}
	}
	urlStr := c.BaseURL + "/leases/" + req.LeaseID + "/freeze"
	resp, err := c.doPost(ctx, urlStr, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// TerminateLease terminates an active lease (POST /leases/{leaseId}/terminate)
func (c *Client) TerminateLease(ctx context.Context, req *TerminateLeaseRequest) error {
	if req == nil || req.LeaseID == "" {
		return &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseID is required")}
	}
	urlStr := c.BaseURL + "/leases/" + req.LeaseID + "/terminate"
	resp, err := c.doPost(ctx, urlStr, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// UpdateLeaseTemplate updates a lease template (PUT /leaseTemplates/{leaseTemplateId})
func (c *Client) UpdateLeaseTemplate(ctx context.Context, req *UpdateLeaseTemplateRequest) (*UpdateLeaseTemplateResponse, error) {
	if req == nil || req.LeaseTemplateID == "" {
		return nil, &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseTemplateID is required")}
	}
	urlStr := c.BaseURL + "/leaseTemplates/" + req.LeaseTemplateID
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &APIRequestError{Op: "marshal", URL: urlStr, Err: err}
	}
	resp, err := c.doPut(ctx, urlStr, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var wrapper struct {
		Status string        `json:"status"`
		Data   LeaseTemplate `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}
	return &UpdateLeaseTemplateResponse{LeaseTemplate: wrapper.Data}, nil
}

// DeleteLeaseTemplate deletes a lease template (DELETE /leaseTemplates/{leaseTemplateId})
func (c *Client) DeleteLeaseTemplate(ctx context.Context, req *DeleteLeaseTemplateRequest) error {
	if req == nil || req.LeaseTemplateID == "" {
		return &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("LeaseTemplateID is required")}
	}
	urlStr := c.BaseURL + "/leaseTemplates/" + req.LeaseTemplateID
	resp, err := c.doDelete(ctx, urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// RegisterAccount registers an account (POST /accounts)
func (c *Client) RegisterAccount(ctx context.Context, req *RegisterAccountRequest) (*RegisterAccountResponse, error) {
	urlStr := c.BaseURL + "/accounts"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &APIRequestError{Op: "marshal", URL: urlStr, Err: err}
	}
	resp, err := c.doPost(ctx, urlStr, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var wrapper struct {
		Status string  `json:"status"`
		Data   Account `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, &JSONDecodingError{Err: err}
	}
	return &RegisterAccountResponse{Account: wrapper.Data}, nil
}

// RetryCleanup retries cleanup for an account (POST /accounts/{awsAccountId}/retryCleanup)
func (c *Client) RetryCleanup(ctx context.Context, req *RetryCleanupRequest) error {
	if req == nil || req.AwsAccountId == "" {
		return &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("AwsAccountId is required")}
	}
	urlStr := c.BaseURL + "/accounts/" + req.AwsAccountId + "/retryCleanup"
	resp, err := c.doPost(ctx, urlStr, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// EjectAccount ejects an account from the sandbox (POST /accounts/{awsAccountId}/eject)
func (c *Client) EjectAccount(ctx context.Context, req *EjectAccountRequest) error {
	if req == nil || req.AwsAccountId == "" {
		return &APIRequestError{Op: "param", URL: "", Err: fmt.Errorf("AwsAccountId is required")}
	}
	urlStr := c.BaseURL + "/accounts/" + req.AwsAccountId + "/eject"
	resp, err := c.doPost(ctx, urlStr, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// paginateAll is a generic helper for paginated API fetches (no reflection needed)
func paginateAll[T any, R PageIdentifiable](
	ctx context.Context,
	req R,
	fetchPage func(context.Context, R) ([]T, string, error),
) ([]T, error) {
	var allItems []T
	for {
		items, nextPage, err := fetchPage(ctx, req)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items...)
		if nextPage == "" {
			break
		}
		req.SetPageIdentifier(nextPage)
	}
	return allItems, nil
}

// doGet is a helper for making GET requests and handling common errors.
func (c *Client) doGet(ctx context.Context, url string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, &APIRequestError{Op: "new_request", URL: url, Err: err}
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &APIRequestError{Op: "do", URL: url, Err: err}
	}

	if isJSON, body := isJSONResponse(resp); !isJSON {
		defer resp.Body.Close()
		return nil, &APIRequestError{
			Op:  "doGet",
			URL: url,
			Err: fmt.Errorf("non-JSON response (%s): %s", resp.Header.Get("Content-Type"), body),
		}
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, DecodeAPIError(nil, resp)
	}
	return resp, nil
}

// doPost is a helper for making POST requests and handling common errors.
func (c *Client) doPost(ctx context.Context, url string, body []byte) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, &APIRequestError{Op: "new_request", URL: url, Err: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &APIRequestError{Op: "do", URL: url, Err: err}
	}

	if isJSON, body := isJSONResponse(resp); !isJSON {
		defer resp.Body.Close()
		return nil, &APIRequestError{
			Op:  "doPost",
			URL: url,
			Err: fmt.Errorf("non-JSON response (%s): %s", resp.Header.Get("Content-Type"), body),
		}
	}

	// Accept 200 or 201 as success for POST
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		return nil, DecodeAPIError(body, resp)
	}
	return resp, nil
}

// doPatch is a helper for making PATCH requests and handling common errors.
func (c *Client) doPatch(ctx context.Context, url string, body []byte) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(body))
	if err != nil {
		return nil, &APIRequestError{Op: "new_request", URL: url, Err: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &APIRequestError{Op: "do", URL: url, Err: err}
	}

	if isJSON, body := isJSONResponse(resp); !isJSON {
		defer resp.Body.Close()
		return nil, &APIRequestError{
			Op:  "doPatch",
			URL: url,
			Err: fmt.Errorf("non-JSON response (%s): %s", resp.Header.Get("Content-Type"), body),
		}
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, DecodeAPIError(body, resp)
	}
	return resp, nil
}

// doPut is a helper for making PUT requests and handling common errors.
func (c *Client) doPut(ctx context.Context, url string, body []byte) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return nil, &APIRequestError{Op: "new_request", URL: url, Err: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &APIRequestError{Op: "do", URL: url, Err: err}
	}

	if isJSON, body := isJSONResponse(resp); !isJSON {
		defer resp.Body.Close()
		return nil, &APIRequestError{
			Op:  "doPut",
			URL: url,
			Err: fmt.Errorf("non-JSON response (%s): %s", resp.Header.Get("Content-Type"), body),
		}
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, DecodeAPIError(body, resp)
	}
	return resp, nil
}

// doDelete is a helper for making DELETE requests and handling common errors.
func (c *Client) doDelete(ctx context.Context, url string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, &APIRequestError{Op: "new_request", URL: url, Err: err}
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, &APIRequestError{Op: "do", URL: url, Err: err}
	}

	if isJSON, body := isJSONResponse(resp); !isJSON {
		defer resp.Body.Close()
		return nil, &APIRequestError{
			Op:  "doDelete",
			URL: url,
			Err: fmt.Errorf("non-JSON response (%s): %s", resp.Header.Get("Content-Type"), body),
		}
	}

	// Accept 200 or 204 as success for DELETE
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		return nil, DecodeAPIError(nil, resp)
	}
	return resp, nil
}

// isJSONResponse returns true if the Content-Type is json or the body is empty
func isJSONResponse(resp *http.Response) (bool, string) {
	// Handle reuse of body stream
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	resp.Body = io.NopCloser(&buf)

	bodyBytes, _ := io.ReadAll(tee)
	// Body is empty, no need to check Content-Type
	if len(bodyBytes) == 0 {
		return true, ""
	}

	contentType := resp.Header.Get("Content-Type")
	if len(contentType) < len("application/json") || !strings.Contains(contentType, "application/json") {
		return false, ""
	}

	return true, ""
}
