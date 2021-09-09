package main

import (
	"flag"
	"fmt"
	"log"
	"myapp/internal/driver"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// holds configuration information for the app
type config struct {
	port int
	env  string // dev or prod
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
	version       string
}

func (app *application) serve() error {
	srv := &http.Server {
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: 30 * time.Second,
		ReadTimeout: 10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	app.infoLog.Println(fmt.Sprintf("Starting Back end server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}

func main() {
	var cfg config // command line flag

	flag.IntVar(&cfg.port, "port", 4001, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.db.dsn, "dsn", "widgets:widgets@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "DSN")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		config: cfg,
		infoLog: infoLog,
		errorLog: errorLog,
		version: version,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}