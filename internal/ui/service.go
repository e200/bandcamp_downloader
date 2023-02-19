package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{
		Config: config,
	}, nil
}

func (s *Service) Run() error {
	paddingRight := 2

	tableColumns := []table.Column{
		{Title: "Track", Width: 20 + paddingRight},
		{Title: "Artist", Width: 20 + paddingRight},
		{Title: "Size", Width: 5 + paddingRight},
		{Title: "Down", Width: 5 + paddingRight},
		{Title: "Progress", Width: 15},
	}

	tableRows := []table.Row{
		{"Breath Easy", "Marie Joly", "10mb", "5.1mb", "===========>···"},
	}

	trackTable := table.New(
		table.WithColumns(tableColumns),
		table.WithRows(tableRows),
		table.WithHeight(7),
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

	p := tea.NewProgram(viewModel{
		table: trackTable,
	})
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to initialize terminal ui: %v", err)
	}

	return nil
}
