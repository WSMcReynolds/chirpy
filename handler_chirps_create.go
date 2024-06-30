package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

func (cfg *apiConfig) chirpsCreateHandler(w http.ResponseWriter, r *http.Request) {

	user, err := cfg.confirmUserAuthenticaiton(r)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleanChirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleanChirp, user.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	responseWithJSON(w, http.StatusCreated, Chirp{
		Id:       chirp.Id,
		AuthorId: user.Id,
		Body:     chirp.Body,
	})
}

func validateChirp(chirp string) (string, error) {
	if len(chirp) >= 140 {
		return "", errors.New("Chirp is too long")
	}

	cleanChirp := sanitizeChrip(chirp)
	return cleanChirp, nil
}

func sanitizeChrip(chirp string) string {
	badWords := make(map[string]bool)
	badWords["fornax"] = true
	badWords["sharbert"] = true
	badWords["kerfuffle"] = true

	originalWords := strings.Split(chirp, " ")
	chirp = strings.ToLower(chirp)
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if badWords[word] {
			originalWords[i] = "****"
		}
	}

	cleanChrip := strings.Join(originalWords, " ")
	return cleanChrip
}
