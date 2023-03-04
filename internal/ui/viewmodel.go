package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	humanize "github.com/dustin/go-humanize"
)

var (
	baseStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		PaddingLeft(1).
		PaddingRight(1).
		BorderForeground(lipgloss.Color("240"))
)

func (v Model) Init() tea.Cmd {
	return v.UIReadyCallback
}

func (v Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg == nil {
		return v.updateTable(msg)
	}

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		}
	case State:
		if msgType.Completed {
			return v, tea.Quit
		}

		return v.updateTable(msgType)
	}

	return v, cmd
}

func (v Model) getTableOptions(columns []table.Column, rows []table.Row) []table.Option {
	return []table.Option{
		table.WithRows(rows),
		table.WithHeight(len(rows)),
		table.WithColumns(columns),
		table.WithStyles(table.Styles{
			Header: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderTop(false).
				BorderLeft(false).
				BorderRight(false),
			Cell: lipgloss.NewStyle(),
		}),
	}
}

func (v Model) updateTable(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	paddingRight := 2

	columns := []table.Column{
		{Title: "Track", Width: 20 + paddingRight},
		{Title: "Artist", Width: 20 + paddingRight},
		{Title: "Downloaded", Width: 10 + paddingRight},
		{Title: "  ", Width: 2},
	}

	rows := make([]table.Row, 0)

	if uiState, ok := msg.(State); ok {
		numTracks := len(uiState.Tracks)

		rows = make([]table.Row, 0, numTracks)

		for i := range uiState.Tracks {
			track := uiState.Tracks[i]

			var downloadIcon string

			if track.DownloadCompleted {
				downloadIcon = "✅"
			} else if track.DownloadError != nil {
				downloadIcon = "❌"
			} else {
				downloadIcon = "🔥"
			}

			rows = append(rows, table.Row{
				track.Title,
				track.Artist,
				humanize.Bytes(track.DownloadProgress),
				downloadIcon,
			})
		}

		v.Tracks = uiState.Tracks
	}

	table := table.New(v.getTableOptions(columns, rows)...)

	v.Table = table

	v.Table, cmd = v.Table.Update(msg)

	return v, cmd
}

func (v Model) View() string {
	var ui string

	ui += baseStyle.Render(v.Table.View())

	return fmt.Sprintf(
		"%s\nquit (ctrl+q)",
		ui,
	)
}
