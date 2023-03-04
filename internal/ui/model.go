package ui

import (
	"bandcamp_downloader/internal/audiosmetadatafetcher"

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
	Error     error
	Tracks    []TrackState
	Completed bool
}

type TrackState struct {
	DownloadProgress  uint64
	DownloadCompleted bool
	DownloadError     error
	audiosmetadatafetcher.AudioMeta
}

type Model struct {
	Tracks               []TrackState
	AllDownloadsComplete bool
	UIReadyCallback      func() tea.Msg

	Spinner spinner.Model
	Table   table.Model
}
