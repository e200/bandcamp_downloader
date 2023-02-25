package urlfetcher

type Config struct {
}

type Dependencies struct{}

type Options struct {
}

type AudioMeta struct {
	Title  string
	Artist string
	URL    string
}

type Service struct {
	fetchingListeners []func()
	fetchedListeners  []func(meta AudioMeta)
}
