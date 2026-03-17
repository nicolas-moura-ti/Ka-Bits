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
