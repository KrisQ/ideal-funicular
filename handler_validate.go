package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profanities := make(map[string]interface{})
	profanities["kerfuffle"] = nil
	profanities["sharbert"] = nil
	profanities["fornax"] = nil

	words := strings.Split(params.Body, " ")
	cleanedBody := make([]string, len(words))

	for i, word := range words {
		if _, ok := profanities[strings.ToLower(word)]; ok {
			cleanedBody[i] = "****"
		} else {
			cleanedBody[i] = word
		}
	}

	data := returnVals{
		CleanedBody: strings.Join(cleanedBody, " "),
	}

	respondWithJSON(w, http.StatusOK, data)
}
