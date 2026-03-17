package game

import (
	"fmt"
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
	e.Player.LastUpdate = now
	
	return earnings, offlineTime
}

func (e *Engine) Update(delta time.Duration) {
	bps := e.Player.CalculateBPS(e.Registry)
	e.Player.Bits += bps * delta.Seconds()
	e.Player.LastUpdate = time.Now()
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
		return true, fmt.Sprintf("You acquired: %s!", upgrade.Name)
	}

	return false, fmt.Sprintf("Not enough bits for %s, gunslinger!", upgrade.Name)
}
