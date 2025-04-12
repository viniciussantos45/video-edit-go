package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/mowshon/moviego"
)

func main() {
	fmt.Println("Hello, World!")

	// Generate the GIF with transparency
	gifFile := "animation.gif"
	err := generateAnimationMP4(gifFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("GIF gerado com sucesso!")

	first, err := moviego.Load("forest.mp4")
	if err != nil {
		fmt.Errorf("error loading video: %v", err)
	}

	var wg sync.WaitGroup

	// Process first subclip
	wg.Add(1)
	go func() {
		defer wg.Done()
		first.SubClip(0, 5).FadeOut(0.2).Output("cutFirst.mp4").Run()
	}()

	// Process second subclip
	wg.Add(1)
	go func() {
		defer wg.Done()
		first.SubClip(5, 10).FadeIn(1, 0.2).Output("cutSecond.mp4").Run()
	}()

	// Wait for both subclips to complete
	wg.Wait()

	cutFirst, _ := moviego.Load("cutFirst.mp4")
	cutSecond, _ := moviego.Load("cutSecond.mp4")

	concat, _ := moviego.Concat(
		[]moviego.Video{
			cutFirst,
			cutSecond,
		},
	)

	// Generate the concatenated video
	concat.Output("concat.mp4").Run()

	// Now overlay the GIF on the final video
	overlayGIFOnVideo("concat.mp4", gifFile, "final_with_gif.mp4")

	fmt.Println("Video final com GIF gerado com sucesso!")
}

func generateAnimationMP4(outputFile string) error {
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

	// Criar arquivo HTML temporário
	tmpDir := "frames"
	os.MkdirAll(tmpDir, 0755)
	tmpHtml := filepath.Join(tmpDir, "index.html")
	err := os.WriteFile(tmpHtml, []byte(html), 0644)
	if err != nil {
		return err
	}

	// Inicia o navegador headless com flags para transparência
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

	// Espera animação inicializar
	time.Sleep(1 * time.Second)

	// Captura de frames with transparency
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

	// Gera o GIF com ffmpeg, ensuring transparency
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

// OverlayGIFOnVideo adds a GIF to a video using ffmpeg
func overlayGIFOnVideo(videoFile, gifFile, outputFile string) error {
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

func getCWD() string {
	dir, _ := os.Getwd()
	return dir
}
