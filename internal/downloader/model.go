package downloader

import "time"

type Config struct {
}

type Dependencies struct{}

type Options struct {
	Timeout   time.Duration
	Filename  string
	OutputDir string
}

type Service struct{}
