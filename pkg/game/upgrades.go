package game

type Upgrade struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BaseCost    float64 `json:"base_cost"`
	BaseBPS     float64 `json:"base_bps"`
	Type        string  `json:"type"` // "hardware", "software", "cosmic"
}
