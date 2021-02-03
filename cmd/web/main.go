package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"time"

	"awesomeProject/pkg/models"

	"github.com/alexedwards/scs"
	_ "github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "aserty"
	dbname   = "snippetbox"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("CONN"), "PostGreSQL")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	secret := flag.String("secret", "s6Nd%+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets directory")
	flag.Parse()

	db, err := connect(*dsn)
	defer db.Close()

	sessionManager := scs.NewCookieManager(*secret)
	sessionManager.Lifetime(12 * time.Hour)
	sessionManager.Persist(true)

	app := &App{
		Database:  &models.Database{DB: db},
		HTMLDir:   *htmlDir,
		Sessions:  sessionManager,
		StaticDir: *staticDir,
	}

	log.Printf("Server listening on %s", *addr)
	err = http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}

func connect(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
