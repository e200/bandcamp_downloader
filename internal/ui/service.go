package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{
		Config: config,
	}, nil
}

func (s *Service) Run(uiViewModelChan chan UIModel) error {
	teaProgram := tea.NewProgram(s.initUIModel())

	go func() {
		for uiViewModel := range uiViewModelChan {
			if uiViewModel.Loading {
				sp := spinner.New()
				sp.Spinner = spinner.Dot

				uiViewModel.spinner = sp

				teaProgram.Send(loadingState)
			} else {
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

				uiViewModel.table = trackTable

				teaProgram.Send(downloadingState)
			}
		}
	}()

	if _, err := teaProgram.Run(); err != nil {
		log.Fatalf("failed to initialize terminal ui: %v", err)
	}

	return nil
}

func (s *Service) initUIModel() UIModel {
	sp := spinner.New()

	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	return UIModel{
		Loading: true,
		spinner: sp,
	}
}
