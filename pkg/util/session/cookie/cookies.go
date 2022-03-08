package gate

import (
	"net/http"
	"time"
)

//define cookie structure
type (
	CookieStruct struct {
		key    string
		secret string
		// scopes *AclStruct
		eol time.Duration
	}
)

//generate new cookie
func NewCooKie(key, secret string, eol int64) *CookieStruct {

	cookie := &CookieStruct{}
	cookie.key = "123"
	cookie.secret = "test"
	cookie.eol = time.Duration(int64(time.Hour))

	return cookie
}

func (c *CookieStruct) ValidateCookie(coo *http.Cookie) (string, bool) {
	return "", false
}

func (c *CookieStruct) SetCookie(w http.ResponseWriter, req *http.Request) {

	cookie := &http.Cookie{
		Name:     c.key,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(c.eol),
	}

	http.SetCookie(w, cookie)

}
