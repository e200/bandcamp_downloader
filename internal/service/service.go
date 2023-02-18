package downloader

import (
	"context"
)

const ()

var ()

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{
		Config:      config,
		urlFetcher: deps.URLFetcher,
		downloader: deps.Downloader,
	}, nil
}

func (s *Service) DownloadTrack(
	context context.Context,
	URL string,
	options Options,
) error {
	return nil
}

func (s *Service) DownloadMany(
	context context.Context,
	URLs []string,
	options Options,
) error {
	return nil
}
