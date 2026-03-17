package storage

import (
	"encoding/json"
	"ka-bits/pkg/game"
	"os"
)

const SaveFilePath = "save.json"

func Save(player *game.Player) error {
	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}

	// Create backup before writing
	if _, err := os.Stat(SaveFilePath); err == nil {
		backupData, _ := os.ReadFile(SaveFilePath)
		_ = os.WriteFile(SaveFilePath+".bak", backupData, 0644)
	}

	return os.WriteFile(SaveFilePath, data, 0644)
}

func Load() (*game.Player, error) {
	if _, err := os.Stat(SaveFilePath); os.IsNotExist(err) {
		// Try to load from backup if primary is missing (unexpected)
		if _, errBak := os.Stat(SaveFilePath + ".bak"); errBak == nil {
			return loadFile(SaveFilePath + ".bak")
		}
		return game.NewPlayer(), nil
	}

	return loadFile(SaveFilePath)
}

func loadFile(path string) (*game.Player, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var player game.Player
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, err
	}

	return &player, nil
}
