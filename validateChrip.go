package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
		ValidChrip  bool   `json:"valid"`
	}

	resBody := returnVals{}

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	chirp := sanitizeChrip(params.Body)
	resBody.CleanedBody = chirp

	if !validChirpLenght(chirp) {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	resBody.ValidChrip = true
	responseWithJSON(w, 200, resBody)

}

func validChirpLenght(chrip string) bool {
	return len(chrip) <= 140
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Error marshalling JSON")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
