package main

import (
	"flag"
	"log"
	"time"

	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/service"
	"bandcamp_downloader/internal/ui"
	"bandcamp_downloader/internal/urlfetcher"
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
	flag.StringVar(&outputDir, "output", ".", "Playlist URL")

	flag.Int64Var(&timeout, "timeout", 60, "Timeout duration in seconds")

	flag.Parse()

	if trackURL == "" && playlistURL == "" {
		log.Fatal("either a track or playlist URL must be provided")
	}

	if trackURL != "" && playlistURL != "" {
		log.Fatal("you must provide a track or playlist URL, not both")
	}

	urlFetcher, err := urlfetcher.New(
		&urlfetcher.Config{},
		&urlfetcher.Dependencies{},
	)
	if err != nil {
		log.Fatalf("unable to instantiate urlfetcher service: %v", err)
	}

	down, err := downloader.New(
		&downloader.Config{},
		&downloader.Dependencies{},
	)
	if err != nil {
		log.Fatalf("unable to instantiate downloader service: %v", err)
	}

	UI, err := ui.New(
		&ui.Config{},
		&ui.Dependencies{},
	)
	if err != nil {
		log.Fatalf("unable to instantiate terminal ui service: %v", err)
	}

	uiModelChan := make(chan ui.UIModel)

	svc, err := service.New(
		&service.Config{},
		&service.Dependencies{
			URLFetcher: urlFetcher,
			Downloader: down,
			UI:         UI,
		},
	)
	if err != nil {
		log.Fatalf("unable to instantiate downloader service: %v", err)
	}

	svc.OnFetchMeta(func(meta urlfetcher.AudioMeta) {
		uiModelChan <- ui.UIModel{
			Loading: true,
		}
	})

	svc.OnDownloadTrack(func() {
		uiModelChan <- ui.UIModel{
			Downloading: true,
		}
	})

	svc.OnDownloadPlaylist(func() {
		uiModelChan <- ui.UIModel{
			Downloading: true,
		}
	})

	go func() {
		if trackURL != "" {
			uiModelChan <- ui.UIModel{
				Loading: true,
			}

			if err := svc.DownloadTrack(trackURL, service.Options{
				Timeout:   time.Duration(timeout) * time.Second,
				OutputDir: outputDir,
			}); err != nil {
				log.Fatalf("error downloading track: %v", err)
			}
		}

		if playlistURL != "" {
			if err := svc.DownloadPlaylist(playlistURL, &service.Options{
				Timeout: time.Duration(timeout) * time.Second,
			}); err != nil {
				if err != nil {
					log.Fatalf("error downloading playlist: %v", err)
				}
			}
		}
	}()

	if err = svc.Init(uiModelChan); err != nil {
		log.Fatalf("unable to initiate service: %v", err)
	}
}
