package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	jwtAccessKey  = []byte(os.Getenv("JWT_SECRET"))
	jwtRefreshKey = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

type Claims struct {
	Id      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Version uint32    `json:"version"`
	Type    string    `json:"type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(id uuid.UUID, email string, version uint32) (string, error) {
	claims := &Claims{
		Id:      id,
		Email:   email,
		Version: version,
		Type:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtAccessKey)
}

func GenerateRefreshToken(id uuid.UUID, email string, version uint32) (string, error) {
	claims := &Claims{
		Id:      id,
		Email:   email,
		Version: version,
		Type:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtRefreshKey)
}
