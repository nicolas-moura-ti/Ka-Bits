package ui

import (
	"errors"
	"ka-bits/pkg/game"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate_TickMsg(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	modelResult, cmd := model.Update(tickMsg(time.Now()))

	_, ok := modelResult.(Model)
	if !ok {
		t.Errorf("Expected model to be of type Model")
	}

	if cmd == nil {
		t.Errorf("Expected non-nil tea.Cmd for tickMsg")
	}
}

func TestUpdate_FrameMsg(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	model.MiningEffect = 1

	modelResult, _ := model.Update(frameMsg(time.Now()))
	updatedModel := modelResult.(Model)

	if updatedModel.AnimationTick != 1 {
		t.Errorf("Expected AnimationTick to be 1, got %d", updatedModel.AnimationTick)
	}

	if updatedModel.MiningEffect != 0 {
		t.Errorf("Expected MiningEffect to be 0, got %d", updatedModel.MiningEffect)
	}
}

func TestUpdate_AutoSaveMsg(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	modelResult, _ := model.Update(autoSaveMsg(time.Now()))
	updatedModel := modelResult.(Model)

	lastLog := updatedModel.Logs[len(updatedModel.Logs)-1]
	if lastLog != "[INFO] Game auto-saved." {
		t.Errorf("Expected auto-save log, got %s", lastLog)
	}
}

func TestUpdate_RandomLogMsg(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	modelResult, _ := model.Update(randomLogMsg("Test random log"))
	updatedModel := modelResult.(Model)

	lastLog := updatedModel.Logs[len(updatedModel.Logs)-1]
	if lastLog != "Test random log" {
		t.Errorf("Expected random log, got %s", lastLog)
	}
}

func TestUpdate_SaveResultMsg(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	err := errors.New("test error")
	modelResult, _ := model.Update(saveResultMsg{err})
	updatedModel := modelResult.(Model)

	lastLog := updatedModel.Logs[len(updatedModel.Logs)-1]
	if lastLog != "[WARN] Error saving game: test error" {
		t.Errorf("Expected error log, got %s", lastLog)
	}
}
func TestHandleKey_Mining(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	modelResult, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	updatedModel := modelResult.(Model)

	if updatedModel.Engine.Player.Bits != 1 {
		t.Errorf("Expected Bits to be 1, got %f", updatedModel.Engine.Player.Bits)
	}

	if updatedModel.MiningEffect != 3 {
		t.Errorf("Expected MiningEffect to be 3, got %d", updatedModel.MiningEffect)
	}
}

func TestHandleKey_CursorMovement(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	// Move down
	modelResult, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	updatedModel := modelResult.(Model)
	if updatedModel.Cursor != 1 {
		t.Errorf("Expected cursor to be 1, got %d", updatedModel.Cursor)
	}

	// Move up
	modelResult, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	updatedModel = modelResult.(Model)
	if updatedModel.Cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", updatedModel.Cursor)
	}
}

func TestHandleKey_ResetConfirm(t *testing.T) {
	player := game.NewPlayer()
	player.Bits = 100
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	// Press 'r' to trigger reset confirmation
	modelResult, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updatedModel := modelResult.(Model)
	if !updatedModel.ConfirmingReset {
		t.Errorf("Expected ConfirmingReset to be true")
	}

	// Press 'n' to cancel
	modelResult, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	updatedModel = modelResult.(Model)
	if updatedModel.ConfirmingReset {
		t.Errorf("Expected ConfirmingReset to be false")
	}

	if updatedModel.Engine.Player.Bits != 100 {
		t.Errorf("Expected Bits to be 100, got %f", updatedModel.Engine.Player.Bits)
	}

	// Trigger again and press 'y' to confirm
	updatedModel.ConfirmingReset = true
	modelResult, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	updatedModel = modelResult.(Model)
	if updatedModel.ConfirmingReset {
		t.Errorf("Expected ConfirmingReset to be false")
	}

	if updatedModel.Engine.Player.Bits != 0 {
		t.Errorf("Expected Bits to be 0 after reset, got %f", updatedModel.Engine.Player.Bits)
	}
}

func TestHandleKey_Prestige(t *testing.T) {
	player := game.NewPlayer()
	// Need at least 500k for prestige
	player.TotalBitsEver = 500000
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	// Press 'p'
	modelResult, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	updatedModel := modelResult.(Model)

	if !updatedModel.ConfirmingPrestige {
		t.Errorf("Expected ConfirmingPrestige to be true")
	}

	// Press 'y' to confirm prestige
	modelResult, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	updatedModel = modelResult.(Model)

	if updatedModel.ConfirmingPrestige {
		t.Errorf("Expected ConfirmingPrestige to be false")
	}

	if updatedModel.Engine.Player.KaPoints <= 0 {
		t.Errorf("Expected KaPoints > 0, got %d", updatedModel.Engine.Player.KaPoints)
	}

	// Try prestige when not enough
	player2 := game.NewPlayer()
	engine2 := game.NewEngine(player2, registry)
	model2 := NewModel(engine2)
	modelResult2, _ := model2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	updatedModel2 := modelResult2.(Model)

	if updatedModel2.ConfirmingPrestige {
		t.Errorf("Expected ConfirmingPrestige to be false when bits < 500k")
	}

	lastLog := updatedModel2.Logs[len(updatedModel2.Logs)-1]
	if lastLog != "[WARN] Not enough total bits for Beam Rescue (Min: 500k)." {
		t.Errorf("Expected insufficient bits log, got %s", lastLog)
	}
}

func TestHandleKey_BuyUpgrade(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	player.Bits = 1000

	// Set cursor to the first upgrade
	model.Cursor = 0
	upgradeID := engine.Registry.Order[model.Cursor]

	// Press 'enter'
	modelResult, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	updatedModel := modelResult.(Model)

	if updatedModel.Engine.Player.UpgradesOwned[upgradeID] != 1 {
		t.Errorf("Expected upgrade to be purchased")
	}
}

func TestHandleKey_Quit(t *testing.T) {
	player := game.NewPlayer()
	registry := game.NewRegistry()
	engine := game.NewEngine(player, registry)
	model := NewModel(engine)

	modelResult, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	_, ok := modelResult.(Model)
	if !ok {
		t.Errorf("Expected model to be of type Model")
	}

	if cmd == nil {
		t.Errorf("Expected non-nil tea.Cmd for quit")
	}
}
