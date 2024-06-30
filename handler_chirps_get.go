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
		Id:       dbChirp.Id,
		AuthorId: dbChirp.AuthorId,
		Body:     dbChirp.Body,
	})

}

func (cfg *apiConfig) chirpsRetrieveHandler(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.DB.GetChirps()

	chirps := []Chirp{}
	authorId := 0

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chrips")
		return
	}

	authorIdString := r.URL.Query().Get("author_id")
	if authorIdString != "" {
		authorId, err = strconv.Atoi(authorIdString)
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid author id")
	}

	if authorId > 0 {
		for _, chirp := range dbChirps {
			if chirp.AuthorId == authorId {
				chirps = append(chirps, Chirp{
					Id:       chirp.Id,
					AuthorId: chirp.AuthorId,
					Body:     chirp.Body,
				})
			}
		}
	} else {
		for _, chirp := range dbChirps {
			chirps = append(chirps, Chirp{
				Id:       chirp.Id,
				AuthorId: chirp.AuthorId,
				Body:     chirp.Body,
			})
		}
	}

	sortType := r.URL.Query().Get("sort")

	if sortType == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})
	}

	responseWithJSON(w, http.StatusOK, chirps)
}
