package storage

import (
	"ka-bits/pkg/game"
	"os"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	player := game.NewPlayer()
	player.Bits = 100
	player.UpgradesOwned["test"] = 5

	err := Save(player)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if loaded.Bits != 100 {
		t.Errorf("Expected 100 bits, got %f", loaded.Bits)
	}

	if loaded.UpgradesOwned["test"] != 5 {
		t.Errorf("Expected 5 upgrades, got %d", loaded.UpgradesOwned["test"])
	}

	// Limpar
	os.Remove(SaveFilePath)
	os.Remove(SaveFilePath + ".bak")
}

func TestSaveBackupFailure(t *testing.T) {
	// Create a directory with the same name as SaveFilePath to cause os.ReadFile to fail
	os.Mkdir(SaveFilePath, 0755)
	defer os.RemoveAll(SaveFilePath)

	player := game.NewPlayer()
	err := Save(player)

	if err == nil {
		t.Errorf("Expected error when backup fails, but got nil")
	}
}
