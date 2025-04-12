package video

import (
	"fmt"
	"os"
	"os/exec"
)

// OverlayGIFOnVideo adds a GIF to a video using ffmpeg
func OverlayGIFOnVideo(videoFile, gifFile, outputFile string) error {
	// Use ffmpeg to overlay the GIF on the video
	// Position: 10px from top-right corner
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-i", videoFile,
		"-ignore_loop", "0", // Make sure animated GIF loops
		"-i", gifFile,
		"-filter_complex", "[0:v][1:v]overlay=main_w-overlay_w-10:10:shortest=1",
		"-c:a", "copy", // Copy audio stream
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreateSubclip creates a subclip from a video file
func CreateSubclip(inputFile string, startTime, duration float64, outputFile string) error {
	// Format start time and duration as strings
	startTimeStr := formatTime(startTime)
	durationStr := formatTime(duration)

	// Create subclip using ffmpeg
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-i", inputFile,
		"-ss", startTimeStr,
		"-t", durationStr,
		"-c", "copy", // Copy streams without re-encoding
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ConcatVideos concatenates multiple video files
func ConcatVideos(inputFiles []string, outputFile string) error {
	// Create a temporary concat file
	concatFile := "concat.txt"
	var concatContent string
	for _, file := range inputFiles {
		concatContent += "file '" + file + "'\n"
	}

	// Write the concat file
	err := os.WriteFile(concatFile, []byte(concatContent), 0644)
	if err != nil {
		return err
	}

	// Concatenate videos using the concat demuxer
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	// Clean up concat file
	os.Remove(concatFile)

	return err
}

// ApplyEffect applies visual effects to a video
func ApplyEffect(inputFile, effect, outputFile string) error {
	var filterArgs string

	// Choose filter based on effect name
	switch effect {
	case "fadeIn":
		filterArgs = "fade=t=in:st=0:d=1"
	case "fadeOut":
		filterArgs = "fade=t=out:st=0:d=1"
	case "blur":
		filterArgs = "boxblur=10:1:cr=0:ar=0"
	// Add more effects as needed
	default:
		// No effect, just copy
		return simpleCopy(inputFile, outputFile)
	}

	// Apply the effect using ffmpeg
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-i", inputFile,
		"-vf", filterArgs,
		"-c:a", "copy", // Copy audio
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Helper function to format time as string
func formatTime(seconds float64) string {
	return fmt.Sprintf("%.2f", seconds)
}

// Helper function for simple copy when no effect is applied
func simpleCopy(inputFile, outputFile string) error {
	cmd := exec.Command(
		"ffmpeg", "-y",
		"-i", inputFile,
		"-c", "copy",
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
