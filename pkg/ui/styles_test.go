package ui

import (
	"testing"
)

func TestGetPulseColor(t *testing.T) {
	// The pulseColors array has 6 elements.
	// tick/4 % len(pulseColors) -> repeats every 4 * len(pulseColors) ticks
	// Since pulseColors is unexported, we can just use 6, or compare dynamically
	cycleLength := 4 * len(pulseColors)
	for tick := 0; tick <= 100; tick++ {
		color := GetPulseColor(tick)
		if color == "" {
			t.Errorf("GetPulseColor(%d) returned empty color", tick)
		}

		// cyclical value check
		expectedTick := tick % cycleLength
		expectedColor := GetPulseColor(expectedTick)
		if color != expectedColor {
			t.Errorf("GetPulseColor(%d) = %v, expected %v", tick, color, expectedColor)
		}
	}
}

func TestGetGlowColor(t *testing.T) {
	// The glowColors array has 6 elements.
	// tick/6 % len(glowColors) -> repeats every 6 * len(glowColors) ticks
	cycleLength := 6 * len(glowColors)
	for tick := 0; tick <= 100; tick++ {
		color := GetGlowColor(tick)
		if color == "" {
			t.Errorf("GetGlowColor(%d) returned empty color", tick)
		}

		// cyclical value check
		expectedTick := tick % cycleLength
		expectedColor := GetGlowColor(expectedTick)
		if color != expectedColor {
			t.Errorf("GetGlowColor(%d) = %v, expected %v", tick, color, expectedColor)
		}
	}
}

func TestGetGlowStyle(t *testing.T) {
	// GlowStyles has length 6.
	// tick/6 % len(GlowStyles) -> repeats every 6 * len(GlowStyles) ticks
	cycleLength := 6 * len(GlowStyles)
	for tick := 0; tick <= 100; tick++ {
		style := GetGlowStyle(tick)

		// cyclical value check
		expectedTick := tick % cycleLength
		expectedStyle := GetGlowStyle(expectedTick)
		// Instead of comparing style directly, let's compare a rendered string
		if style.Render("test") != expectedStyle.Render("test") {
			t.Errorf("GetGlowStyle(%d) rendered string = %v, expected %v", tick, style.Render("test"), expectedStyle.Render("test"))
		}
	}
}
