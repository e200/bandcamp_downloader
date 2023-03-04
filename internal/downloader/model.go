package downloader

type Config struct {
}

type Dependencies struct{}

type Options struct {
	Filepath string
	ProgressListener func(progress uint64)
}

type Service struct {}

type DownloadWriter struct {
	total     int
	listeners []func(progress uint64)
}
