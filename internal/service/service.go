package service

import (
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/ui"
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
		config:     config,
		urlFetcher: deps.URLFetcher,
		downloader: deps.Downloader,
	}, nil
}

func (s *Service) Init(uiMsgChan chan any) error {
	if err := s.ui.Run(uiMsgChan); err != nil {
		return fmt.Errorf("unable to initiate terminal ui, error: %w", err)
	}

	return nil
}

func (s *Service) DownloadTrack(
	trackURL string,
	options Options,
) error {
	uiMsgChan := s.config.UIStateChan

	if uiMsgChan != nil {
		s.urlFetcher.AddFetchingListener(func() {
			uiMsgChan <- ui.State{
				FetchingMeta: true,
			}
		})

		s.urlFetcher.AddFetchedListener(func(meta urlfetcher.AudioMeta) {
			uiMsgChan <- ui.State{
				FetchedMeta: meta,
			}
		})

		s.AddDownloadTrackListener(func(progress int) {
			uiMsgChan <- ui.State{
				Downloading:      true,
				DownloadProgress: progress,
			}
		})
	}

	s.resolveOptions(options)

	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	audioMeta, err := s.urlFetcher.FetchAudioURL(ctx, trackURL, nil)
	if err != nil {
		return err
	}

	filename := s.getFilename(*audioMeta)

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

func (s *Service) AddDownloadTrackListener(listener func(progress int)) {
	s.downloader.AddDownloadListener(listener)
}

func (s *Service) AddDownloadPlaylistListener(listener func(progress int)) {
	// s.onDownloadPlaylistEvents = append(s.onDownloadPlaylistEvents, listener)
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

func (s *Service) resolveOptions(options Options) error {
	if options.OutputDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		options.OutputDir = cwd
	}

	return nil
}
