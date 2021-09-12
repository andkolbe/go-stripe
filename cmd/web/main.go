package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"myapp/internal/driver"
	"myapp/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const version = "1.0.0"

// const cssVersion = "1"
var session *scs.SessionManager

// holds configuration information for the app
type config struct {
	port int
	env  string // dev or prod
	api  string // the url we call to the backend api
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

// receiver
type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}

func main() {
	gob.Register(TransactionData{})
	var cfg config // command line flag

	flag.IntVar(&cfg.port, "port", 4000, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.db.dsn, "dsn", "widgets:widgets@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "DSN")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	flag.Parse()


	

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	// set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
		Session: session,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
