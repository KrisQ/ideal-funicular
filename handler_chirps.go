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
