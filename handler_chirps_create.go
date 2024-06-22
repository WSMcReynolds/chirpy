package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) chirpsCreateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleanChirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleanChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	responseWithJSON(w, http.StatusCreated, Chirp{
		Id:   chirp.Id,
		Body: chirp.Body,
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
