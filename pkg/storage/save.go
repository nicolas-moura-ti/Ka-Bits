package storage

import (
	"encoding/json"
	"ka-bits/pkg/game"
	"os"
	"sync"
)

const SaveFilePath = "save.json"

var saveMutex sync.Mutex

func Save(player *game.Player) error {
	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}

	return WriteData(data)
}

func WriteData(data []byte) error {
	saveMutex.Lock()
	defer saveMutex.Unlock()

	// Create backup before writing
	if _, err := os.Stat(SaveFilePath); err == nil {
		backupData, err := os.ReadFile(SaveFilePath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(SaveFilePath+".bak", backupData, 0600); err != nil {
			return err
		}
	}

	return os.WriteFile(SaveFilePath, data, 0600)
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

	player.InvalidateCache()

	return &player, nil
}
