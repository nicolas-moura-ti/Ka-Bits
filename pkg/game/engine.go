package game

import (
	"fmt"
	"math"
	"time"
)

type Engine struct {
	Player   *Player
	Registry *UpgradeRegistry
}

func NewEngine(player *Player, registry *UpgradeRegistry) *Engine {
	return &Engine{
		Player:   player,
		Registry: registry,
	}
}

func (e *Engine) ProcessOfflineEarnings() (float64, time.Duration) {
	now := time.Now()
	offlineTime := now.Sub(e.Player.LastUpdate)

	if offlineTime < time.Second {
		return 0, 0
	}

	bps := e.Player.CalculateBPS(e.Registry)
	earnings := offlineTime.Seconds() * bps * 0.75

	e.Player.Bits += earnings
	e.Player.TotalBitsEver += earnings
	e.Player.LastUpdate = now

	return earnings, offlineTime
}

func (e *Engine) Update(delta time.Duration) {
	bps := e.Player.CalculateBPS(e.Registry)
	generated := bps * delta.Seconds()
	e.Player.Bits += generated
	e.Player.TotalBitsEver += generated
	e.Player.LastUpdate = time.Now()
}

func (e *Engine) CalculatePrestigeGain() int {
	// Formula: sqrt(TotalBits / 500,000)
	// You need at least 500k total bits to start gaining points.
	if e.Player.TotalBitsEver < 500000 {
		return 0
	}
	points := math.Sqrt(e.Player.TotalBitsEver / 500000)
	return int(points)
}

func (e *Engine) TryBuyUpgrade(id string) (bool, string) {
	upgrade, ok := e.Registry.Get(id)
	if !ok {
		return false, "Upgrade not found."
	}

	owned := e.Player.UpgradesOwned[id]
	cost := CalculateUpgradeCost(upgrade.BaseCost, owned)

	if e.Player.Bits >= cost {
		e.Player.Bits -= cost
		e.Player.UpgradesOwned[id]++
		e.Player.InvalidateCache()
		return true, fmt.Sprintf("You acquired: %s!", upgrade.Name)
	}

	return false, fmt.Sprintf("Not enough bits for %s, gunslinger!", upgrade.Name)
}
