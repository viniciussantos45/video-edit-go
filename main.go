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
	"github.com/mowshon/moviego"
)

func main() {
	fmt.Println("Hello, World!")

	err := generateAnimationMP4("output.gif")
	if err != nil {
		panic(err)
	}
	fmt.Println("MP4 gerado com sucesso!")

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

	concat.Output("concat.mp4").Run()
}

func generateAnimationMP4(outputFile string) error {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<style>
				body { background: transparent; }
				.container { display: flex; flex-direction: column; gap: 10px; }
				.square { width: 30px; height: 30px; background: #f43f5e; }
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

	// Inicia o navegador headless
	browser := rod.New().MustConnect()
	page := browser.MustPage("file://" + filepath.Join(getCWD(), tmpHtml))
	page.MustWaitLoad()

	// Espera animação inicializar
	time.Sleep(1 * time.Second)

	// Captura de frames
	frameCount := 60
	frameRate := 20
	for i := 0; i < frameCount; i++ {
		imgPath := filepath.Join(tmpDir, fmt.Sprintf("frame-%03d.png", i))
		page.MustScreenshot(imgPath)
		time.Sleep(time.Second / time.Duration(frameRate))
	}

	// Gera o GIF com ffmpeg
	cmd := exec.Command("ffmpeg", "-y", "-framerate", fmt.Sprint(frameRate),
		"-i", filepath.Join(tmpDir, "frame-%03d.png"),
		"-vf", "split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse",
		"-loop", "0",
		strings.Replace(outputFile, ".gif", ".gif", 1),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCWD() string {
	dir, _ := os.Getwd()
	return dir
}
