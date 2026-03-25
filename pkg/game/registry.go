package game

import "math"

type UpgradeRegistry struct {
	Upgrades map[string]Upgrade
	Order    []string
}

func NewRegistry() *UpgradeRegistry {
	r := &UpgradeRegistry{
		Upgrades: make(map[string]Upgrade),
		Order:    []string{},
	}

	for _, u := range InitialUpgrades {
		r.Register(u)
	}

	return r
}

func (r *UpgradeRegistry) Register(u Upgrade) {
	r.Upgrades[u.ID] = u
	r.Order = append(r.Order, u.ID)
}

func (r *UpgradeRegistry) Get(id string) (Upgrade, bool) {
	u, ok := r.Upgrades[id]
	return u, ok
}

var InitialUpgrades = []Upgrade{
	{
		ID:          "terminal_gilead",
		Name:        "Gilead Terminal",
		Description: "A rusty terminal pulsing with the rhythm of Ka. It smells like old parchment.",
		BaseCost:    100, // Increased from 50
		BaseBPS:     0.1, // Decreased from 0.5
		Type:        "Hardware",
	},
	{
		ID:          "servidor_mid_world",
		Name:        "Mid-World Server",
		Description: "Servers located in the thin spots of reality. Sometimes you hear voices in the fans.",
		BaseCost:    1000, // Increased from 500
		BaseBPS:     0.5,  // Decreased from 2
		Type:        "Hardware",
	},
	{
		ID:          "lobstros_optimizer",
		Name:        "Lobstros Optimizer",
		Description: "Highly efficient scripts. Did-a-chick? Dum-a-chum? Better not touch the crab-claws.",
		BaseCost:    5000, // Increased from 2500
		BaseBPS:     2,    // Decreased from 10
		Type:        "Software",
	},
	{
		ID:          "blaine_engine",
		Name:        "Blaine's Logic Engine",
		Description: "A mono-train AI that solves riddles for processing power. Don't mention silly stuff.",
		BaseCost:    20000, // Increased from 7500
		BaseBPS:     5,     // Decreased from 25
		Type:        "Software",
	},
	{
		ID:          "quantum_rose",
		Name:        "Quantum Rose",
		Description: "The focal point of existence, blooming in binary. It sings a song of 19.",
		BaseCost:    100000, // Increased from 20000
		BaseBPS:     15,     // Decreased from 75
		Type:        "Cosmic",
	},
	{
		ID:          "crimson_king_proxy",
		Name:        "Crimson King Proxy",
		Description: "Routing your data through the Red King's own network. Risky, but powerful.",
		BaseCost:    500000, // Increased from 100000
		BaseBPS:     50,     // Decreased from 250
		Type:        "Cosmic",
	},
}

func (p *Player) CalculateBPS(r *UpgradeRegistry) float64 {
	if p.CacheValid {
		return p.CachedBPS
	}

	bps := 0.0
	bonusMultiplier := 1.0
	totalUpgrades := 0

	for id, count := range p.UpgradesOwned {
		totalUpgrades += count

		// Special "Sincronicidade" bonus for 19 or 99
		if count == 19 || count == 99 {
			bonusMultiplier = 1.19
		}

		upgrade, ok := r.Get(id)
		if ok {
			bps += upgrade.BaseBPS * float64(count)
		}
	}

	// Prestige multiplier: 5% more per Ka-Point
	prestigeMultiplier := 1.0 + (float64(p.KaPoints) * 0.05)

	p.CachedBPS = bps * bonusMultiplier * prestigeMultiplier
	p.CachedTotalUpgrades = totalUpgrades
	p.CacheValid = true

	return p.CachedBPS
}

func (p *Player) GetTotalUpgrades(r *UpgradeRegistry) int {
	p.CalculateBPS(r)
	return p.CachedTotalUpgrades
}

func CalculateUpgradeCost(baseCost float64, owned int) float64 {
	// Brutal scaling (2.2 instead of 1.5)
	// Every level more than doubles the cost.
	return baseCost * math.Pow(2.2, float64(owned))
}

var RandomLogs = []string{
	"[INFO] Backup performed on Lud's servers.",
	"[WARN] The Crimson King tried to intercept data packets.",
	"[ERROR] Critical Failure: The world has moved on.",
	"[INFO] Lobstros detected on the network. Did-a-chick?",
	"[INFO] Blaine is a pain, but he sure can calculate.",
	"[WARN] Thin spot detected in the database cluster.",
	"[INFO] Feeding the Rose with fresh binary data.",
	"[INFO] Slow mutants cleared from the cache.",
	"[WARN] Oy the Brave found a bug in the routing table.",
	"[INFO] Shardik.exe is consuming too much memory.",
	"[INFO] Re-aligning the Beam for better throughput.",
	"[WARN] Tak is shouting at the server room door.",
	"[INFO] Low Men in Yellow Coats are monitoring your traffic.",
	"[INFO] Sombra Corporation attempted a hostile takeover.",
	"[WARN] The Calla is calling for a system update.",
	"[INFO] Commala-come-come, the bits are on the run.",
}
