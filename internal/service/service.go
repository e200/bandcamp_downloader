package service

import (
	"bandcamp_downloader/internal/downloader"
	"context"
	"fmt"
	"os"
	"path"
)

const (
	outputFileFormat = "mp3"
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{
		Config:     config,
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

	audioMeta, err := s.urlFetcher.FetchAudioURL(ctx, trackURL, nil)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(
		"%s - %s.%s",
		audioMeta.Artist,
		audioMeta.Title,
		outputFileFormat,
	)

	if err := s.downloader.Download(ctx, audioMeta.URL, downloader.Options{
		Filepath: path.Join(options.OutputDir, filename),
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
