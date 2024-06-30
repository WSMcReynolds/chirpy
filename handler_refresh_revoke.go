package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) refreshRevokeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitHeader := strings.Split(authHeader, " ")
	authToken := splitHeader[1]

	user, err := cfg.DB.GetUserByToken(authToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue getting refresh token")
		return
	}

	_, err = cfg.DB.UpdateUser(user.Id, user.Email, user.Password, "", user.IsChirpyRed, false)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue revoking refresh token")
	}

	w.WriteHeader(204)

}
