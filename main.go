// This project is largely based on Eli Bendersky's series `REST Servers in Go`.
// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/wagnojunior/taskstore/internal/taskstore"
)

// `appConfig` configures all env variables required for this application to run
type appConfig struct {
	PSQL   taskstore.PostgresConfig
	Server struct {
		Address string
	}
}

// `loadEnvVar` loads the required environment variables
func loadEnvVar() (appConfig, error) {
	var appCfg appConfig

	err := godotenv.Load()
	if err != nil {
		return appCfg, err
	}

	// /////////////////////////////////////////////////////////////////////////
	// DATABASE
	// /////////////////////////////////////////////////////////////////////////
	appCfg.PSQL.Host = os.Getenv("DB_HOST")
	appCfg.PSQL.Port = os.Getenv("DB_PORT")
	appCfg.PSQL.User = os.Getenv("DB_USER")
	appCfg.PSQL.Password = os.Getenv("DB_PASSWORD")
	appCfg.PSQL.Database = os.Getenv("DB_DATABASE")
	appCfg.PSQL.SSLMode = os.Getenv("DB_SSLMODE")

	// /////////////////////////////////////////////////////////////////////////
	// SERVER
	// /////////////////////////////////////////////////////////////////////////
	appCfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return appCfg, nil
}

func main() {
	// Initiates the appConfig
	appCfg, err := loadEnvVar()
	if err != nil {
		panic(err)
	}

	// Opens a connection to the databrase
	db, err := taskstore.Open(appCfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Starts Gorilla mux
	r := mux.NewRouter()
	r.StrictSlash(true)

	// Starts the task server
	server := taskstore.NewTaskServer()

	// Defines the routes
	r.HandleFunc("/task/", server.CreateTaskHandler).Methods(http.MethodPost)
	r.HandleFunc("/task/", server.GetAllTasksHandler).Methods(http.MethodGet)
	r.HandleFunc("/task/", server.DeleteAllTasksHandler).Methods(http.MethodDelete)
	r.HandleFunc("/task/{id:[0-9]+}/", server.GetTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/task/{id:[0-9]+}/", server.DeleteTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/tag/{tag}/", server.GetTaskByTagHandler).Methods(http.MethodGet)
	r.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]}/", server.GetTaskByDueHandler).
		Methods(http.MethodGet)

	fmt.Println("Starting server at port :3030...")
	log.Fatal(http.ListenAndServe("localhost:3030", taskstore.LoggingMiddleware(r.ServeHTTP)))
}
