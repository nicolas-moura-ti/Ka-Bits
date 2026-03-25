package storage

import (
	"encoding/json"
	"ka-bits/pkg/game"
	"os"
	"testing"
)

func BenchmarkJSONMarshal(b *testing.B) {
	player := game.NewPlayer()
	player.Bits = 1000
	player.UpgradesOwned["test"] = 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.MarshalIndent(player, "", "  ")
	}
}

func BenchmarkDiskIO(b *testing.B) {
	player := game.NewPlayer()
	player.Bits = 1000
	player.UpgradesOwned["test"] = 50
	data, _ := json.MarshalIndent(player, "", "  ")
	_ = os.WriteFile(SaveFilePath, data, 0600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := os.Stat(SaveFilePath); err == nil {
			backupData, _ := os.ReadFile(SaveFilePath)
			_ = os.WriteFile(SaveFilePath+".bak", backupData, 0600)
		}
		_ = os.WriteFile(SaveFilePath, data, 0600)
	}
	os.Remove(SaveFilePath)
	os.Remove(SaveFilePath + ".bak")
}
