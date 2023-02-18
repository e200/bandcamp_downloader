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
) (*AudioMeta, error) {
	ctx, cancel := chromedp.NewContext(context)
	defer cancel()

	var (
		isTrackAudioURLAvailable bool

		trackAudioURL string
		trackTitle    string
		trackArtist   string
	)

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
		chromedp.Text(".trackTitle", &trackTitle),
		chromedp.Text(".albumTitle > span > a", &trackArtist),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", audioURLGettingErrorMessage, err)
	}

	if response.Status >= http.StatusBadRequest {
		switch response.Status {
		case http.StatusNotFound:
			return nil, fmt.Errorf("%s. %v", audioURLGettingErrorMessage, ErrPageNotFound)
		default:
			return nil, fmt.Errorf("%s. %v, status code: %d, status message: %s",
				audioURLGettingErrorMessage,
				ErrBadResponse,
				response.Status,
				response.StatusText,
			)
		}
	}

	if !isTrackAudioURLAvailable {
		return nil, ErrAudioURLNotAvailable
	}

	return &AudioMeta{
		Title:  trackTitle,
		Artist: trackArtist,
		URL:    trackAudioURL,
	}, nil
}

func (s *Service) FetchAudioURLS(
	context context.Context,
	playlistURL string,
	options *Options,
) ([]string, error) {
	return nil, nil
}
