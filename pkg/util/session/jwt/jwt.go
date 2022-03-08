package jwt

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type contextKey struct {
	name string
}

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

type (
	JWTStruct struct {
		SignKey    interface{}
		verifyKey  interface{}
		SignMethod jwt.SigningMethod
		parser     *jwt.Parser
	}
)

func NewAuth() *JWTStruct {

	key := os.Getenv("SCRT_KET")

	return &JWTStruct{
		SignKey:    []byte(key),
		SignMethod: jwt.GetSigningMethod("HS256"),
		parser:     &jwt.Parser{},
	}
}

func (a *JWTStruct) AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		header := GetTokenFromHeader(r)
		validate, err := a.ValidateKey(header)
		if err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(401), 401)
		} else {
			ctx = NewContext(ctx, validate, err)
			if validate == nil || !validate.Valid {
				http.Error(w, http.StatusText(401), 401)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}
func (a *JWTStruct) ValidateKey(tokenn string) (*jwt.Token, error) {

	if tokenn == "" {
		return nil, errors.New("Token not found")
	}

	// Verify the token
	token, err := a.Decode(tokenn)

	if err != nil {

		verr := err.(*jwt.ValidationError)

		if verr != nil {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				errors.New("Token expired")
			} else if verr.Errors&jwt.ValidationErrorIssuedAt > 0 {
				errors.New("xtau sapa bagi nih")
			} else if verr.Errors&jwt.ValidationErrorIssuer > 0 {
				errors.New("xtau sapa bagi nih")
			}
		}

		return token, err
	}

	if token == nil || !token.Valid {
		err = errors.New("Token not valid")

		return token, err
	}

	// Verify signing algorithm
	if token.Method != a.SignMethod {
		return token, errors.New("Sing method not calid")
	}

	// Valid!
	return token, nil

}

func (a *JWTStruct) Decode(tokenString string) (t *jwt.Token, err error) {

	t, err = a.parser.Parse(tokenString, a.keyFunc)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (a *JWTStruct) keyFunc(t *jwt.Token) (interface{}, error) {
	return a.SignKey, nil
}
