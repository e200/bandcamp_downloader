package urlfetcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/chromedp/cdproto/cdp"
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

func (s *Service) AddFetchingListener(listener func()) {
	s.fetchingListeners = append(s.fetchingListeners, listener)
}

func (s *Service) AddFetchedListener(listener func(meta AudioMeta)) {
	s.fetchedListeners = append(s.fetchedListeners, listener)
}

func (s *Service) FetchAudioURL(
	ctx context.Context,
	trackURL string,
	options *Options,
) (*AudioMeta, error) {
	if len(s.fetchingListeners) > 0 {
		for i := range s.fetchingListeners {
			go s.fetchingListeners[i]()
		}
	}

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var (
		isTrackAudioURLAvailable bool

		trackAudioURL string
		trackTitle    string
		trackArtist   string
	)

	if err := s.dispatch(
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
	); err != nil {
		return nil, err
	}

	if !isTrackAudioURLAvailable {
		return nil, ErrAudioURLNotAvailable
	}

	audioMeta := AudioMeta{
		Title:  trackTitle,
		Artist: trackArtist,
		URL:    trackAudioURL,
	}

	if len(s.fetchedListeners) > 0 {
		for i := range s.fetchedListeners {
			go s.fetchedListeners[i](audioMeta)
		}
	}

	return &audioMeta, nil
}

func (s *Service) FetchAudioURLS(
	ctx context.Context,
	playlistURL string,
	options *Options,
) ([]AudioMeta, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var nodes []*cdp.Node

	if err := s.dispatch(
		ctx,
		chromedp.Navigate(playlistURL),
		chromedp.WaitVisible(".track_list"),
		chromedp.Nodes(".track_list .play_status", &nodes),
	); err != nil {
		return nil, err
	}

	tracksActions := []chromedp.Action{
		chromedp.Navigate(playlistURL),
		chromedp.WaitVisible(".track_list"),
	}

	audioTrackAvailabilities := make([]bool, len(nodes))
	audioMetas := make([]AudioMeta, len(nodes))

	for i := range nodes {
		audioMetas[i] = AudioMeta{}

		nodeSelector := fmt.Sprintf(
			".track_list .track_row_view:nth-child(%d) .play_status",
			i+1,
		)

		tracksActions = append(
			tracksActions,
			chromedp.Click(nodeSelector),
			chromedp.AttributeValue(
				audioSelector,
				"src",
				&audioMetas[i].URL,
				&audioTrackAvailabilities[i],
			),
			chromedp.Text(".track_info .title", &audioMetas[i].Title),
			chromedp.Text("#name-section > h3 > span > a", &audioMetas[i].Artist),
		)
	}

	if err := s.dispatch(
		ctx,
		tracksActions...,
	); err != nil {
		return nil, err
	}

	// some track may not be downloadble
	/* if !isTrackAudioURLAvailable {
		return nil, ErrAudioURLNotAvailable
	} */

	return audioMetas, nil
}

func (s *Service) dispatch(
	ctx context.Context,
	actions ...chromedp.Action,
) error {
	response, err := chromedp.RunResponse(
		ctx,
		actions...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", audioURLGettingErrorMessage, err)
	}

	if response.Status >= http.StatusBadRequest {
		if response.Status == http.StatusNotFound {
			return fmt.Errorf("%s. %v", audioURLGettingErrorMessage, ErrPageNotFound)
		}

		return fmt.Errorf("%s. %v, status code: %d, status message: %s",
			audioURLGettingErrorMessage,
			ErrBadResponse,
			response.Status,
			response.StatusText,
		)
	}

	return nil
}
