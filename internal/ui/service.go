package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func New(config *Config, deps *Dependencies) (*Service, error) {
	return &Service{
		Config: config,
	}, nil
}

func (s *Service) Run(uiMsgChan chan any) error {
	teaProgram := tea.NewProgram(s.initUIModel())

	go func() {
		for uiState := range uiMsgChan {
			teaProgram.Send(uiState)
		}
	}()

	if _, err := teaProgram.Run(); err != nil {
		log.Fatalf("failed to initialize terminal ui: %v", err)
	}

	return nil
}

func (s *Service) initUIModel() Model {
	state := s.Config.InitialState

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e445ff"))

	state.Spinner = sp

	return state
}
