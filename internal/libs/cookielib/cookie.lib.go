package cookielib

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
