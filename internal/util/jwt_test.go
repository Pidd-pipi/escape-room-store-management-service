package util

import "testing"

func TestJWTRoundTrip(t *testing.T) {
	secret := "unit-test-secret"
	token, err := GenerateToken(secret, 7, "13800000000", "player", 1)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserID != 7 || claims.Role != "player" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if _, err := ParseToken(secret, token+"x"); err == nil {
		t.Fatalf("expected tampered token error")
	}
}
