package downloader

type Config struct {
}

type Dependencies struct{}

type Options struct {
	Filepath string
}

type Service struct {
	downloadProgressListeners []func(progress uint64)
	downloadCompleteListeners []func()
}

type DownloadWriter struct {
	total     int
	listeners []func(progress uint64)
}
