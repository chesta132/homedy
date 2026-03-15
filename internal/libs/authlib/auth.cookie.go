package authlib

import (
	"homedy/config"
	"net/http"
	"time"
)

// set expires < 0 to invalidate cookie
func ToCookie(name string, str string, expires time.Duration) (cookie http.Cookie) {
	cookie = http.Cookie{
		Name:     name,
		Value:    str,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   config.IsEnvProd(),
		HttpOnly: true,
	}
	if expires > 0 {
		cookie.Expires = time.Now().Add(expires)
	} else if expires < 0 {
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
	}
	return
}

func Invalidate(name string) http.Cookie {
	return ToCookie(name, "", -1)
}

func CreateRefreshCookie(id string, rememberMe bool) http.Cookie {
	expires := time.Duration(0)
	str := CreateRefreshToken(id, rememberMe)
	if rememberMe {
		expires = config.REFRESH_TOKEN_EXPIRY
	}

	return ToCookie(config.REFRESH_TOKEN_KEY, str, expires)
}

func CreateAccessCookie(id string, rememberMe bool) http.Cookie {
	expires := time.Duration(0)
	str := CreateAccessToken(id, rememberMe)
	if rememberMe {
		expires = config.ACCESS_TOKEN_EXPIRY
	}

	return ToCookie(config.ACCESS_TOKEN_KEY, str, expires)
}

func CreateTokenCookie(id string, rememberMe bool) []http.Cookie {
	return []http.Cookie{CreateRefreshCookie(id, rememberMe), CreateAccessCookie(id, rememberMe)}
}

func InvalidateTokenCookie() []http.Cookie {
	return []http.Cookie{Invalidate(config.ACCESS_TOKEN_KEY), Invalidate(config.REFRESH_TOKEN_KEY)}
}
