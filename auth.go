package isbclient

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims represents the user information to embed in the JWT.
type UserClaims struct {
	DisplayName string   `json:"displayName"`
	UserName    string   `json:"userName"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
}

// Claims is the JWT claims structure for the API.
type Claims struct {
	User UserClaims `json:"user"`
	jwt.RegisteredClaims
}

// ISB roles
const (
	RoleAdmin   = "Admin"
	RoleManager = "Manager"
	RoleUser    = "User"
)

// GenerateJWT generates a JWT token string with the given user claims, secret, and expiry duration.
func GenerateJWT(user UserClaims, secret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// NewAdminUserClaims returns a UserClaims struct for an admin user with the given email.
func NewAdminUserClaims(email string) UserClaims {
	return UserClaims{
		DisplayName: "Admin",
		UserName:    email,
		Email:       email,
		Roles:       []string{RoleAdmin},
	}
}

// NewUserUserClaims returns a UserClaims struct for a regular user with the given email.
func NewUserUserClaims(email string) UserClaims {
	return UserClaims{
		DisplayName: "GitHub",
		UserName:    email,
		Email:       email,
		Roles:       []string{RoleUser},
	}
}
