package internal

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

var logger *log.Logger

func Log(msg string, args ...interface{}) {
	if logger != nil {
		logger.Printf(msg+"\n", args...)
	}
}

//type M map[string]interface{} // Just a simple shortcut

////////////////////////////////////////////////////////////////////////////////

type WebError struct {
	Code int
	Err  error
}

func (e WebError) Error() string {
	return e.Err.Error()
}

func (e WebError) Status() int {
	return e.Code
}

type H func(w http.ResponseWriter, r *http.Request) error

func (s *Server) Handler(h H) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // default
		err := h(w, r)
		if err != nil {
			switch e := err.(type) {
			case WebError:
				http.Error(w, e.Error(), e.Status())
			default:
				Log("server error: %s", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	})
}

func webClientAddr(r *http.Request) string {
	addr := r.Header.Get("x-forwarded-for") // Possible proxy addr
	if addr == "" {
		addr = r.RemoteAddr // fallback
	}
	return addr
}

////////////////////////////////////////////////////////////////////////////////

const (
	cookieName   = "session"
	cookieExpire = 24 * time.Hour
)

var (
	hashKey       = securecookie.GenerateRandomKey(64)
	blockKey      = securecookie.GenerateRandomKey(32)
	cookieEncoder = securecookie.New(hashKey, blockKey).MaxAge(int(cookieExpire.Seconds()))
)

func defaultCookie() *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode, // Prevents CSRF
		HttpOnly: true,
		//Secure:   true, // NOTE: only works when we run the site behind TLS

		// NOTE: don't forget to set these values on the new cookie
		//Value:    token,
		//Expires:  time.Now().Add(expire),
		//MaxAge:   int(expire.Seconds()),
	}
}

func setCookie(w http.ResponseWriter, data string) {
	cookie := defaultCookie()
	if data != "" {
		val, err := cookieEncoder.Encode(cookieName, data)
		if err != nil {
			// silently ignore errors, no cookies will be set. sue me
			// this should never happen with a string?
			return
		}
		cookie.Value = val
		cookie.Expires = time.Now().Add(cookieExpire)
		cookie.MaxAge = int(cookieExpire.Seconds())
	}
	http.SetCookie(w, cookie)
}

func getCookie(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	var data string
	err = cookieEncoder.Decode(cookieName, cookie.Value, &data)
	if err != nil {
		return ""
	}
	return data
}
