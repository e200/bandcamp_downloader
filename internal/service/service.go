package service

import (
	"bandcamp_downloader/internal/audiosmetadatafetcher"
	"bandcamp_downloader/internal/downloader"
	"bandcamp_downloader/internal/ui"
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
		config:                config,
		audiosmetadatafetcher: deps.AMF,
		downloader:            deps.Downloader,
	}, nil
}

func (s *Service) Init(uiMsgChan chan any) error {
	if err := s.ui.Run(uiMsgChan); err != nil {
		return fmt.Errorf("unable to initiate terminal ui, error: %w", err)
	}

	return nil
}

func (s *Service) DownloadTracks(
	sourceURL string,
	options Options,
) error {
	uiMsgChan := s.config.UIStateChan

	s.resolveOptions(options)

	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	tracksMetadata, err := s.audiosmetadatafetcher.Fetch(ctx, sourceURL)
	if err != nil {
		uiMsgChan <- ui.State{
			Error:     err,
			Completed: true,
		}

		close(uiMsgChan)
	}

	var wg sync.WaitGroup

	numTracks := len(tracksMetadata)

	wg.Add(numTracks)

	tracksStates := make([]ui.TrackState, 0, numTracks)

	for i := range tracksMetadata {
		tracksStates = append(tracksStates, ui.TrackState{
			AudioMeta: tracksMetadata[i],
		})

		go func(currentTrackIndex int, tracksStates []ui.TrackState) {
			if err != nil {
				tracksStates[currentTrackIndex].DownloadError = err

				wg.Done()

				return
			}

			filename := s.getFilename(tracksStates[currentTrackIndex].AudioMeta)

			if err := s.downloader.Download(
				ctx, tracksStates[currentTrackIndex].URL, downloader.Options{
					Filepath: path.Join(options.OutputDir, filename),
					ProgressListener: func(progress uint64) {
						tracksStates[currentTrackIndex].DownloadProgress = progress

						uiMsgChan <- ui.State{
							Tracks: tracksStates,
						}
					},
				}); err != nil {
				tracksStates[currentTrackIndex].DownloadError = err

				wg.Done()

				return
			}

			tracksStates[currentTrackIndex].DownloadCompleted = true

			uiMsgChan <- ui.State{
				Tracks: tracksStates,
			}

			wg.Done()
		}(i, tracksStates)
	}

	wg.Wait()

	uiMsgChan <- ui.State{
		Completed: true,
	}

	close(uiMsgChan)

	return nil
}

func (*Service) getFilename(audioMeta audiosmetadatafetcher.AudioMeta) string {
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
