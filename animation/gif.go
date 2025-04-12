package animation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// GenerateGIF creates an animated GIF with transparent background using a browser animation
func GenerateGIF(outputFile string) error {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<style>
				html, body { 
					background-color: rgba(0,0,0,0) !important;
					margin: 0;
					padding: 0;
					width: 500px;
					height: 300px;
				}
				.container { 
					display: flex; 
					flex-direction: column; 
					gap: 10px;
					background-color: rgba(0,0,0,0) !important;
				}
				.square { 
					width: 30px; 
					height: 30px; 
					background: #f43f5e; 
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="square"></div>
				<div class="square"></div>
				<div class="square"></div>
				<div class="square"></div>
				<div class="square"></div>
			</div>
			<script src="https://cdn.jsdelivr.net/npm/animejs@4.0.1/lib/anime.umd.min.js"></script>
			<script>
				const { animate, stagger } = anime;
				animate('.square', {
					x: 320,
					rotate: { from: -180 },
					duration: 1250,
					delay: stagger(65, { from: 'center' }),
					ease: 'inOutQuint',
					loop: true,
					alternate: true
				});
			</script>
		</body>
		</html>
	`

	// Create temporary directory for frames
	tmpDir := "frames"
	os.MkdirAll(tmpDir, 0755)
	tmpHtml := filepath.Join(tmpDir, "index.html")
	err := os.WriteFile(tmpHtml, []byte(html), 0644)
	if err != nil {
		return err
	}

	// Initialize headless browser with transparency flags
	browser := rod.New().MustConnect()

	// Configure browser for transparency with specific options
	page := browser.MustPage("")
	page.MustSetViewport(500, 300, 1, false)
	page.MustNavigate("file://" + filepath.Join(getCWD(), tmpHtml))
	page.MustWaitLoad()

	// Force transparent background via JavaScript
	page.MustEval(`() => {
		document.documentElement.style.backgroundColor = 'rgba(0,0,0,0)';
		document.body.style.backgroundColor = 'rgba(0,0,0,0)';
		document.querySelector('.container').style.backgroundColor = 'rgba(0,0,0,0)';
	}`)

	// Wait for animation to initialize
	time.Sleep(1 * time.Second)

	// Capture frames with transparency
	frameCount := 60
	frameRate := 20
	for i := 0; i < frameCount; i++ {
		imgPath := filepath.Join(tmpDir, fmt.Sprintf("frame-%03d.png", i))

		// Use the correct Screenshot method with transparency
		data, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format:                proto.PageCaptureScreenshotFormatPng,
			FromSurface:           true,
			CaptureBeyondViewport: true,
		})

		// Save the screenshot
		os.WriteFile(imgPath, data, 0644)

		time.Sleep(time.Second / time.Duration(frameRate))
	}

	// Use ImageMagick to ensure transparency in PNGs
	for i := 0; i < frameCount; i++ {
		imgPath := filepath.Join(tmpDir, fmt.Sprintf("frame-%03d.png", i))
		cmd := exec.Command("convert", imgPath, "-transparent", "white", imgPath)
		cmd.Run()
	}

	// Generate GIF with ffmpeg, ensuring transparency
	cmd := exec.Command("ffmpeg", "-y", "-framerate", fmt.Sprint(frameRate),
		"-i", filepath.Join(tmpDir, "frame-%03d.png"),
		"-vf", "split[s0][s1];[s0]palettegen=reserve_transparent=1:transparency_color=ffffff[p];[s1][p]paletteuse=alpha_threshold=128",
		"-loop", "0",
		strings.Replace(outputFile, ".mp4", ".gif", 1),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Get current working directory
func getCWD() string {
	dir, _ := os.Getwd()
	return dir
}

// CleanupFrames removes the temporary frames directory and its contents
func CleanupFrames() error {
	tmpDir := "frames"
	// Remove all files in the directory
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err = os.Remove(filepath.Join(tmpDir, entry.Name()))
		if err != nil {
			return err
		}
	}

	// Remove the directory itself
	return os.RemoveAll(tmpDir)
}
