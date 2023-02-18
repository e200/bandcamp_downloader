package main

import (
	"flag"
	"log"
	"time"

	"bandcamp_downloader/internal/urlfetcher"
)

func main() {
	var (
		trackURL    string
		playlistURL string
		// outputDir   string

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

	bcdown, err := urlfetcher.New(
		&urlfetcher.Config{},
		&urlfetcher.Dependencies{},
	)
	if err != nil {
		log.Fatal(err)
	}

	timeoutDuration := time.Duration(timeout) * time.Second

	if trackURL != "" {
		_, err := bcdown.FetchAudioURL(
			trackURL, urlfetcher.Options{
				Timeout:   timeoutDuration,
			})
		if err != nil {
			log.Fatalf("error downloading track: %v", err)
		}
	}

	if playlistURL != "" {
		_, err := bcdown.FetchAudioURLS(playlistURL, urlfetcher.Options{
			Timeout:   timeoutDuration,
		})
		if err != nil {
			log.Fatalf("error downloading playlist: %v", err)
		}
	}
}
