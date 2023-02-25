package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		BorderForeground(lipgloss.Color("240"))
)

func (v Model) Init() tea.Cmd {
	return tea.Batch(v.UIReadyCallback, v.Spinner.Tick)
}

func (v Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg == nil {
		v.Initial = true

		return v, cmd
	}

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		}
	case State:
		if msgType.FetchingMeta {
			v.FetchingMeta = msgType.FetchingMeta

			return v, v.Spinner.Tick
		}

		if msgType.Downloading {
			v.FetchingMeta = false

			v.FetchedMeta = msgType.FetchedMeta
			v.Downloading = msgType.Downloading
			v.DownloadProgress = msgType.DownloadProgress

			return v, cmd
		}
	}

	if v.Initial {
		v.Spinner, cmd = v.Spinner.Update(msg)

		return v, cmd
	}

	if v.FetchingMeta {
		v.Spinner, cmd = v.Spinner.Update(msg)

		return v, cmd
	}

	if v.Downloading {
		v.Table, cmd = v.Table.Update(msg)

		return v, cmd
	}

	return v, cmd
}

func (v Model) View() string {
	if v.FetchingMeta {
		return v.Spinner.View()
	}
	
	if v.Downloading {
		return baseStyle.Render(fmt.Sprint(
			v.Table.View(),
		))
	}

	return "Something is wrong"
}
