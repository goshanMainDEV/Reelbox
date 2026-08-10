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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /reel/{id}", reelHandler)
	mux.HandleFunc("GET /proxy/video", proxyVideoHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func reelHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	instagramURL := "https://www.instagram.com/reel/" + id + "/"

	info, err := fetchVideoInfo(instagramURL)
	if err != nil {
		http.Error(w, "Не удалось получить видео", http.StatusInternalServerError)
		fmt.Println("Ошибка fetchVideoInfo:", err)
		return
	}

	meta := VideoMeta{
		Title:    info.Title,
		ThumbURL: info.Thumbnail,
		VideoURL: "http://" + r.Host + "/proxy/video?id=" + id,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ogTemplate.Execute(w, meta)
}

func proxyVideoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("proxyVideoHandler вызван, id:", r.URL.Query().Get("id"))
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	instagramURL := "https://www.instagram.com/reel/" + id + "/"

	info, err := fetchVideoInfo(instagramURL)
	if err != nil {
		http.Error(w, "Не удалось получить видео", http.StatusInternalServerError)
		fmt.Println("Ошибка fetchVideoInfo:", err)
		return
	}

	videoURL, audioURL := findBestVideoAndAudio(info.Formats)

	w.Header().Set("Content-Type", "video/mp4")

	if err := streamMergedVideo(videoURL, audioURL, w); err != nil {
		fmt.Println("Ошибка стриминга:", err)
		http.Error(w, "Video streaming failed", http.StatusInternalServerError)
		return
	}
	// Важно: если заголовки уже отправлены (стриминг начался),
	// http.Error здесь может не сработать как ожидается —
	// это нормальное ограничение потоковых ответов (спорно)

}
func truncate(s string, max int) string {
	if len(s) > max {
		return s
	}
	return s[:max] + "..."
}
