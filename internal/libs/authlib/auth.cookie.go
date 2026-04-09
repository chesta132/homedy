package authlib

import (
	"homedy/config"
	"homedy/internal/libs/cookielib"
	"net/http"
	"time"
)

func Invalidate(name string) http.Cookie {
	return cookielib.ToCookie(name, "", -1)
}

func CreateRefreshCookie(id string, rememberMe bool) http.Cookie {
	expires := time.Duration(0)
	str := CreateRefreshToken(id, rememberMe)
	if rememberMe {
		expires = config.REFRESH_TOKEN_EXPIRY
	}

	return cookielib.ToCookie(config.REFRESH_TOKEN_KEY, str, expires)
}

func CreateAccessCookie(id string, rememberMe bool) http.Cookie {
	expires := time.Duration(0)
	str := CreateAccessToken(id, rememberMe)
	if rememberMe {
		expires = config.ACCESS_TOKEN_EXPIRY
	}

	return cookielib.ToCookie(config.ACCESS_TOKEN_KEY, str, expires)
}

func CreateTokenCookie(id string, rememberMe bool) []http.Cookie {
	return []http.Cookie{CreateRefreshCookie(id, rememberMe), CreateAccessCookie(id, rememberMe)}
}

func InvalidateTokenCookie() []http.Cookie {
	return []http.Cookie{Invalidate(config.ACCESS_TOKEN_KEY), Invalidate(config.REFRESH_TOKEN_KEY)}
}
