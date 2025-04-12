package animation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateGIF tests the GIF generation function
func TestGenerateGIF(t *testing.T) {
	// Create test output file
	outputFile := "test-animation.gif"

	// Remove output file if it exists
	_ = os.Remove(outputFile)

	// Generate the animation
	err := GenerateGIF(outputFile)
	if err != nil {
		t.Fatalf("Failed to generate animation: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s was not created", outputFile)
	}

	// Check file size (should be non-zero)
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatalf("Output file is empty")
	}

	fmt.Printf("GIF animation created successfully: %s (%d bytes)\n", outputFile, fileInfo.Size())

	// Clean up
	_ = os.Remove(outputFile)
}

// TestCleanupFrames tests the frame cleanup function
func TestCleanupFrames(t *testing.T) {
	// Create test directory and files
	tmpDir := "frames"
	os.MkdirAll(tmpDir, 0755)

	// Create some test files
	testFiles := []string{"test1.png", "test2.png", "test3.png"}
	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		err := os.WriteFile(path, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	// Call cleanup
	err := CleanupFrames()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Check if directory was removed
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Fatalf("Temporary directory was not cleaned up")
	}
}
