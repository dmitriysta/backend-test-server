package main

import (
	"back-end/internal/database"
	"back-end/internal/handlers"
	"back-end/logs"
	"fmt"
	"net/http"
)

func main() {
	logs.SetupLogging()
	go logs.AutoRotateLogs()

	config := database.LoadConfig("./internal/config/config.json")

	db := database.ConnectToDataBase(config.Database)
	if db != nil {
		defer db.Close()
		fmt.Println("Successfully connected to the database!")
	}

	env := &handlers.Environment{DB: db}

	http.HandleFunc("/login", env.LoginHandler)
	http.HandleFunc("/start_test", env.AuthMiddleware(env.StartTestHandler))
	http.HandleFunc("/test_result", env.AuthMiddleware(env.GetTestResultsHandler))
	http.HandleFunc("/logout", env.LogoutHandler)
	http.ListenAndServe(config.Server.Port, nil)

}
