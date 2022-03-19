package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = []byte("my_secret_key")

var users = map[string]string{
	"evgsol":   "kek",
	"solarka":  "lol",
	"psmirnov": "arbidol",
	"miracle":  "lavanda",
}

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

const LIFETIME = 12 * time.Hour

func login(w http.ResponseWriter, request *http.Request) {
	var creds Creds
	if err := json.NewDecoder(request.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expected, ok := users[creds.Login]
	if !ok || expected != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(LIFETIME)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Login: creds.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}).SignedString(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})
}

func loginRequired(f func(http.ResponseWriter, *http.Request, string)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var claims Claims
		token, err := jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		f(w, r, claims.Login)
	})
}
