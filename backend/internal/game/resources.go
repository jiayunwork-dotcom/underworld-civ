package game

type Resources struct {
	Stone       int `json:"stone"`
	Metal       int `json:"metal"`
	GlowMushroom int `json:"glow_mushroom"`
	Water       int `json:"water"`
	MagicCrystal int `json:"magic_crystal"`
	FossilFuel  int `json:"fossil_fuel"`
}

type ResourceStorage struct {
	Resources
	Capacity int `json:"capacity"`
}

type ResourceProduction struct {
	Stone       int `json:"stone"`
	Metal       int `json:"metal"`
	GlowMushroom int `json:"glow_mushroom"`
	Water       int `json:"water"`
	MagicCrystal int `json:"magic_crystal"`
	FossilFuel  int `json:"fossil_fuel"`
}
