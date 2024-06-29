package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitHeader := strings.Split(authHeader, " ")
	authToken := splitHeader[1]
	claims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token")
		return
	}

	userIdString, err := jwtToken.Claims.GetSubject()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue parsing auth token")
		return
	}

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid User Id")
		return
	}

	user, err := cfg.DB.GetUserById(userId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user found")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err = cfg.DB.UpdateUser(userId, params.Email, params.Password, user.RefreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	responseWithJSON(w, http.StatusOK, User{
		Id:    user.Id,
		Email: user.Email,
	})

}
