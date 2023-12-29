package main

import (
	"fmt"
	"html/template"
	"net/http"
	"url-shortener/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var tpl = template.Must(template.ParseFiles("../../internal/web-site/public/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// TODO: добавить graceful shutdown for server

func main() {
	cfg := config.MustLoad()

	// db, err := postgres.New(cfg.Storage_name)
	// if err != nil {
	// 	fmt.Printf("failde to init storage: %s", cfg.Storage_name)
	// 	return
	// }

	// _ = db

	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	fs := http.FileServer(http.Dir("../../internal/web-site/public/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", indexHandler)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Idle_timeout,
	}

	fmt.Println("Its okay")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed ro start server")
	}
}
