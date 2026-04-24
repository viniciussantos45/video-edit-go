# video-edit-go

A Go-based video editing tool that automates common video processing tasks such as trimming, concatenation, animated GIF generation, and GIF overlay on video.

## Objective

The objective of this project is to explore video manipulation in Go, leveraging the `moviego` library (a Go wrapper around FFmpeg) to perform editing operations programmatically — cutting clips, applying fade effects, concatenating segments, generating animated GIFs with transparency, and compositing them onto video.

## Features

- Cut video into sub-clips with configurable start and end times
- Apply fade-in and fade-out effects on clips
- Concatenate multiple video clips into a single output
- Generate animated GIFs from frames
- Overlay a GIF animation on top of a video file
- Concurrent processing of video clips using goroutines

## Tech Stack

- Go
- [moviego](https://github.com/mowshon/moviego) — FFmpeg wrapper for Go
- FFmpeg (required as system dependency)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/viniciussantos45/video-edit-go.git
cd video-edit-go

# Install dependencies
go mod download

# Place your source video as forest.mp4 and run
go run main.go
```
