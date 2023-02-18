package downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	fileDownloadErrorMessage = "error while downloading file"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Download(
	context context.Context,
	URL string,
	options Options,
) error {
	request, err := http.NewRequestWithContext(context, http.MethodGet, URL, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	defer request.Body.Close()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		switch response.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("%s. %v", fileDownloadErrorMessage, ErrFileNotFound)
		default:
			return fmt.Errorf(
				"%s. status code: %d, status message: %s",
				fileDownloadErrorMessage,
				response.StatusCode,
				response.Status,
			)
		}
	}

	return nil
}

func (s *Service) DownloadMany(
	context context.Context,
	URLs []string,
	options Options,
) error {
	return nil
}
