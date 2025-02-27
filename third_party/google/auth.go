package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/idtoken"
	"io/ioutil"
	"net/http"
	"time"
)

type Claims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

type AuthResponse struct {
	Name    string `json:"name"`
	UserId  string `json:"sub"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func Auth(ctx context.Context, clientId string, code string) (*AuthResponse, error) {
	payload, err := idtoken.Validate(ctx, code, clientId)
	if err != nil {
		return nil, err
	}

	resp := &AuthResponse{
		Name:   payload.Claims["name"].(string),
		UserId: payload.Claims["sub"].(string),
	}
	if payload.Claims["picture"] != nil {
		resp.Picture = payload.Claims["picture"].(string)
	}
	if payload.Claims["email"] != nil {
		resp.Email = payload.Claims["email"].(string)
	}
	return resp, nil
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

func ParseCredentials(credentials string, webClientId string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(credentials, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return nil, errors.New("iss is invalid")
	}

	if claims.Audience != webClientId {
		return nil, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("JWT is expired")
	}
	return claims, nil
}
