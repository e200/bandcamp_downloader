package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

const (
	fileDownloadErrorMessage = "error while downloading file"
	fileSavingErrorMessage   = "error while saving downloaded file"
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
	bytes, err := s.getFileBytes(context, URL)
	if err != nil {
		return err
	}

	filePath := path.Join(options.OutputDir, options.Filename)

	err = s.saveFile(bytes, filePath)
	if err != nil {
		return err
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

func (s *Service) getFileBytes(context context.Context, URL string) ([]byte, error) {
	request, err := http.NewRequestWithContext(context, http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	defer request.Body.Close()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		switch response.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("%s. %v", fileDownloadErrorMessage, ErrFileNotFound)
		default:
			return nil, fmt.Errorf(
				"%s. status code: %d, status message: %s",
				fileDownloadErrorMessage,
				response.StatusCode,
				response.Status,
			)
		}
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	return bytes, nil
}

func (s *Service) saveFile(bytes []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		return fmt.Errorf("%s: %w", fileSavingErrorMessage, err)
	}

	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("%s: %w", fileSavingErrorMessage, err)
	}

	return nil
}
