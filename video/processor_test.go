package video

import (
	"fmt"
	"os"
	"testing"
)

// TestOverlayGIFOnVideo tests the GIF overlay function
func TestOverlayGIFOnVideo(t *testing.T) {
	// Skip if video file doesn't exist
	if _, err := os.Stat("../forest.mp4"); os.IsNotExist(err) {
		t.Skip("Test video file forest.mp4 not found, skipping test")
	}

	// Create a simple test GIF file
	gifFile := "test-overlay.gif"
	createTestGIF(gifFile, t)

	// Output file
	outputFile := "test-with-overlay.mp4"

	// Remove output file if it exists
	_ = os.Remove(outputFile)

	// Overlay GIF on video
	err := OverlayGIFOnVideo("../forest.mp4", gifFile, outputFile)
	if err != nil {
		t.Fatalf("Failed to overlay GIF on video: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s was not created", outputFile)
	}

	// Check file size
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatalf("Output file is empty")
	}

	fmt.Printf("Video with GIF overlay created successfully: %s (%d bytes)\n", outputFile, fileInfo.Size())

	// Clean up
	_ = os.Remove(gifFile)
	_ = os.Remove(outputFile)
}

// TestCreateSubclip tests the subclip function
func TestCreateSubclip(t *testing.T) {
	// Skip if video file doesn't exist
	if _, err := os.Stat("../forest.mp4"); os.IsNotExist(err) {
		t.Skip("Test video file forest.mp4 not found, skipping test")
	}

	outputFile := "test-subclip.mp4"

	// Remove output file if it exists
	_ = os.Remove(outputFile)

	// Create the subclip
	err := CreateSubclip("../forest.mp4", 0, 2, outputFile)
	if err != nil {
		t.Fatalf("Failed to create subclip: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s was not created", outputFile)
	}

	// Check file size
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatalf("Output file is empty")
	}

	fmt.Printf("Subclip created successfully: %s (%d bytes)\n", outputFile, fileInfo.Size())

	// Clean up
	_ = os.Remove(outputFile)
}

// TestConcatVideos tests the video concatenation function
func TestConcatVideos(t *testing.T) {
	// Skip if video file doesn't exist
	if _, err := os.Stat("../forest.mp4"); os.IsNotExist(err) {
		t.Skip("Test video file forest.mp4 not found, skipping test")
	}

	// Create test subclips
	clip1 := "test-clip1.mp4"
	clip2 := "test-clip2.mp4"
	outputFile := "test-concat.mp4"

	// Remove output files if they exist
	_ = os.Remove(clip1)
	_ = os.Remove(clip2)
	_ = os.Remove(outputFile)

	// Create two test clips
	err := CreateSubclip("../forest.mp4", 0, 2, clip1)
	if err != nil {
		t.Fatalf("Failed to create first test clip: %v", err)
	}

	err = CreateSubclip("../forest.mp4", 2, 2, clip2)
	if err != nil {
		t.Fatalf("Failed to create second test clip: %v", err)
	}

	// Concatenate the clips
	err = ConcatVideos([]string{clip1, clip2}, outputFile)
	if err != nil {
		t.Fatalf("Failed to concatenate videos: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s was not created", outputFile)
	}

	// Check file size
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatalf("Output file is empty")
	}

	fmt.Printf("Concatenated video created successfully: %s (%d bytes)\n", outputFile, fileInfo.Size())

	// Clean up
	_ = os.Remove(clip1)
	_ = os.Remove(clip2)
	_ = os.Remove(outputFile)
}

// Helper function to create a test GIF file
func createTestGIF(filename string, t *testing.T) {
	// Generate a simple 1x1 pixel GIF file
	gifData := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00, 0x00, 0xff, 0xff, 0xff,
		0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x21, 0xf9, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2c, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44, 0x01, 0x00, 0x3b,
	}

	err := os.WriteFile(filename, gifData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test GIF file: %v", err)
	}
}
