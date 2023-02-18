package service

import (
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/urlfetcher"
	"time"
)

type Config struct {
}

type Dependencies struct {
	URLFetcher *urlfetcher.Service
	Downloader *downloader.Service
}

type Options struct {
	Timeout   time.Duration
	OutputDir string
	Filename  string
}

type Service struct {
	Config      *Config
	urlFetcher *urlfetcher.Service
	downloader *downloader.Service
}
