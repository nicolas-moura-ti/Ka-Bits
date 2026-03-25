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

func TestSavePermissions(t *testing.T) {
	player := game.NewPlayer()

	// Initial save should create SaveFilePath
	err := Save(player)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Second save should create SaveFilePath.bak
	err = Save(player)
	if err != nil {
		t.Fatalf("Failed to save second time (to trigger backup): %v", err)
	}

	// Check SaveFilePath permissions
	info, err := os.Stat(SaveFilePath)
	if err != nil {
		t.Fatalf("Failed to stat %s: %v", SaveFilePath, err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected %s to have permissions 0600, got %v", SaveFilePath, info.Mode().Perm())
	}

	// Check SaveFilePath.bak permissions
	infoBak, err := os.Stat(SaveFilePath + ".bak")
	if err != nil {
		t.Fatalf("Failed to stat %s.bak: %v", SaveFilePath, err)
	}
	if infoBak.Mode().Perm() != 0600 {
		t.Errorf("Expected %s.bak to have permissions 0600, got %v", SaveFilePath, infoBak.Mode().Perm())
	}

	// Clean up
	os.Remove(SaveFilePath)
	os.Remove(SaveFilePath + ".bak")
}
