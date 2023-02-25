package downloader

func (w *DownloadWriter) Write(p []byte) (n int, err error) {
	total := len(p)

	w.total += total
	w.UpdateProgress()

	return total, nil
}

func (w *DownloadWriter) UpdateProgress() {
	for i := range w.listeners {
		w.listeners[i](uint64(w.total))
	}
}

func (w *DownloadWriter) AddListener(listener func(progress uint64)) {
	w.listeners = append(w.listeners, listener)
}
