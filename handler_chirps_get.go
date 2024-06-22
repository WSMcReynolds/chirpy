package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) chirpsGetHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chrip Id")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chrip")
		return
	}

	responseWithJSON(w, http.StatusOK, Chirp{
		Id:   dbChirp.Id,
		Body: dbChirp.Body,
	})

}

func (cfg *apiConfig) chirpsRetrieveHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chrips")
		return
	}

	chirps := []Chirp{}

	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Id:   chirp.Id,
			Body: chirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	responseWithJSON(w, http.StatusOK, chirps)
}
