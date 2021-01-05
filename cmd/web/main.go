package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"snippetbox/pkg/models"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "aserty1234"
	dbname   = "snippetbox"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := "user=postgres dbname=snippetbox password=aserty1234 host=localhost sslmode=disable"
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets directory")
	flag.Parse()

	db := connect(dsn)
	defer db.Close()


	app := &App{
		Database:  &models.Database{db},
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}


	log.Printf("Server listening on %s", *addr)
	err := http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}


func connect(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
