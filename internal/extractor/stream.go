package extractor

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// StreamMergedVideo — была streamMergedVideo, тоже стала публичной.
func StreamMergedVideo(videoURL, audioURL string, dst io.Writer) error {
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