package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type VideoMeta struct {
	Title    string
	ThumbURL string
	VideoURL string
}

var ogTemplate = template.Must(template.New("og").Parse(`<!DOCTYPE html>
<html prefix="og: https://ogp.me/ns#">
<head>
	<meta charset="utf-8">
	<meta property="og:type" content="video.other">
	<meta property="og:title" content="{{.Title}}">
	<meta property="og:image" content="{{.ThumbURL}}">
	<meta property="og:video" content="{{.VideoURL}}">
	<meta property="og:video:type" content="video/mp4">
</head>
<body>
	<video controls poster="{{.ThumbURL}}">
		<source src="{{.VideoURL}}" type="video/mp4">
	</video>
</body>
</html>`))

func main() {
	info, err := fetchVideoInfo("https://www.instagram.com/reel/DbtqWSnBKOe/")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Получено видео:", info.Title)
		fmt.Println("URL:", info.URL)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /reel/{id}", reelHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func reelHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fmt.Println("Запрошен рилс с ID:", id)

	meta := VideoMeta{
		Title:    "Рилс " + id,
		ThumbURL: "https://placehold.co/720x1280.png",
		VideoURL: "https://example.com/video.mp4",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ogTemplate.Execute(w, meta)
}