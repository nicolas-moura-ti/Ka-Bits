package main

import (
	"fmt"
	"ka-bits/pkg/game"
	"ka-bits/pkg/storage"
	"ka-bits/pkg/ui"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	player, err := storage.Load()
	if err != nil {
		log.Fatalf("Error loading the wheel of Ka: %v", err)
	}

	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	earnings, offlineTime := engine.ProcessOfflineEarnings()

	m := ui.NewModel(engine)
	if earnings > 0 {
		m.Logs = append(m.Logs, fmt.Sprintf("[INFO] While you were away, you mined %.2f bits in %v.", earnings, offlineTime.Round(time.Second)))
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("A glitch in the matrix occurred: %v", err)
		os.Exit(1)
	}
}
