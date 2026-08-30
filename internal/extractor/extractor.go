package extractor

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Format struct {
	FormatID string `json:"format_id"`
	VCodec   string `json:"vcodec"`
	ACodec   string `json:"acodec"`
	URL      string `json:"url"`
}

type YtDlpResult struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Thumbnail string   `json:"thumbnail"`
	URL       string   `json:"url"`
	Formats   []Format `json:"formats"`
}

// FetchVideoInfo — экспортируемая версия (была fetchVideoInfo с маленькой буквы).
// Теперь её будет вызывать пакет httpserver, поэтому она обязана быть публичной.
func FetchVideoInfo(videoURL string) (*YtDlpResult, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", videoURL)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("yt-dlp failed: %w\nstderr: %s", err, exitErr.Stderr)
		}
		return nil, fmt.Errorf("yt-dlp failed: %w", err)
	}

	var result YtDlpResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &result, nil
}

// FindBestVideoAndAudio — тоже стала публичной (была findBestVideoAndAudio).
func FindBestVideoAndAudio(formats []Format) (videoURL string, audioURL string) {
	for _, f := range formats {
		if f.VCodec != "none" && f.VCodec != "" && videoURL == "" {
			videoURL = f.URL
		}
		if f.ACodec != "none" && f.ACodec != "" && audioURL == "" {
			audioURL = f.URL
		}
	}
	return videoURL, audioURL
}