package main

import (
	"fmt"
	"os"
	"time"

	"ka-bits/pkg/game"
	"ka-bits/pkg/storage"
)

func main() {
	player := game.NewPlayer()
	start := time.Now()
	storage.Save(player)
	elapsed := time.Since(start)
	fmt.Printf("Initial Save: %s\n", elapsed)

	start = time.Now()
	storage.Save(player)
	elapsed = time.Since(start)
	fmt.Printf("Save with Backup: %s\n", elapsed)
	os.Remove("save.json")
	os.Remove("save.json.bak")
}
