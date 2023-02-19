package service

import (
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/ui"
	"bandcamp_downloader/internal/urlfetcher"
	"time"
)

type Config struct {
	UIModelChan chan ui.UIModel
}

type Dependencies struct {
	URLFetcher *urlfetcher.Service
	Downloader *downloader.Service
	UI         *ui.Service
}

type Options struct {
	Timeout   time.Duration
	OutputDir string
}

type Service struct {
	config                   *Config
	urlFetcher               *urlfetcher.Service
	downloader               *downloader.Service
	ui                       *ui.Service
	onFetchMetaEvents        []func(meta urlfetcher.AudioMeta)
	onDownloadTrackEvents    []func()
	onDownloadPlaylistEvents []func()
}
