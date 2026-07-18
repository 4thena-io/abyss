package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignVerify_Roundtrip(t *testing.T) {
	token, err := Sign(42, "alice", true, "https://example.com/avatar.png", "secret")
	if err != nil {
		t.Fatalf("unexpected error signing: %v", err)
	}

	claims, err := Verify(token, "secret")
	if err != nil {
		t.Fatalf("unexpected error verifying: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "alice" || !claims.IsAdmin || claims.AvatarURL != "https://example.com/avatar.png" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Fatalf("expected a future expiry, got %v", claims.ExpiresAt)
	}
}

func TestVerify_WrongSecret(t *testing.T) {
	token, err := Sign(1, "bob", false, "", "secret-a")
	if err != nil {
		t.Fatalf("unexpected error signing: %v", err)
	}

	if _, err := Verify(token, "secret-b"); err == nil {
		t.Fatal("expected verification to fail with the wrong secret")
	}
}

func TestVerify_TamperedToken(t *testing.T) {
	token, err := Sign(1, "bob", false, "", "secret")
	if err != nil {
		t.Fatalf("unexpected error signing: %v", err)
	}

	tampered := token[:len(token)-2] + "xy"
	if _, err := Verify(tampered, "secret"); err == nil {
		t.Fatal("expected verification to fail for a tampered token")
	}
}

func TestVerify_ExpiredToken(t *testing.T) {
	claims := Claims{
		UserID:   1,
		Username: "bob",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("unexpected error signing: %v", err)
	}

	if _, err := Verify(signed, "secret"); err == nil {
		t.Fatal("expected verification to fail for an expired token")
	}
}

func TestVerify_WrongSigningMethod(t *testing.T) {
	claims := Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("unexpected error signing: %v", err)
	}

	if _, err := Verify(signed, "secret"); err == nil {
		t.Fatal("expected verification to reject a non-HMAC signing method")
	}
}
