package urlfetcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/chromedp/chromedp"
)

const (
	audioSelector = "audio[src^=\"https://\"]"
)

var (
	ErrPageNotFound         = errors.New("page not found")
	ErrGettingAudioURL      = errors.New("error while getting audio URL")
	ErrRequestError         = errors.New("request error while getting audio URL")
	ErrAudioURLNotAvailable = errors.New("audio URL not available")
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) FetchAudioURL(trackURL string, options DownloadOptions) (string, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	var isTrackAudioURLAvailable bool
	var trackAudioURL string

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
		return "", fmt.Errorf("%s: %w", ErrGettingAudioURL, err)
	}

	if response.Status >= http.StatusBadRequest {
		switch response.Status {
		case http.StatusNotFound:
			return "", fmt.Errorf("request error while getting track audio URL. %v", ErrPageNotFound)
		default:
			return "", fmt.Errorf("%s. status code: %d", ErrRequestError, response.Status)
		}
	}

	if !isTrackAudioURLAvailable {
		return "", ErrAudioURLNotAvailable
	}

	return trackAudioURL, nil
}

func (s *Service) FetchAudioURLS(playlistURL string, options DownloadOptions) ([]string, error) {
	return nil, nil
}
