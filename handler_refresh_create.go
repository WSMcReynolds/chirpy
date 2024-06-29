package main

import (
	"net/http"
	"strings"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshCreateHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitHeader := strings.Split(authHeader, " ")
	authToken := splitHeader[1]

	user, err := cfg.DB.GetUserByToken(authToken)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := cfg.getAccessToken(user.Id, 3600)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Issue creating access token")
		return
	}

	responseWithJSON(w, http.StatusOK, TokenResponse{
		Token: accessToken,
	})

}
