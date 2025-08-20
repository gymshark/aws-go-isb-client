# AWS Innovation Sandbox Go Client

This is a Go client library for the AWS Innovation Sandbox API that provides JWT-based authentication and simplifies making requests to the API.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

- **CRITICAL**: Set timeouts of 60+ minutes for builds and 30+ minutes for test commands. NEVER CANCEL long-running operations.

### Bootstrap, Build, and Test the Repository:
- Verify Go installation: `go version` (requires Go 1.24.5+)
- Setup dependencies: `make setup` -- takes <1 second when cached, ~5 seconds on first run
- Run all tests: `make test` -- takes ~19 seconds total. NEVER CANCEL. Set timeout to 60+ minutes.
- Format code: `make fmt` -- takes <1 second. Runs gofmt on all Go files.
- Run linting: `go vet ./...` -- takes <1 second. No additional linters configured.

### Build Verification:
- Test compilation: `go build .` -- takes <1 second. No binary output (library only).
- Run fast tests only: `go test -short ./...` -- takes <1 second. Skips slow network timeout tests.
- Clear test cache: `go clean -testcache` -- use before `make test` to get accurate timing measurements.

### Manual Validation Scenarios:
ALWAYS test basic client functionality after making changes by running:
```bash
cd /tmp && cat > test_client.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/gymshark/aws-go-isb-client"
)

func main() {
	// Test JWT generation
	user := isbclient.NewAdminUserClaims("admin@example.com")
	secret := "test-secret"
	expiresIn := 2 * time.Hour
	
	jwtToken, err := isbclient.GenerateJWT(user, secret, expiresIn)
	if err != nil {
		log.Fatalf("Failed to generate JWT: %v", err)
	}
	fmt.Printf("JWT generated successfully: %s...\n", jwtToken[:50])
	
	// Test client creation
	client := isbclient.NewClient("https://example.com/api", jwtToken)
	fmt.Printf("Client created successfully with BaseURL: %s\n", client.BaseURL)
	
	// Test request building
	leaseReq := &isbclient.GetLeasesRequest{PageIdentifier: "next", PageSize: "20"}
	query := leaseReq.BuildQuery()
	fmt.Printf("Query params built: %v\n", query)
	fmt.Println("Basic client functionality test completed successfully!")
}
EOF
go run test_client.go
```

Expected output should show JWT generation, client creation, and query building without errors.

## Timing and Timeout Requirements

- **NEVER CANCEL**: Build and test operations may take significant time
- `make setup`: <1 second when dependencies cached, ~5 seconds on first run
- `make test`: ~19 seconds total (~18.7s actual test time). NEVER CANCEL. Set timeout to 60+ minutes.
- `make fmt`: <1 second
- `go vet ./...`: <1 second
- `go build .`: <1 second
- `go test -short ./...`: <1 second (skips slow tests)

**WARNING**: One test (`TestNonJSONResponses`) takes ~18 seconds due to network timeout testing. This is normal behavior. Use `go test -short ./...` to skip slow tests during development.

## Repository Structure and Navigation

### Key Files and Locations:
- `README.md` - Comprehensive API documentation and usage examples
- `Makefile` - Build automation (setup, test, fmt, update-spec targets)
- `go.mod`/`go.sum` - Go module definition and dependencies
- `spec.yaml` - OpenAPI specification (38KB) defining the AWS Innovation Sandbox API
- `client.go` - Main HTTP client implementation with all API methods
- `types.go` - Request/response types and data structures
- `auth.go` - JWT authentication helpers and user claims
- `errors.go` - Custom error types for API and client errors
- `*_test.go` - Comprehensive test suite covering all functionality

### Project Structure:
```
.
├── README.md           # Complete API documentation
├── Makefile           # Build system (4 targets: setup, test, fmt, update-spec)
├── go.mod/go.sum      # Go module (single dependency: jwt/v5)
├── spec.yaml          # OpenAPI spec (downloaded from AWS GitHub)
├── client.go          # Main client (~21k lines with all API methods)
├── types.go           # Request/response types (~13k lines)
├── auth.go            # JWT authentication (~1.6k lines)
├── errors.go          # Error handling (~10k lines)
└── *_test.go          # Test files (37+ test cases)
```

### Important Code Areas:
- **Client Methods**: All in `client.go` - GetLeases, CreateLease, GetLeaseTemplates, etc.
- **Request Types**: All in `types.go` - *Request structs with BuildQuery() methods
- **Response Types**: All in `types.go` - *Response structs matching API responses
- **JWT Helpers**: In `auth.go` - NewAdminUserClaims, NewUserUserClaims, GenerateJWT
- **Error Handling**: In `errors.go` - APIRequestError, JSONDecodingError, etc.

## Dependencies

- **Only dependency**: `github.com/golang-jwt/jwt/v5 v5.3.0`
- **Go version**: Requires Go 1.24.5+ (currently using Go 1.24.6)
- **No additional tools**: No linters, formatters, or build tools beyond standard Go toolchain

## Validation Steps

Always run these validation steps before completing changes:
1. `make setup` - Ensure dependencies are current
2. `make fmt` - Format code according to Go standards  
3. `go vet ./...` - Run static analysis
4. `make test` - Run full test suite (NEVER CANCEL, ~36 seconds)
5. Manual validation scenario above - Test basic functionality

## Common Development Tasks

### Making API Changes:
- **ALWAYS** check `spec.yaml` for API contract details
- **ALWAYS** update corresponding types in `types.go` if changing request/response structures
- **ALWAYS** add tests in appropriate `*_test.go` file
- **ALWAYS** run manual validation scenario after changes

### Key API Methods Available:
- Leases: `GetLeases`, `GetLeaseByID`, `CreateLease`, `CreateLeaseAsUser`, `UpdateLease`, `ReviewLease`, `FreezeLease`, `TerminateLease`
- Lease Templates: `GetLeaseTemplates`, `UpdateLeaseTemplate`, `DeleteLeaseTemplate`
- Accounts: `GetAccounts`, `RegisterAccount`, `RetryCleanup`, `EjectAccount`
- Utilities: `FetchAllLeases`, `FetchAllLeaseTemplates`, `FetchAllAccounts` (pagination helpers)

### JWT Authentication:
- Admin users: `isbclient.NewAdminUserClaims("admin@example.com")`
- Regular users: `isbclient.NewUserUserClaims("user@example.com")`
- Generate tokens: `isbclient.GenerateJWT(claims, secret, duration)`

## Known Issues

- **Makefile bug**: `make update-spec` fails because it references `pkg/isb/spec.yaml` but `spec.yaml` is in root directory. Do not rely on this command.
- **Slow test**: `TestNonJSONResponses` takes ~18 seconds due to intentional timeout testing. This is expected behavior.

## Architecture Notes

- **Library only**: This is a client library, not an executable application
- **HTTP client**: Uses standard `net/http` with custom auth transport for Bearer tokens
- **Error handling**: Comprehensive custom error types for different failure scenarios
- **Testing**: Uses `httptest.NewServer` for mocking HTTP responses in tests
- **Thread safety**: Client is safe for concurrent use

Always build and test your changes thoroughly using the validation steps above.