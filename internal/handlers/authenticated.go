package handlers

import (
	"back-end/logs"
	"database/sql"
	"net/http"
)

func (env *Environment) isAuthenticated(r *http.Request) (bool, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil
		}
		return false, err
	}

	var userId int
	err = env.DB.QueryRow("SELECT user_id FROM sessions WHERE token = $1", cookie.Value).Scan(&userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (env *Environment) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, err := env.isAuthenticated(r)
		if err != nil {
			logs.SendError(w, http.StatusBadRequest, "Internal server error")
			return
		}

		if !authenticated {
			logs.SendError(w, http.StatusBadRequest, "User not authenticated")
			return
		}

		next.ServeHTTP(w, r)
	}
}
