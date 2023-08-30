package handlers

import (
	"back-end/logs"
	"encoding/json"
	"net/http"
)

func (env *Environment) GetTestResultsHandler(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		logs.SendError(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	var userId int
	err = env.DB.QueryRow("SELECT user_id FROM sessions WHERE token = $1", cookie.Value).Scan(&userId)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Invalid session token")
		return
	}

	testId := r.URL.Query().Get("test_id")
	if testId == "" {
		logs.SendError(w, http.StatusBadRequest, "Missing test_id")
		return
	}

	var correctAnswersCount int
	err = env.DB.QueryRow("SELECT COUNT(*) FROM user_answers WHERE user_id = $1 AND test_id = $2 AND is_correct = true", userId, testId).Scan(&correctAnswersCount)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Error retrieving correct answers count")
		return
	}

	var totalAnswersCount int
	err = env.DB.QueryRow("SELECT COUNT(*) FROM user_answers WHERE user_id = $1 AND test_id = $2", userId, testId).Scan(&totalAnswersCount)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Error retrieving total answers count")
		return
	}

	if totalAnswersCount == 0 {
		logs.SendError(w, http.StatusBadRequest, "Total answers count cannot be zero")
		return
	}

	percentageCorrect := (float64(correctAnswersCount) / float64(totalAnswersCount)) * 100

	result := map[string]float64{
		"percentage_correct": percentageCorrect,
	}
	json.NewEncoder(w).Encode(result)
}
