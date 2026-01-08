package main

import (
	"fmt"
	"my_own_shortener/internal/config"
	redirecth "my_own_shortener/internal/http-server/handlers/redirect"
	deleteh "my_own_shortener/internal/http-server/handlers/url/delete"
	saveh "my_own_shortener/internal/http-server/handlers/url/save"
	"my_own_shortener/internal/logger"

	"my_own_shortener/internal/storage/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New()
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		cfg.Postgresql.Host, cfg.Postgresql.Port, cfg.Postgresql.User,
		cfg.Postgresql.Dbname, cfg.Postgresql.Sslmode)
	db, err := postgresql.New(connStr)
	defer db.Close()
	if err != nil {
		fmt.Printf("failed to init storage: %s", err)
		return
	}
	redirectHandler := redirecth.NewRedirectHandler(log, db)
	saveHandler := saveh.NewSaveHandler(log, db)
	deleteHandler := deleteh.NewDeleteHandler(log, db)
	r := chi.NewRouter()
	r.Method("GET", "/{alias}", redirectHandler)
	r.Route("/url", func(r chi.Router) {
		r.Post("/", saveHandler)
		r.Delete("/", deleteHandler)
	})
	log.Info("Starting server...")
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
