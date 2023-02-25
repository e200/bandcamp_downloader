package service

import (
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/ui"
	"bandcamp_downloader/internal/urlfetcher"
	"time"
)

type Config struct {
	UIStateChan chan any
}

type Dependencies struct {
	URLFetcher *urlfetcher.Service
	Downloader *downloader.Service
}

type Options struct {
	Timeout   time.Duration
	OutputDir string
}

type Service struct {
	config     *Config
	urlFetcher *urlfetcher.Service
	downloader *downloader.Service
	ui         *ui.Service
}
