package main

import (
	"os"
	"testing"
	"time"

	"github.com/viniciussantos45/video-edit-go/animation"
	"github.com/viniciussantos45/video-edit-go/video"
)

// TestEndToEndVideoProcessing performs an integration test of the entire workflow
func TestEndToEndVideoProcessing(t *testing.T) {
	// Skip if video file doesn't exist
	if _, err := os.Stat("forest.mp4"); os.IsNotExist(err) {
		t.Skip("Test video file forest.mp4 not found, skipping test")
	}

	// Define test output files
	gifFile := "test-end-to-end.gif"
	finalVideo := "test-final.mp4"

	// Clean up previous test files
	_ = os.Remove(gifFile)
	_ = os.Remove(finalVideo)

	// Test both functions in sequence with a short timeout
	testTimeout := 60 * time.Second // Timeout after 60 seconds
	done := make(chan bool, 1)

	go func() {
		// 1. Generate GIF
		err := animation.GenerateGIF(gifFile)
		if err != nil {
			t.Errorf("Failed to generate GIF: %v", err)
			done <- false
			return
		}

		// 2. Overlay GIF on the video file
		err = video.OverlayGIFOnVideo("forest.mp4", gifFile, finalVideo)
		if err != nil {
			t.Errorf("Failed to overlay GIF: %v", err)
			done <- false
			return
		}

		// 3. Check if file exists
		if _, err := os.Stat(finalVideo); os.IsNotExist(err) {
			t.Errorf("Final video was not created")
			done <- false
			return
		}

		// Get file info to report size
		fi, err := os.Stat(finalVideo)
		if err == nil {
			t.Logf("Final video created successfully: %s (%d bytes)", finalVideo, fi.Size())
		}

		done <- true
	}()

	// Wait for completion or timeout
	select {
	case success := <-done:
		if !success {
			t.Fatal("End-to-end test failed")
		}
	case <-time.After(testTimeout):
		t.Fatal("End-to-end test timed out after", testTimeout)
	}

	// Clean up
	_ = os.Remove(gifFile)
	_ = os.Remove(finalVideo)
	animation.CleanupFrames() // Clean up any temporary frames
}
