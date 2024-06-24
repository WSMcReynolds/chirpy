package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

func (cfg *apiConfig) usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	existingUser, _ := cfg.DB.GetUserByEmail(params.Email)

	if existingUser.Email == params.Email {
		respondWithError(w, http.StatusBadRequest, "User with email already exists")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	responseWithJSON(w, http.StatusCreated, User{
		Id:    user.Id,
		Email: user.Email,
	})

}
