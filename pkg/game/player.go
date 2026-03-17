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
}

func NewPlayer() *Player {
	return &Player{
		Version:       "0.1.0",
		Bits:          0,
		TotalBitsEver: 0,
		KaPoints:      0,
		Resources:     make(map[string]float64),
		LastUpdate:    time.Now(),
		UpgradesOwned: make(map[string]int),
	}
}

func (p *Player) Reset() {
	p.Bits = 0
	p.TotalBitsEver = 0
	p.KaPoints = 0
	p.Resources = make(map[string]float64)
	p.LastUpdate = time.Now()
	p.UpgradesOwned = make(map[string]int)
}

func (p *Player) BeamRescue(gainedPoints int) {
	p.KaPoints += gainedPoints
	p.Bits = 0
	// Keep TotalBitsEver if you want it to be cumulative across resets, 
	// or reset it if you want points to be based on the current run.
	// Let's reset it for "per-run" point calculation.
	p.TotalBitsEver = 0 
	p.Resources = make(map[string]float64)
	p.LastUpdate = time.Now()
	p.UpgradesOwned = make(map[string]int)
}
