package service

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{}, nil
}

type Service struct {}
