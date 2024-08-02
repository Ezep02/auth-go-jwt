package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	SecretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{SecretKey: secretKey}
}

func (maker *JWTMaker) CreateToken(id int64, email string, isAdmin bool, duration time.Duration) (string, *UserClaims, error) {

	claims, err := NewUserClaim(id, email, isAdmin, duration)

	if err != nil {
		return "", nil, err
	}

	// Crear un nuevo token con las claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token usando la clave secreta
	tokenStr, err := token.SignedString([]byte(maker.SecretKey))

	if err != nil {
		return "", nil, fmt.Errorf("error signing token %w", err)
	}

	return tokenStr, claims, nil

}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	// Analizar el token y verificar las claims
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: invalid signing method")
		}

		return []byte(maker.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	// Verificar las claims
	claim, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claim, nil
}
