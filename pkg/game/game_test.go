package game

import "testing"

func TestCalculateBPS(t *testing.T) {
	player := NewPlayer()
	registry := NewRegistry()
	player.UpgradesOwned["terminal_gilead"] = 2    // 0.1 BPS * 2 = 0.2
	player.UpgradesOwned["servidor_mid_world"] = 1 // 0.5 BPS * 1 = 0.5
	                                               // Total = 0.7 BPS

	bps := player.CalculateBPS(registry)
	if bps != 0.7 {
		t.Errorf("Expected 0.7 BPS, got %f", bps)
	}
}

func TestTryBuyUpgrade(t *testing.T) {
	player := NewPlayer()
	registry := NewRegistry()
	engine := NewEngine(player, registry)
	player.Bits = 90

	// Gilead Terminal base cost is 100.
	success, _ := engine.TryBuyUpgrade("terminal_gilead")
	if success {
		t.Error("Should have failed (insufficient bits)")
	}

	player.Bits = 110
	success, _ = engine.TryBuyUpgrade("terminal_gilead")
	if !success {
		t.Error("Should have worked (sufficient bits)")
	}

	if player.UpgradesOwned["terminal_gilead"] != 1 {
		t.Errorf("Expected 1 upgrade, got %d", player.UpgradesOwned["terminal_gilead"])
	}

	if player.Bits >= 20 {
		t.Error("Bits were not deducted correctly")
	}
}

func TestCalculatePrestigeGain(t *testing.T) {
	tests := []struct {
		name          string
		totalBitsEver float64
		want          int
	}{
		{"Zero bits", 0, 0},
		{"Just below threshold", 499999, 0},
		{"Exactly threshold", 500000, 1},
		{"One million bits (sqrt(2) approx 1.41)", 1000000, 1},
		{"Two million bits (sqrt(4) = 2)", 2000000, 2},
		{"4.5 million bits (sqrt(9) = 3)", 4500000, 3},
		{"Very large amount", 500000000, 31}, // sqrt(1000) approx 31.62
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer()
			player.TotalBitsEver = tt.totalBitsEver
			engine := &Engine{Player: player}
			if got := engine.CalculatePrestigeGain(); got != tt.want {
				t.Errorf("CalculatePrestigeGain() with %f bits = %v, want %v", tt.totalBitsEver, got, tt.want)
			}
		})
	}
}
