package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	loadingState     UIState = "loading-state"
	downloadingState UIState = "downloading-state"
)

var (
	baseStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		BorderForeground(lipgloss.Color("240"))
)

func (v UIModel) Init() tea.Cmd {
	return v.spinner.Tick
}

func (v UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		}
	case UIState:
		switch msgType {
		case loadingState:
			v.Loading = true

			return v, tea.Batch(v.spinner.Tick)
		}
	}

	if v.Loading {
		v.spinner, cmd = v.spinner.Update(msg)

		return v, cmd
	}

	if v.Downloading {
		v.table, cmd = v.table.Update(msg)
	}

	return v, cmd
}

func (v UIModel) View() string {
	var ui string

	if v.Loading {
		ui = fmt.Sprint(
			v.spinner.View(),
			"Fetching tracks metadata...",
		)
	} else {
		ui = baseStyle.Render(fmt.Sprint(
			// "Bandcamp Downloader (v0.1)\n",
			v.table.View(),
		))
	}

	return fmt.Sprint(
		ui,
		"\n(ctrl+c to quit)",
	)
}
