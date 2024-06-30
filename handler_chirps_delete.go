package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) chirpsDeleteHandler(w http.ResponseWriter, r *http.Request) {

	user, err := cfg.confirmUserAuthenticaiton(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpId, err := strconv.Atoi(chirpIDString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chrip id")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp Id invalid")
		return
	}

	if user.Id != chirp.AuthorId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.DB.DeleteChirp(chirpId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue deleting chrip")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
