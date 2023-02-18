package downloader

import (
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/urlfetcher"
)

type Config struct {
}

type Dependencies struct {
	URLFetcher *urlfetcher.Service
	Downloader *downloader.Service
}

type Options struct {
}

type Service struct {
	Config      *Config
	urlFetcher *urlfetcher.Service
	downloader *downloader.Service
}
