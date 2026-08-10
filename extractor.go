package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func fetchVideoInfo(videoURL string) (*YtDlpResult, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", videoURL)

	output, err := cmd.Output()
	if err != nil {
		// Если ошибка именно от завершения процесса — достанем stderr для диагностики.
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

func findBestVideoAndAudio(formats []Format) (videoURL string, audioURL string) {
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

func streamMergedVideo(videoURL, audioURL string, dst io.Writer) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", videoURL,
		"-i", audioURL,
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-c:a", "copy",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	if _, err := io.Copy(dst, stdout); err != nil {
		return fmt.Errorf("failed to stream ffmpeg output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg exited with error: %w\nstderr: %s", err, stderrBuf.String())
	}

	fmt.Println("ffmpeg stderr log:", stderrBuf.String())

	return nil
}
