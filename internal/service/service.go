package service

import (
	"context"
	"net/http"

	"github.com/chromedp/chromedp"
)

const (
	audioSelector = "audio[src^=\"https://\"]"
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) DownloadTrack(trackURL string, options DownloadOptions) error {
	var isTrackAudioURLAvailable bool
	var trackAudioURL string

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	response, err := chromedp.RunResponse(
		ctx,
		chromedp.Navigate(trackURL),
		chromedp.WaitVisible(".playbutton"),
		chromedp.Click(".playbutton"),
		chromedp.AttributeValue(
			audioSelector,
			"src",
			&trackAudioURL,
			&isTrackAudioURLAvailable,
		),
	)
	if err != nil {
		return err
	}

	if response.Status >= http.StatusBadRequest {
		return nil
	}

	return nil
}

func (s *Service) DownloadPlaylist(trackURL string, options DownloadOptions) error {
	return nil
}
