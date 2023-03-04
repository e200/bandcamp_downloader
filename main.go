package main

import (
	"flag"
	"log"
	"time"

	"bandcamp_downloader/internal/audiosmetadatafetcher"
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/service"
	"bandcamp_downloader/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var (
		URL       string
		outputDir string

		timeout int64
	)

	flag.StringVar(&URL, "url", "", "Track or playlist URL")
	flag.StringVar(&outputDir, "output", ".", "Output directory")

	flag.Int64Var(&timeout, "timeout", 60, "Timeout duration in seconds")

	flag.Parse()

	if URL == "" {
		log.Fatal("a bandcamp URL (track or playlist) must be provided")
	}

	audiosmetadatafetcher, err := audiosmetadatafetcher.New(
		&audiosmetadatafetcher.Config{},
		&audiosmetadatafetcher.Dependencies{},
	)
	if err != nil {
		log.Fatalf("unable to instantiate audiosmetadatafetcher service: %v", err)
	}

	down, err := downloader.New(
		&downloader.Config{},
		&downloader.Dependencies{},
	)
	if err != nil {
		log.Fatalf("unable to instantiate downloader service: %v", err)
	}

	uiMsgChan := make(chan any)

	svc, err := service.New(
		&service.Config{
			UIStateChan: uiMsgChan,
		},
		&service.Dependencies{
			AMF:        audiosmetadatafetcher,
			Downloader: down,
		},
	)
	if err != nil {
		log.Fatalf("unable to instantiate downloader service: %v", err)
	}

	UI, err := ui.New(
		&ui.Config{
			InitialState: ui.Model{
				UIReadyCallback: func() tea.Msg {
					go func() {
						if err := svc.DownloadTracks(URL, service.Options{
							Timeout:   time.Duration(timeout) * time.Second,
							OutputDir: outputDir,
						}); err != nil {
							log.Fatalf("error downloading track: %v", err)
						}
					}()

					return nil
				},
			},
		},
		&ui.Dependencies{},
	)
	if err != nil {
		log.Fatalf("unable to instantiate terminal ui service: %v", err)
	}

	if err = UI.Run(uiMsgChan); err != nil {
		log.Fatalf("unable to initiate service: %v", err)
	}
}
