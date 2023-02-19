package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
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

type viewModel struct {
	table table.Model
}

func (v viewModel) Init() tea.Cmd {
	return nil
}

func (v viewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		}
	}

	v.table, cmd = v.table.Update(msg)

	return v, cmd
}

func (v viewModel) View() string {
	ui := fmt.Sprint(
		// "Bandcamp Downloader (v0.1)\n",
		v.table.View(),
	)

	baseUi := baseStyle.Render(ui)

	return fmt.Sprint(
		baseUi,
		"\n(ctrl+c to quit)",
	)
}
