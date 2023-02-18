package main

import (
	"flag"
	"log"
	"time"

	"bandcamp_downloader/internal/service"
)

func main() {
	var (
		trackURL    string
		playlistURL string
		outputDir   string

		timeout int64
	)

	flag.StringVar(&trackURL, "track", "", "Track URL")
	flag.StringVar(&playlistURL, "playlist", "", "Playlist URL")

	flag.Int64Var(&timeout, "timeout", 60, "Timeout duration in seconds")

	flag.Parse()

	if trackURL == "" && playlistURL == "" {
		log.Fatal("either a track or playlist URL must be provided")
	}

	if trackURL != "" && playlistURL != "" {
		log.Fatal("you must provide a track or playlist URL, not both")
	}

	bcdown, err := service.New(
		&service.Config{},
		&service.Dependencies{},
	)
	if err != nil {
		log.Fatal(err)
	}

	timeoutDuration := time.Duration(timeout) * time.Second

	if trackURL != "" {
		if err := bcdown.DownloadTrack(
			trackURL, service.DownloadOptions{
				Timeout:   timeoutDuration,
				OutputDir: outputDir,
			}); err != nil {
			log.Fatalf("error downloading track: %v", err)
		}
	}

	if playlistURL != "" {
		if err := bcdown.DownloadPlaylist(playlistURL, service.DownloadOptions{
			Timeout:   timeoutDuration,
			OutputDir: outputDir,
		}); err != nil {
			log.Fatalf("error downloading playlist: %v", err)
		}
	}
}
