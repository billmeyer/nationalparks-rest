package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"nationalparks-rest/pkg"
	"nationalparks-rest/pkg/db"
	http2 "nationalparks-rest/pkg/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"net/http"
	"os"

	"github.com/XSAM/otelsql"
)

var dtb *sql.DB
var ctx context.Context

var dbHost string
var dbPort int
var httpHost string
var httpPort int

func main() {
	var err error

	cleanup := pkg.InitOTEL("nationalparks-rest", "development")
	defer cleanup(context.Background())

	processEnvVariables()

	log.Printf("Using MySQL instance at %s:%d", dbHost, dbPort)

	// Initialize the database connection
	var driverName string
	// Register an OTel driver
	driverName, err = otelsql.Register("mysql", semconv.DBSystemMySQL.Value.AsString())

	var connString = db.GetConnectionString("nationalparks_user", "nationalparks_user", dbHost, dbPort, "nationalparks_db")
	dtb, err = sql.Open(driverName, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %+v\n", err)
		os.Exit(-1)
	}
	defer dtb.Close()
	http2.SetDB(dtb)

	// Initialize the HTTP Router
	router := mux.NewRouter()
	var muxMiddleware = otelmux.Middleware("nationalparks-rest")
	router.Use(muxMiddleware)
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", http2.RouteHealthCheck).Methods(http.MethodGet)
	api.HandleFunc("/health-check", http2.RouteHealthCheck).Methods("GET")
	api.HandleFunc("/nationalpark/{id:[0-9]+}", http2.RouteGetNationalParkById).Methods(http.MethodGet)
	api.HandleFunc("/nationalparks", http2.RouteGetNationalParks).Methods(http.MethodGet)
	//api.HandleFunc("/nationalparks/name/{parkname}", http2.RouteGetNationalParkByName).Methods(http.MethodGet)
	//api.HandleFunc("/nationalparks/city/{city}", http2.RouteGetNationalParksByCity).Methods(http.MethodGet)
	//api.HandleFunc("/nationalparks/state/{stateabbr}", http2.RouteGetNationalParksByState).Methods(http.MethodGet)
	//api.HandleFunc("/nationalparks/zipcode/{zipcode}", http2.RouteGetNationalParksByZipCode).Methods(http.MethodGet)

	// Setup HTTP server
	var httpServer = httpHost + ":" + strconv.Itoa(httpPort)
	server := &http.Server{
		Handler: router,
		Addr:    httpServer,
		// timeouts so the server never waits forever...
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start accepting connections...
	log.Printf("Server started at %s", httpServer)
	log.Fatal(server.ListenAndServe())
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

func processEnvVariables() {
	var exists bool
	var val string

	dbHost = os.Getenv("DBHOST")
	if val, exists = os.LookupEnv("DBPORT"); exists == true {
		dbPort, _ = strconv.Atoi(val)
	} else {
		dbPort = 3306
	}
	httpHost = os.Getenv("HTTPHOST")
	if val, exists = os.LookupEnv("HTTPPORT"); exists == true {
		httpPort, _ = strconv.Atoi(val)
	} else {
		httpPort = 8080
	}
}

//func processCmdLine() {
//	flag.StringVar(&dbHost, "dbhost", "", "Hostname or IP Address of the MySQL server.")
//	flag.IntVar(&dbPort, "dbport", 3306, "Port number MySQL server is listening on.")
//	flag.StringVar(&httpHost,"httphost", "127.0.0.1", "TCP address for the server to listen.")
//	flag.IntVar(&httpPort,"httpport", 8080, "Port number the HTTP server should accept connections on.")
//	flag.Parse()
//
//	if len(dbHost) == 0 {
//		fmt.Fprintf(os.Stderr, "ERROR: dbhost not specified.\n")
//		fmt.Fprintf(os.Stderr, "Usage: %s <options>, where <options> are:\n", os.Args[0])
//		flag.PrintDefaults()
//		os.Exit(-1)
//	}
//}
