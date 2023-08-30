package handlers

import (
	"back-end/logs"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Environment struct {
	DB *sql.DB
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (env *Environment) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		logs.SendError(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var storedPassword string
	err = env.DB.QueryRow("SELECT password_hash FROM users WHERE login = $1", loginReq.Login).Scan(&storedPassword)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "User not found")
		return
	}

	if !сheckPasswordHash(loginReq.Password, storedPassword) {
		logs.SendError(w, http.StatusBadRequest, "Incorrect password")
		return
	}

	token, err := generateSessionToken()
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = env.DB.Exec("INSERT INTO sessions (token, user_id) VALUES ($1, $2)", token, loginReq.Login)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Error creating session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: token,
		Path:  "/",
	})

	w.WriteHeader(http.StatusOK)
	logs.SendJSONResponse(w, http.StatusOK, "Login successful")
}

func generateSessionToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("error generating session token: %v", err)
	}

	return hex.EncodeToString(token), nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func сheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
