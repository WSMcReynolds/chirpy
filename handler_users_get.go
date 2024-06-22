package main

import (
	"net/http"
	"sort"
)

func (cfg *apiConfig) usersRetrieveHandler(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve Users")
		return
	}

	users := []User{}

	for _, user := range dbUsers {
		users = append(users, User{
			Id:    user.Id,
			Email: user.Email,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	responseWithJSON(w, http.StatusOK, users)
}
