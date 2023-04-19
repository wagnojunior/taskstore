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
	"github.com/wagnojunior/taskstore/internal/handlers"
	"github.com/wagnojunior/taskstore/internal/migrations"
	"github.com/wagnojunior/taskstore/internal/models"
	"github.com/wagnojunior/taskstore/internal/repository/dbrepo"
	"github.com/wagnojunior/taskstore/internal/taskstore"
)

// `appConfig` configures all env variables required for this application to run
type appConfig struct {
	PSQL   models.PostgresConfig
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
	db, err := models.Open(appCfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Runs the goose database migration
	err = dbrepo.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Creates the necessary services
	storeService := &dbrepo.PostgresRepo{
		DB: db,
	}
	taskService := &dbrepo.PostgresRepo{
		DB: db,
	}

	// Creates the necessary handlers
	storeHandler := &handlers.Store{
		StoreService: storeService,
	}
	taskHanler := &handlers.Tasks{
		TaskService: taskService,
	}

	// Starts Gorilla mux
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/store", storeHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/delete/{store_id:[0-9]+}", storeHandler.DeleteByID).Methods(http.MethodGet)

	r.HandleFunc("/task", taskHanler.Create).Methods(http.MethodPost)
	r.HandleFunc("/task/{store_id:[0-9]+}/{id:[0-9]+}", taskHanler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/task/{store_id:[0-9]+}/tags", taskHanler.GetByTags).Methods(http.MethodPost)
	r.HandleFunc("/task/{store_id:[0-9]+}", taskHanler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/delete/{store_id:[0-9]+}/{id:[0-9]+}", taskHanler.DeleteByID).Methods(http.MethodGet)
	r.HandleFunc("/delete/{store_id:[0-9]+}/all", taskHanler.DeleAll).Methods(http.MethodGet)

	fmt.Println("Starting server at port :3030...")
	log.Fatal(http.ListenAndServe("localhost:3030",
		taskstore.LoggingMiddleware(r.ServeHTTP)))
}
