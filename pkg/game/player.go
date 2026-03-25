package game

import (
	"time"
)

type Player struct {
	Version       string             `json:"version"`
	Bits          float64            `json:"bits"`
	TotalBitsEver float64            `json:"total_bits_ever"` // Track total for prestige
	KaPoints      int                `json:"ka_points"`
	Resources     map[string]float64 `json:"resources"`
	LastUpdate    time.Time          `json:"last_update"`
	UpgradesOwned map[string]int     `json:"upgrades"`

	CachedBPS           float64 `json:"-"`
	CachedTotalUpgrades int     `json:"-"`
	CacheValid          bool    `json:"-"`
}

func NewPlayer() *Player {
	p := &Player{
		Version: "0.1.0",
	}
	p.resetRunState()
	p.KaPoints = 0
	return p
}

func (p *Player) resetRunState() {
	p.Bits = 0
	p.TotalBitsEver = 0
	p.Resources = make(map[string]float64)
	p.LastUpdate = time.Now()
	p.UpgradesOwned = make(map[string]int)
	p.InvalidateCache()
}

func (p *Player) InvalidateCache() {
	p.CacheValid = false
}

func (p *Player) Reset() {
	p.resetRunState()
	p.KaPoints = 0
}

func (p *Player) BeamRescue(gainedPoints int) {
	p.KaPoints += gainedPoints
	p.resetRunState()
}
