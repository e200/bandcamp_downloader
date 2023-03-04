package service

import (
	"bandcamp_downloader/internal/audiosmetadatafetcher"
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/ui"
	"time"
)

type Config struct {
	UIStateChan chan any
}

type Dependencies struct {
	AMF        *audiosmetadatafetcher.Service
	Downloader *downloader.Service
}

type Options struct {
	Timeout   time.Duration
	OutputDir string
}

type Service struct {
	config                *Config
	audiosmetadatafetcher *audiosmetadatafetcher.Service
	downloader            *downloader.Service
	ui                    *ui.Service
}
