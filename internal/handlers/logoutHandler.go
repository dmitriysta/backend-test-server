package handlers

import (
	"back-end/logs"
	"net/http"
	"time"
)

func (env *Environment) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		logs.SendError(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		logs.SendError(w, http.StatusBadRequest, "No session token found")
		return
	}

	_, err = env.DB.Exec("DELETE FROM sessions WHERE token = $1", cookie.Value)
	if err != nil {
		logs.SendError(w, http.StatusBadRequest, "Error during logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	w.WriteHeader(http.StatusOK)
	logs.SendJSONResponse(w, http.StatusOK, "Logged out successfully")
}
