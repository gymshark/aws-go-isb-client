package isbclient

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewAdminUserClaims(t *testing.T) {
	email := "admin@example.com"
	claims := NewAdminUserClaims(email)
	if claims.DisplayName != "Admin" {
		t.Errorf("expected DisplayName 'Admin', got '%s'", claims.DisplayName)
	}
	if claims.UserName != email {
		t.Errorf("expected UserName '%s', got '%s'", email, claims.UserName)
	}
	if claims.Email != email {
		t.Errorf("expected Email '%s', got '%s'", email, claims.Email)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != RoleAdmin {
		t.Errorf("expected Roles ['%s'], got %v", RoleAdmin, claims.Roles)
	}
}

func TestGenerateJWT(t *testing.T) {
	secret := "testsecret"
	email := "admin@example.com"
	user := NewAdminUserClaims(email)
	expiresIn := time.Hour
	tokenStr, err := GenerateJWT(user, secret, expiresIn)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	if strings.Count(tokenStr, ".") != 2 {
		t.Errorf("expected JWT to have 2 dots, got: %s", tokenStr)
	}

	// Parse and validate
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !token.Valid {
		t.Error("token is not valid")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatal("claims type assertion failed")
	}
	if claims.User.Email != email {
		t.Errorf("expected email '%s', got '%s'", email, claims.User.Email)
	}
	if claims.User.DisplayName != "Admin" {
		t.Errorf("expected DisplayName 'Admin', got '%s'", claims.User.DisplayName)
	}
	if len(claims.User.Roles) != 1 || claims.User.Roles[0] != RoleAdmin {
		t.Errorf("expected Roles ['%s'], got %v", RoleAdmin, claims.User.Roles)
	}
}
