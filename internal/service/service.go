package service

import (
	"bandcamp_downloader/internal/downloader"
	"context"
	"os"
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
	trackURL string,
	options *Options,
) error {
	s.resolveOptions(options)

	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	audioFileURL, err := s.urlFetcher.FetchAudioURL(ctx, trackURL, nil)
	if err != nil {
		return err
	}

	if err := s.downloader.Download(ctx, audioFileURL, &downloader.Options{
		OutputDir: options.OutputDir,
		Filename: "first_track.mp3",
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) DownloadPlaylist(
	playlistURL string,
	options *Options,
) error {
	return nil
}

func (s *Service) resolveOptions(options *Options) error {
	if options == nil {
		options = &Options{}
	}

	if options.OutputDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		options.OutputDir = cwd
	}

	return nil
}
