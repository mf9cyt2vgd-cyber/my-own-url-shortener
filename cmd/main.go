package main

import (
	"fmt"
	"my_own_shortener/internal/config"
	redirecth "my_own_shortener/internal/http-server/handlers/redirect"
	deleteh "my_own_shortener/internal/http-server/handlers/url/delete"
	saveh "my_own_shortener/internal/http-server/handlers/url/save"
	updateh "my_own_shortener/internal/http-server/handlers/url/update"
	"my_own_shortener/internal/logger"
	"time"

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
	r := chi.NewRouter()
	r.Method("GET", "/{alias}", redirecth.NewRedirectHandler(log, db))
	r.Route("/url", func(r chi.Router) {
		r.Patch("/", updateh.NewUpdateHandler(log, db))
		r.Post("/", saveh.NewSaveHandler(log, db))
		r.Delete("/", deleteh.NewDeleteHandler(log, db))
	})
	server := &http.Server{
		Addr:              "localhost:8080", // или cfg.HTTP.Address если есть
		Handler:           r,
		ReadTimeout:       cfg.HttpServer.Timeout,     // Макс время на чтение запроса
		ReadHeaderTimeout: 2 * time.Second,            // Макс время на чтение заголовков
		WriteTimeout:      10 * time.Second,           // Макс время на запись ответа
		IdleTimeout:       cfg.HttpServer.IdleTimeout, // Макс время простоя
	}

	log.Info("Starting server...")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
