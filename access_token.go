package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) getAccessToken(userId, expires_in_seconds int) (string, error) {
	userIdString := strconv.Itoa(userId)

	if expires_in_seconds == 0 || expires_in_seconds > 86400 {
		expires_in_seconds = 86400
	}

	now := time.Now()
	expireTime := now.Add(time.Second * time.Duration(expires_in_seconds))

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expireTime),
		Subject:   userIdString,
	})

	signedToken, err := jwtToken.SignedString([]byte(cfg.JWTSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil

}
