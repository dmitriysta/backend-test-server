package handlers

import (
	"back-end/logs"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type Question struct {
	ID             int    `json:"id"`
	Text           string `json:"text"`
	VariantID      int    `json:"variant_id"`
	QuestionNumber int    `json:"question_number"`
}

func (env *Environment) StartTestHandler(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		logs.SendError(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	variantID := r.URL.Query().Get("variant_id")
	questionNumber := r.URL.Query().Get("question_number")
	if variantID == "" || questionNumber == "" {
		logs.SendError(w, http.StatusBadRequest, "Missing variant_id or question_number")
		return
	}

	qNum, err := strconv.Atoi(questionNumber)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Invalid question number format")
		return
	}

	var question Question
	err = env.DB.QueryRow("SELECT id, question_text, variant_id, question_number FROM test_questions WHERE variant_id = $1 AND question_number = $2", variantID, qNum).Scan(&question.ID, &question.Text, &question.VariantID, &question.QuestionNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			logs.SendError(w, http.StatusBadRequest, "No questions found for this variant and question number")
		} else {
			logs.SendError(w, http.StatusBadRequest, "Internal server error")
		}
		return
	}

	json.NewEncoder(w).Encode(question)
}
