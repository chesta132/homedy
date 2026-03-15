package authlib

import (
	"errors"
	"fmt"
	"homedy/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     string     `json:"user_id"`
	RotateAt   *time.Time `json:"rotate_at"`
	RememberMe bool       `json:"remember_me"`
	jwt.RegisteredClaims
}

func createClaim(id string, rememberMe bool, expiresAt time.Duration, rotateAt *time.Duration) *jwt.Token {
	var rotate *time.Time

	if rotateAt != nil {
		r := time.Now().Add(*rotateAt)
		rotate = &r

		if !rememberMe {
			// max 1 day expiry for not remember me refresh token
			maxSessionExpiry := time.Hour * 24
			if expiresAt > maxSessionExpiry {
				expiresAt = maxSessionExpiry
			}
		}
	}

	return jwt.NewWithClaims(config.SIGN_METHOD, Claims{
		UserID:     id,
		RotateAt:   rotate,
		RememberMe: rememberMe,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.APP_NAME,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	})
}

func createKeyFunc(secret string) jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		if t.Method != config.SIGN_METHOD {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secret), nil
	}
}

func CreateRefreshToken(id string, rememberMe bool) string {
	token := createClaim(id, rememberMe, config.REFRESH_TOKEN_EXPIRY, &config.ROTATE_REFRESH_TOKEN_AFTER)
	str, _ := token.SignedString([]byte(config.REFRESH_SECRET))
	return str
}

func CreateAccessToken(id string, rememberMe bool) string {
	token := createClaim(id, rememberMe, config.ACCESS_TOKEN_EXPIRY, nil)
	str, _ := token.SignedString([]byte(config.ACCESS_SECRET))
	return str
}

func ParseRefreshToken(str string) (claims Claims, err error) {
	token, err := jwt.ParseWithClaims(str, &claims, createKeyFunc(config.REFRESH_SECRET))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = fmt.Errorf("%w: refresh token", ErrTokenExpired)
		}
		return
	}
	if !token.Valid {
		err = fmt.Errorf("%w: refresh token", ErrInvalidToken)
	}
	return
}

func ParseAccessToken(str string) (claims Claims, err error) {
	token, err := jwt.ParseWithClaims(str, &claims, createKeyFunc(config.ACCESS_SECRET))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = fmt.Errorf("%w: access token", ErrTokenExpired)
		}
		return
	}
	if !token.Valid {
		err = fmt.Errorf("%w: access token", ErrInvalidToken)
	}
	return
}
