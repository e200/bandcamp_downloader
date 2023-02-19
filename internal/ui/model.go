package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
)

type Config struct {
}

type Dependencies struct {
}

type Service struct {
	Config *Config
}

type UIState string

type UIModel struct {
	Loading     bool
	Downloading bool
	table       table.Model
	spinner     spinner.Model
}
