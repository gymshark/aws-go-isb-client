# AWS Innovation Sandbox Go Client

> **AI Usage Statement:**
> Nearly all client logic and tests in this repository were generated using AI assistance.

## Overview

This Go package provides a client for interacting with the [AWS Innovation Sandbox](https://docs.aws.amazon.com/solutions/latest/innovation-sandbox-on-aws/solution-overview.html) API. It supports JWT-based authentication and simplifies making requests to the API.

## Installation

Add the module to your project:

```sh
go get github.com/gymshark/aws-go-isb-client
```

## Authentication (JWT)

The client uses JWT bearer tokens for authentication. You can generate a JWT for an admin user using the provided helper:

```go
import (
    "time"
	
    "github.com/gymshark/aws-go-isb-client"
)

user := isbclient.NewAdminUserClaims("admin@gymshark.com")
secret := "your-shared-secret" // Secret value from the CloudFormation stack output `JwtSecretArn`
expiresIn := 2*time.Hour

jwtToken, err := isbclient.GenerateJWT(user, secret, expiresIn)
if err != nil {
    // handle error
}
```

> The secret is the value stored in the secret referenced by the CloudFormation stack output `JwtSecretArn`.

## Initialising the Client

Create a new client instance with the API base URL and your JWT token:

```go
client := isbclient.NewClient("https://<CloudFrontDistributionUrl>/api", jwtToken)
```

> `<CloudFrontDistributionUrl>` should be replaced with the `CloudFrontDistributionUrl` output from the CloudFormation compute stack.

## Making Requests

> **Note:** The following client methods are generated from the OpenAPI specification in `spec.yaml`. Refer to the spec for endpoint details and request/response structures.

### GetLeases

Fetch a paginated list of leases:

```go
resp, err := client.GetLeases(ctx, queryBuilder)
```

### GetLeaseByID

Fetch a lease by its ID:

```go
leaseReq := &isbclient.GetLeaseByIDRequest{LeaseID: "lease-id"}
resp, err := client.GetLeaseByID(ctx, leaseReq)
```

### CreateLease

Request a new lease:

```go
leaseReq := &isbclient.CreateLeaseRequest{
    LeaseTemplateUUID: "template-uuid",
    Comments:          "optional comment",
}
resp, err := client.CreateLease(ctx, leaseReq)
```

### CreateLeaseAsUser

Create a lease for another user. See [Acting on Behalf of Another User (Lease Creation)](#acting-on-behalf-of-another-user-lease-creation) for details and usage:

```go
resp, err := client.CreateLeaseAsUser(ctx, leaseReq, "target.user@gymshark.com", jwtSecret)
```

### GetLeaseTemplates

Fetch available lease templates:

```go
resp, err := client.GetLeaseTemplates(ctx, queryBuilder)
```

### FetchAllLeases

Fetch all leases using pagination:

```go
resp, err := client.FetchAllLeases(ctx, getLeasesReq)
```

### FetchAllLeaseTemplates

Fetch all lease templates using pagination:

```go
resp, err := client.FetchAllLeaseTemplates(ctx, getLeaseTemplatesReq)
```

### GetAccounts

Fetch a paginated list of accounts:

```go
resp, err := client.GetAccounts(ctx, queryBuilder)
```

### FetchAllAccounts

Fetch all accounts using pagination:

```go
resp, err := client.FetchAllAccounts(ctx, getAccountsReq)
```

Refer to the source code for available methods and request/response types.

## Acting on Behalf of Another User (Lease Creation)

To create a lease for another user, use the `CreateLeaseAsUser` method. 

> This method generates a JWT for the target user using the `NewUserUserClaims` helper and makes the request with that token. These methods are generated from the OpenAPI specification in `spec.yaml`.

```go
userEmail := "target.user@gymshark.com"
jwtSecret := "your-shared-secret"
leaseReq := &isbclient.CreateLeaseRequest{
    LeaseTemplateUUID: "template-uuid",
    Comments:          "Lease for automation",
}
resp, err := client.CreateLeaseAsUser(ctx, leaseReq, userEmail, jwtSecret)
if err != nil {
    // handle error
}
// process resp
```

This uses the `NewUserUserClaims` helper to generate the JWT for the specified user.

## Roles

Supported roles for JWT claims:
- `Admin`
- `Manager`
- `User`

## Dependencies

- [github.com/golang-jwt/jwt/v5](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)

## Testing

Run all tests:

```sh
go test ./...
```
