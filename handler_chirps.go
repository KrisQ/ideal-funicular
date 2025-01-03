package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	v := r.PathValue("chirpId")
	chirpId, err := uuid.Parse(v)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coulnd't parse id", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find chrip", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find chirps", err)
		return
	}
	chirps := make([]Chirp, len(dbChirps))
	for i, chirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)

}
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	chirp, err := validateChirp(params.Body)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Invalid chirp", err)
		return
	}

	storedChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirp,
		UserID: params.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	data := Chirp{
		ID:        storedChirp.ID,
		CreatedAt: storedChirp.CreatedAt,
		UpdatedAt: storedChirp.UpdatedAt,
		Body:      storedChirp.Body,
		UserId:    storedChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, data)

}

func validateChirp(chirp string) (string, error) {

	const maxChirpLength = 140
	if len(chirp) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	profanities := make(map[string]interface{})
	profanities["kerfuffle"] = nil
	profanities["sharbert"] = nil
	profanities["fornax"] = nil

	words := strings.Split(chirp, " ")
	cleanedBody := make([]string, len(words))

	for i, word := range words {
		if _, ok := profanities[strings.ToLower(word)]; ok {
			cleanedBody[i] = "****"
		} else {
			cleanedBody[i] = word
		}
	}

	return strings.Join(cleanedBody, " "), nil
}
