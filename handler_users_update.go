package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.confirmUserAuthenticaiton(r)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
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

	_, err = cfg.DB.UpdateUser(user.Id, params.Email, params.Password, user.RefreshToken, user.IsChirpyRed, true)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	responseWithJSON(w, http.StatusOK, User{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
