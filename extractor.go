package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type YtDlpResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
}

func fetchVideoInfo(videoURL string) (*YtDlpResult, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", videoURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp failed: %w\noutput: %s", err, output)
	}

	var result YtDlpResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &result, nil
}