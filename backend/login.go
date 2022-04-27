package main

import (
	"encoding/json"
	"errors"
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

type LoginManager struct {
	userManager *UserManager
}

func NewLoginManager(
	userManager *UserManager,
) *LoginManager {
	return &LoginManager{
		userManager: userManager,
	}
}

func (m *LoginManager) Login(w http.ResponseWriter, request *http.Request) {
	var creds Creds
	if err := json.NewDecoder(request.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := m.userManager.Check(request.Context(), creds.Login, creds.Password); err != nil {
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

type RegisterRequest struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (m *LoginManager) Register(request *http.Request) (interface{}, error) {
	var input RegisterRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		return nil, err
	}

	err := m.userManager.Create(request.Context(), input.Login, input.Password, input.Email)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func loginRequired(f func(*http.Request, string) (interface{}, error)) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				return nil, errors.New("no auth token cookie")
			}
			return nil, errors.New("failed to get auth token cookie")
		}

		var claims Claims
		token, err := jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return nil, errors.New("invalid auth signature")
			}
			return nil, errors.New("failed to check auth signature")
		}
		if !token.Valid {
			return nil, errors.New("invalid auth token")
		}

		return f(r, claims.Login)
	}
}
