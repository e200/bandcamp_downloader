package service

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) DownloadTrack(trackURL string, options DownloadOptions) error {
	return nil
}

func (s *Service) DownloadPlaylist(trackURL string, options DownloadOptions) error {
	return nil
}
