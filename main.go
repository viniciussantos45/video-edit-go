package main

import (
	"fmt"
	"sync"

	"github.com/mowshon/moviego"
	"github.com/viniciussantos45/video-edit-go/animation"
	"github.com/viniciussantos45/video-edit-go/video"
)

func main() {
	fmt.Println("Starting video processing...")

	// Generate the GIF with transparency
	gifFile := "animation.gif"
	err := animation.GenerateGIF(gifFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("GIF generated successfully!")

	// Process video using moviego library
	processVideoWithMoviego()

	// Overlay the GIF on the concatenated video
	err = video.OverlayGIFOnVideo("concat.mp4", gifFile, "final_with_gif.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Println("Final video with GIF overlay created successfully!")

	// Optional: cleanup temporary files
	animation.CleanupFrames()
}

// processVideoWithMoviego uses the moviego library to process video clips
func processVideoWithMoviego() {
	// Load the source video
	first, err := moviego.Load("forest.mp4")
	if err != nil {
		fmt.Errorf("error loading video: %v", err)
		return
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

	// Load the processed clips
	cutFirst, _ := moviego.Load("cutFirst.mp4")
	cutSecond, _ := moviego.Load("cutSecond.mp4")

	// Concatenate the clips
	concat, _ := moviego.Concat(
		[]moviego.Video{
			cutFirst,
			cutSecond,
		},
	)

	// Generate the concatenated video
	concat.Output("concat.mp4").Run()

	fmt.Println("Video clips processed and concatenated successfully!")
}
