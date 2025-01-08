package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenFromHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh tokenneeded in the request headers", err)
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenFromHeader)
	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || (refreshToken.RevokedAt.Valid && !refreshToken.RevokedAt.Time.IsZero()) {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Duration(60*60)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't create a new token for this user", err)
		return
	}
	type Res struct {
		Token string `json:"token"`
	}
	data := Res{
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenFromHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh tokenneeded in the request headers", err)
		return
	}
	err = cfg.db.RevokeToken(r.Context(), tokenFromHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if len(params.Email) < 5 || len(params.Password) < 3 {
		respondWithError(w, http.StatusUnauthorized, "Email or Password failed validation", fmt.Errorf("email or password failed validation"))
		return
	}
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode the request", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user with provided email", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(60*60)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	generatedRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating a refresh token", err)
		return
	}

	futureTime := time.Now().AddDate(0, 0, 60)
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     generatedRefreshToken,
		UserID:    user.ID,
		ExpiresAt: futureTime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create a refresh token", err)
		return
	}

	data := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if len(params.Email) < 5 || len(params.Password) < 3 {
		respondWithError(w, http.StatusInternalServerError, "Email or Password failed validation", fmt.Errorf("email or password failed validation"))
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	pw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't secure password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: pw,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a user", err)
		return
	}

	data := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, data)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if len(params.Email) < 5 || len(params.Password) < 3 {
		respondWithError(w, http.StatusInternalServerError, "Email or Password failed validation", fmt.Errorf("email or password failed validation"))
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	pw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't secure password", err)
		return
	}
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: pw,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a user", err)
		return
	}

	data := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	headerKey, err := auth.GetPolkaApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "coudln't find api key", err)
		return
	}
	if headerKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "mismatched keys", err)
		return
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't upgrade user", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = cfg.db.UpgradeUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
