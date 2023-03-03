package ui

import (
	"bandcamp_downloader/internal/urlfetcher"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	InitialState Model
}

type Dependencies struct {
}

type Service struct {
	Config *Config
}

type State struct {
	FetchingMeta         bool
	FetchedMeta          urlfetcher.AudioMeta
	FetchingMetas        bool
	Downloading          bool
	DownloadProgress     uint64
	DownloadingMany      bool
	AllDownloadsComplete bool
}

type Model struct {
	Initial              bool
	FetchingMeta         bool
	FetchedMeta          *urlfetcher.AudioMeta
	FetchingMetas        bool
	Downloading          bool
	DownloadProgress     uint64
	DownloadingMany      bool
	AllDownloadsComplete bool
	UIReadyCallback      func() tea.Msg

	Spinner spinner.Model
	Table   table.Model
}
