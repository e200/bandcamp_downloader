package main

import (
	"flag"
	"log"

	"bandcamp_downloader/internal/service"
)

func main() {
	var trackURL string
	var playlistURL string

	flag.StringVar(&trackURL, "t", "", "Track URL")
	flag.StringVar(&playlistURL, "p", "", "Playlist URL")

	flag.Parse()

	if trackURL == "" && playlistURL == "" {
		log.Fatal("either a track or playlist URL must be provided")
	}

	if trackURL != "" && playlistURL != "" {
		log.Fatal("you must provide a track or playlist URL, not both")
	}

	_, err := service.New(
		&service.Config{},
		&service.Dependencies{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
