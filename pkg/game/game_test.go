package game

import (
	"testing"
	"time"
)

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

func TestPlayer_Reset_Consolidation(t *testing.T) {
	p := NewPlayer()
	p.Bits = 100
	p.TotalBitsEver = 500
	p.KaPoints = 10
	p.Resources["ore"] = 5
	p.UpgradesOwned["test"] = 1

	// Test Reset Total
	p.Reset()
	if p.Bits != 0 || p.TotalBitsEver != 0 || p.KaPoints != 0 || len(p.Resources) != 0 || len(p.UpgradesOwned) != 0 {
		t.Errorf("Reset did not clear all fields: %+v", p)
	}

	// Setup for BeamRescue (Prestige)
	p.Bits = 100
	p.TotalBitsEver = 500
	p.KaPoints = 10
	p.Resources["ore"] = 5
	p.UpgradesOwned["test"] = 1

	// Test BeamRescue
	p.BeamRescue(5) // Ganha 5 KaPoints e reseta o resto
	if p.Bits != 0 || p.TotalBitsEver != 0 || p.KaPoints != 15 || len(p.Resources) != 0 || len(p.UpgradesOwned) != 0 {
		t.Errorf("BeamRescue did not clear/update fields correctly: %+v", p)
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

func TestProcessOfflineEarnings(t *testing.T) {
	registry := NewRegistry()

	t.Run("Normal offline earnings", func(t *testing.T) {
		player := NewPlayer()
		engine := NewEngine(player, registry)

		// Set BPS to a known value by adding an upgrade
		player.UpgradesOwned["terminal_gilead"] = 10 // 0.1 * 10 = 1.0 BPS

		// Simulate being offline for exactly 1 hour
		offlineDuration := time.Hour
		player.LastUpdate = time.Now().Add(-offlineDuration)

		earnings, duration := engine.ProcessOfflineEarnings()

		expectedEarnings := offlineDuration.Seconds() * 1.0 * 0.75

		// Due to small time differences between time.Now() calls, we check if it's close
		if earnings < expectedEarnings*0.99 || earnings > expectedEarnings*1.01 {
			t.Errorf("Expected earnings around %f, got %f", expectedEarnings, earnings)
		}

		if duration < offlineDuration*99/100 || duration > offlineDuration*101/100 {
			t.Errorf("Expected duration around %v, got %v", offlineDuration, duration)
		}

		if player.Bits != earnings {
			t.Errorf("Expected player Bits to be %f, got %f", earnings, player.Bits)
		}

		if player.TotalBitsEver != earnings {
			t.Errorf("Expected player TotalBitsEver to be %f, got %f", earnings, player.TotalBitsEver)
		}

		// Ensure LastUpdate was updated to near time.Now()
		if time.Since(player.LastUpdate) > time.Second {
			t.Errorf("LastUpdate was not updated correctly")
		}
	})

	t.Run("Short offline time", func(t *testing.T) {
		player := NewPlayer()
		engine := NewEngine(player, registry)

		player.UpgradesOwned["terminal_gilead"] = 10 // 1.0 BPS

		// Simulate being offline for 0.5 seconds
		player.LastUpdate = time.Now().Add(-500 * time.Millisecond)

		earnings, duration := engine.ProcessOfflineEarnings()

		if earnings != 0 {
			t.Errorf("Expected 0 earnings for <1s offline time, got %f", earnings)
		}

		if duration != 0 {
			t.Errorf("Expected 0 duration for <1s offline time, got %v", duration)
		}

		if player.Bits != 0 {
			t.Errorf("Expected player Bits to be 0, got %f", player.Bits)
		}
	})

	t.Run("Zero BPS", func(t *testing.T) {
		player := NewPlayer()
		engine := NewEngine(player, registry)

		// 0 BPS

		// Simulate being offline for 1 hour
		offlineDuration := time.Hour
		player.LastUpdate = time.Now().Add(-offlineDuration)

		earnings, duration := engine.ProcessOfflineEarnings()

		if earnings != 0 {
			t.Errorf("Expected 0 earnings for 0 BPS, got %f", earnings)
		}

		if duration < offlineDuration*99/100 || duration > offlineDuration*101/100 {
			t.Errorf("Expected duration around %v, got %v", offlineDuration, duration)
		}

		if player.Bits != 0 {
			t.Errorf("Expected player Bits to be 0, got %f", player.Bits)
		}

		// Ensure LastUpdate was updated
		if time.Since(player.LastUpdate) > time.Second {
			t.Errorf("LastUpdate was not updated correctly")
		}
	})
}