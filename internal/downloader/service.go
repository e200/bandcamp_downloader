package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
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
	bytes, err := s.getFileBytes(context, URL, options.ProgressListener)
	if err != nil {
		return err
	}

	err = s.saveFile(bytes, options.Filepath)
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

func (s *Service) getFileBytes(context context.Context, URL string, listener func(progress uint64)) ([]byte, error) {
	request, err := http.NewRequestWithContext(context, http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	if request.Body != nil {
		defer request.Body.Close()
	}

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

	var writer DownloadWriter

	if listener != nil {
		writer.AddListener(listener)
	}

	bytes, err := io.ReadAll(io.TeeReader(response.Body, &writer))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fileDownloadErrorMessage, err)
	}

	return bytes, nil
}

func (s *Service) saveFile(bytes []byte, filePath string) error {
	err := os.WriteFile(filePath, bytes, fs.FileMode(0777))
	if err != nil {
		return fmt.Errorf("%s: %w", fileSavingErrorMessage, err)
	}

	return nil
}
