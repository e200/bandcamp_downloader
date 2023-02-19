package service

import (
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/urlfetcher"
	"context"
	"fmt"
	"os"
	"path"
	"sync"
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

	if len(s.onFetchMetaEvents) > 0 {
		for i := range s.onFetchMetaEvents {
			go s.onFetchMetaEvents[i](*audioMeta)
		}
	}

	filename := s.getFilename(*audioMeta)

	if err := s.downloader.Download(ctx, audioMeta.URL, downloader.Options{
		Filepath: path.Join(options.OutputDir, filename),
	}); err != nil {
		return err
	}

	return nil
}

func (*Service) getFilename(audioMeta urlfetcher.AudioMeta) string {
	filename := fmt.Sprintf(
		"%s - %s.%s",
		audioMeta.Artist,
		audioMeta.Title,
		outputFileFormat,
	)

	return filename
}

func (s *Service) DownloadPlaylist(
	playlistURL string,
	options *Options,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	audioMetas, err := s.urlFetcher.FetchAudioURLS(
		ctx,
		playlistURL,
		&urlfetcher.Options{},
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for i := range audioMetas {
		wg.Add(i)

		audioMeta := audioMetas[i]

		go func() {
			defer wg.Done()

			err = s.downloader.Download(ctx, audioMeta.URL, downloader.Options{
				Filepath: s.getFilename(audioMeta),
			})
			if err != nil {
				return
			}
		}()
	}

	wg.Wait()

	return nil
}

func (s *Service) OnFetchMeta(callback func(meta urlfetcher.AudioMeta)) {
	s.onFetchMetaEvents = append(s.onFetchMetaEvents, callback)
}

func (s *Service) OnDownloadTrack(callback func()) {
	s.onDownloadTrackEvents = append(s.onDownloadTrackEvents, callback)
}

func (s *Service) OnDownloadPlaylist(callback func()) {
	s.onDownloadPlaylistEvents = append(s.onDownloadPlaylistEvents, callback)
}
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
