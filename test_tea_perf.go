package main

import (
	"fmt"
	"time"

	"ka-bits/pkg/game"
	"ka-bits/pkg/storage"
)

func main() {
	player := game.NewPlayer()

	start := time.Now()
	for i := 0; i < 100; i++ {
		storage.Save(player)
	}
	elapsed := time.Since(start)
	fmt.Printf("100 Synchronous Saves: %s\n", elapsed)
}
