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
		v.FetchingMeta = true

		return v, v.Spinner.Tick

	}

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		}
	case State:
		if msgType.FetchingMeta {
			v.FetchingMeta = true

			return v, v.Spinner.Tick
		}

		if msgType.Downloading {
			v.FetchingMeta = false
			v.Downloading = true

			if v.FetchedMeta == nil {
				v.FetchedMeta = &msgType.FetchedMeta
			}

			v.DownloadProgress = msgType.DownloadProgress

			paddingRight := 2

			columns := []table.Column{
				{Title: "Track", Width: 20 + paddingRight},
				{Title: "Artist", Width: 20 + paddingRight},
				{Title: "Downloaded", Width: 10},
			}

			rows := []table.Row{
				{
					v.FetchedMeta.Title,
					v.FetchedMeta.Artist,
					humanize.Bytes(v.DownloadProgress),
				},
			}

			table := table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithHeight(1),
				table.WithStyles(table.Styles{
					Header: lipgloss.NewStyle().
						Border(lipgloss.NormalBorder()).
						BorderForeground(lipgloss.Color("240")).
						BorderTop(false).
						BorderLeft(false).
						BorderRight(false),
					Cell: lipgloss.NewStyle(),
				}),
			)

			v.Table = table

			v.Table, cmd = v.Table.Update(msg)

			return v, cmd
		}

		if msgType.AllDownloadsComplete {
			return v, tea.Quit
		}
	}

	v.Spinner, cmd = v.Spinner.Update(msg)

	return v, cmd
}

func (v Model) View() string {
	var ui string

	if v.Initial {
		ui += fmt.Sprintf(
			"%sInitializing...",
			v.Spinner.View(),
		)
	}

	if v.FetchingMeta {
		ui += fmt.Sprintf(
			"%sFetching metadata...",
			v.Spinner.View(),
		)
	}

	if v.Downloading {
		ui += baseStyle.Render(v.Table.View())
	}

	return fmt.Sprintf(
		"%s\nquit (ctrl+q)",
		ui,
	)
}
