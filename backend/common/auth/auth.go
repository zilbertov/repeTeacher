package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	RoleTutor   = "tutor"
	RoleStudent = "student"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type Claims struct {
	Role      string `json:"role"`
	TutorID   int64  `json:"tutor_id,omitempty"`
	StudentID int64  `json:"student_id,omitempty"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type contextKey string

const claimsContextKey contextKey = "jwt_claims"

func NewToken(secret string, claims Claims, ttl time.Duration) (string, Claims, error) {
	now := time.Now()
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(ttl).Unix()
	if err := validateClaims(claims); err != nil {
		return "", Claims{}, err
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", Claims{}, err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", Claims{}, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsigned := encodedHeader + "." + encodedClaims
	signature := sign(secret, unsigned)
	return unsigned + "." + signature, claims, nil
}

func ParseToken(secret string, token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrUnauthorized
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSignature := sign(secret, unsigned)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return Claims{}, ErrUnauthorized
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrUnauthorized
	}
	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, ErrUnauthorized
	}
	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return Claims{}, ErrUnauthorized
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrUnauthorized
	}
	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return Claims{}, ErrUnauthorized
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return Claims{}, ErrUnauthorized
	}
	if err := validateClaims(claims); err != nil {
		return Claims{}, ErrUnauthorized
	}
	return claims, nil
}

func Require(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		claims, err := ParseToken(secret, token)
		if err != nil {
			writeAuthError(w, http.StatusUnauthorized, ErrUnauthorized.Error())
			return
		}
		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(Claims)
	return claims, ok
}

func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func RequireTutor(ctx context.Context) (Claims, error) {
	claims, ok := FromContext(ctx)
	if !ok {
		return Claims{}, ErrUnauthorized
	}
	if claims.Role != RoleTutor {
		return Claims{}, ErrForbidden
	}
	return claims, nil
}

func Current(ctx context.Context) (Claims, error) {
	claims, ok := FromContext(ctx)
	if !ok {
		return Claims{}, ErrUnauthorized
	}
	return claims, nil
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func sign(secret string, data string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func validateClaims(claims Claims) error {
	switch claims.Role {
	case RoleTutor:
		if claims.TutorID == 0 {
			return ErrUnauthorized
		}
	case RoleStudent:
		if claims.TutorID == 0 || claims.StudentID == 0 {
			return ErrUnauthorized
		}
	default:
		return ErrUnauthorized
	}
	if claims.ExpiresAt == 0 {
		return ErrUnauthorized
	}
	return nil
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
