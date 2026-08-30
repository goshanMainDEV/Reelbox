package httpserver

import (
	"fmt"
	"net/http"
	"strings"

	"reelbox/internal/cache"
	"reelbox/internal/extractor"
)

type App struct {
	Cache *cache.VideoCache
}

func (app *App) fetchVideoInfoCached(id string) (*extractor.YtDlpResult, error) {
	if cached, ok := app.Cache.Get(id); ok {
		fmt.Println("Кэш-попадание для id:", id)
		return cached, nil
	}

	fmt.Println("Кэш-промах, запрашиваю yt-dlp для id:", id)
	instagramURL := "https://www.instagram.com/reel/" + id + "/"

	info, err := extractor.FetchVideoInfo(instagramURL)
	if err != nil {
		return nil, err
	}

	app.Cache.Set(id, info)
	return info, nil
}

func (app *App) ReelHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "TelegramBot") {
		fmt.Println("Запрос от Telegram-краулера, id:", id)
	} else {
		fmt.Println("Запрос от обычного пользователя, id:", id)
	}

	info, err := app.fetchVideoInfoCached(id)
	if err != nil {
		http.Error(w, "Не удалось получить видео", http.StatusInternalServerError)
		fmt.Println("Ошибка fetchVideoInfoCached:", err)
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

func (app *App) ProxyVideoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	fmt.Println("proxyVideoHandler вызван, id:", id)

	info, err := app.fetchVideoInfoCached(id)
	if err != nil {
		http.Error(w, "Не удалось получить видео", http.StatusInternalServerError)
		fmt.Println("Ошибка fetchVideoInfoCached:", err)
		return
	}

	videoURL, audioURL := extractor.FindBestVideoAndAudio(info.Formats)

	w.Header().Set("Content-Type", "video/mp4")

	if err := extractor.StreamMergedVideo(videoURL, audioURL, w); err != nil {
		fmt.Println("Ошибка стриминга:", err)
		http.Error(w, "Video streaming failed", http.StatusInternalServerError)
		return
	}
}