package urlfetcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/chromedp/chromedp"
)

const (
	audioSelector               = "audio[src^=\"https://\"]"
	audioURLGettingErrorMessage = "error while getting audio URL"
)

var (
	ErrPageNotFound         = errors.New("page not found")
	ErrBadRequest           = errors.New("bad request error")
	ErrBadResponse          = errors.New("bad response error")
	ErrAudioURLNotAvailable = errors.New("audio URL not available")
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) FetchAudioURL(
	context context.Context,
	trackURL string,
	options *Options,
) (string, error) {
	ctx, cancel := chromedp.NewContext(context)
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
		return "", fmt.Errorf("%s: %w", audioURLGettingErrorMessage, err)
	}

	if response.Status >= http.StatusBadRequest {
		switch response.Status {
		case http.StatusNotFound:
			return "", fmt.Errorf("%s. %v", audioURLGettingErrorMessage, ErrPageNotFound)
		default:
			return "", fmt.Errorf("%s. %v, status code: %d, status message: %s",
				audioURLGettingErrorMessage,
				ErrBadResponse,
				response.Status,
				response.StatusText,
			)
		}
	}

	if !isTrackAudioURLAvailable {
		return "", ErrAudioURLNotAvailable
	}

	return trackAudioURL, nil
}

func (s *Service) FetchAudioURLS(
	context context.Context,
	playlistURL string,
	options *Options,
) ([]string, error) {
	return nil, nil
}
