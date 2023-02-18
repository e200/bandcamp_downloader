package urlfetcher

import "time"

type Config struct {
}

type Dependencies struct{}

type DownloadOptions struct {
	Timeout time.Duration
	OutputDir string
}

type Service struct{}
