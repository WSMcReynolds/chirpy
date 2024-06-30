package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) webhookRequestHandler(w http.ResponseWriter, r *http.Request) {

	if len(r.Header.Values("Authorization")) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")
	splitHeader := strings.Split(authHeader, " ")
	apiToken := splitHeader[1]

	if cfg.PolkaKey != apiToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue decoding body")
		return
	}

	if params.Event == "user.upgraded" {
		user, err := cfg.DB.GetUserById(params.Data.UserId)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = cfg.DB.UpdateUser(user.Id, user.Email, user.Password, user.RefreshToken, true, false)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "issue updating user")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
