package downloader

type Config struct {
}

type Dependencies struct{}

type Options struct {
	Filepath string
}

type Service struct{
	downloadListeners []func(progress int)
}

type DownloadWriter struct {
	total int
	listeners []func(progress int)
}
