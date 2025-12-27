package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const defaultTokenTTL = 72 * time.Hour

// HashPassword hashes plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword compares hashed password with plaintext candidate.
func CheckPassword(hashed, candidate string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(candidate))
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTManager builds a JWT manager with a secret and optional ttl (fallback to defaultTokenTTL when ttl <= 0).
func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	if ttl <= 0 {
		ttl = defaultTokenTTL
	}
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

// GenerateToken signs a JWT with user_id claim and exp.
func (m *JWTManager) GenerateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(m.ttl).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken validates token and returns embedded userID.
func (m *JWTManager) ParseToken(tokenStr string) (uint64, error) {
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	uidVal, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("missing user_id in token")
	}

	switch v := uidVal.(type) {
	case float64:
		return uint64(v), nil
	case uint64:
		return v, nil
	case int64:
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("unexpected user_id type: %T", uidVal)
	}
}
