package main

import (
	"errors"
	"fmt"
	"net/http"
	"rest-api/pkg/data"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtTokenExpiry = time.Minute * 15
const refreshTokenExpiry = time.Hour * 24

type TokenPair struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

func (app *application) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// checks the raw Authorization header value
	// don't return the cached response to the wrong client
	// eg: compares request "Bearer xyz" with cached "Bearer xyz"
	// this is just plain string matching
	w.Header().Add("Vary", "Authorization")

	// get authorization header
	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// check to see if we have word "Bearer"
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("unauthorized: no Bearer")
	}

	token := headerParts[1]

	// declare empty Claims
	claims := &Claims{}

	// parse the token with our claims (we read into claims), using our secret (from the receiver)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.JWTSecret), nil
	})

	// check for an error; note that this catches expired tokens as well
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	// make sure that we issued this token
	if claims.Issuer != app.Domain {
		return "", nil, errors.New("incorrect issuer")
	}

	// valid token
	return token, claims, nil
}

func (app *application) generateTokenPairs(user *data.User) (TokenPair, error) {
	// ** create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)

	// set claims
	// NOTE: This returns a pointer
	// So we are setting
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = app.Domain
	claims["iss"] = app.Domain

	if user.IsAdmin == 1 {
		claims["admin"] = true
	} else {
		claims["admin"] = false
	}

	// set expiry
	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()

	// create signed token
	signedAccessToken, err := accessToken.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return TokenPair{}, err
	}

	// ** create refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	// set claims
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)

	// set expiry; must be longer than jwt expiry
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// created signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return TokenPair{}, err
	}

	var tokenPairs = TokenPair{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokenPairs, nil
}
