package main

import (
	"fmt"
	"net/http"
	"time"

	"reelbox/internal/cache"
	"reelbox/internal/httpserver"
)

func main() {
	app := &httpserver.App{
		Cache: cache.NewVideoCache(10 * time.Minute),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /reel/{id}", app.ReelHandler)
	mux.HandleFunc("GET /proxy/video", app.ProxyVideoHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
