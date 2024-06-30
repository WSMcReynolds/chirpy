package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) loginRequestHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expiration int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIdString := strconv.Itoa(user.Id)

	expiresAt := params.Expiration

	if expiresAt == 0 || expiresAt > 86400 {
		expiresAt = 86400
	}

	now := time.Now()
	expireTime := now.Add(time.Second * time.Duration(expiresAt))

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expireTime),
		Subject:   userIdString,
	})

	signedToken, err := jwtToken.SignedString([]byte(cfg.JWTSecret))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create access token")
		return
	}

	randHex := make([]byte, 32)

	_, err = rand.Read(randHex)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue creating refresh token")
		return
	}

	refreshToken := hex.EncodeToString(randHex)

	user, err = cfg.DB.UpdateUser(user.Id, user.Email, user.Password, refreshToken, user.IsChirpyRed, false)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue updating user")
		return
	}

	responseWithJSON(w, http.StatusOK, User{
		Id:           user.Id,
		Email:        user.Email,
		Token:        signedToken,
		RefreshToken: user.RefreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})

}
