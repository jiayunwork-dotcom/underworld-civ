package game

type HexCoord struct {
	Q int `json:"q"`
	R int `json:"r"`
}

type HexCell struct {
	Coord         HexCoord `json:"coord"`
	Layer         int      `json:"layer"`
	Discovered    bool     `json:"discovered"`
	Owner         string   `json:"owner"`
	RockHardness  int      `json:"rock_hardness"`
	WaterContent  float64  `json:"water_content"`
	MineralType   string   `json:"mineral_type"`
	IsWall        bool     `json:"is_wall"`
	IsShaft       bool     `json:"is_shaft"`
	Building      *Building `json:"building"`
	Units         []Unit   `json:"units"`
	DefenseBonus  float64  `json:"defense_bonus"`
	Flooded       bool     `json:"flooded"`
	MiningOwner   string   `json:"mining_owner"`
	MiningProgress int     `json:"mining_progress"`
	MiningTotal   int      `json:"mining_total"`
}

type Layer struct {
	Depth    int                `json:"depth"`
	Name     string             `json:"name"`
	Width    int                `json:"width"`
	Height   int                `json:"height"`
	Cells    map[string]*HexCell `json:"cells"`
	Shafts   []HexCoord         `json:"shafts"`
}

type GameMap struct {
	Layers []*Layer `json:"layers"`
}
