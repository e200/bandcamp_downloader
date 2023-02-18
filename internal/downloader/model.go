package downloader

import "time"

type Config struct {
}

type Dependencies struct{}

type Options struct {
	Timeout   time.Duration
	OutputDir string
}

type Service struct{}
