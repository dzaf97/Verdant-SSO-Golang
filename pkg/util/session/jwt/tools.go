package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/dgrijalva/jwt-go"
)

type (
	MergeClaims struct {
		UserID    string
		FirstName string
		Email     string
		Roles     string

		// Namespace   string   `json:"namespace"`
		// AccessList []string `json:"namespace"`
		jwt.StandardClaims
	}
)

func GetTokenFromHeader(r *http.Request) string {

	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

func GetTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value

}

// //generate new jwt
func GenNewToken(userid, name, email, role string) (string, error) {

	standard := jwt.StandardClaims{

		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "verdant-SSO",
		IssuedAt:  time.Now().Unix(),
	}

	merge := &MergeClaims{
		userid,
		name,
		email,
		role,
		standard,
	}

	ji := jwt.NewWithClaims(jwt.SigningMethodHS256, merge)
	key := os.Getenv("SCRT_KET")

	tokenstring, err := ji.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return tokenstring, nil

}

//create new context
func NewContext(ctx context.Context, t *jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func Get(ctx context.Context) (*simplejson.Json, error) {
	token, _ := ctx.Value(TokenCtxKey).(*jwt.Token)

	var claims jwt.MapClaims
	if token != nil {
		if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
			claims = tokenClaims
		} else {
			panic(fmt.Sprintf("jwtauth: unknown type of Claims: %T", token.Claims))
		}
	} else {
		claims = jwt.MapClaims{}
	}

	err, _ := ctx.Value(ErrorCtxKey).(error)

	b, err := json.Marshal(claims)

	if err != nil {
		fmt.Println("error:", err)
	}

	data, err := simplejson.NewJson(b)
	if err != nil {
		return nil, err
	}

	return data, err
}
