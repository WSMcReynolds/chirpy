package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) confirmUserAuthenticaiton(r *http.Request) (User, error) {
	authHeader := r.Header.Get("Authorization")
	splitHeader := strings.Split(authHeader, " ")
	authToken := splitHeader[1]
	claims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return User{}, err
	}

	userIdString, err := jwtToken.Claims.GetSubject()

	if err != nil {
		return User{}, err
	}

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		return User{}, err
	}

	user, err := cfg.DB.GetUserById(userId)

	if err != nil {
		return User{}, err
	}

	return User{
		Id:           user.Id,
		Email:        user.Email,
		RefreshToken: user.RefreshToken,
	}, nil
}
