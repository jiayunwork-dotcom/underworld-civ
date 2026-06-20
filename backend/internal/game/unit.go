package game

type UnitType string

const (
	UnitSapper       UnitType = "sapper"
	UnitInfantry     UnitType = "infantry"
	UnitArcher       UnitType = "archer"
	UnitSiegeRam     UnitType = "siege_ram"
	UnitIronGuard    UnitType = "iron_guard"
	UnitSporeGrenadier UnitType = "spore_grenadier"
	UnitShadowArcher   UnitType = "shadow_archer"
	UnitRockGiant    UnitType = "rock_giant"
	UnitBurrowBomber UnitType = "burrow_bomber"
)

type Unit struct {
	ID        string   `json:"id"`
	Type      UnitType `json:"type"`
	Owner     string   `json:"owner"`
	HP        int      `json:"hp"`
	MaxHP     int      `json:"max_hp"`
	Attack    int      `json:"attack"`
	Defense   int      `json:"defense"`
	Range     int      `json:"range"`
	Speed     int      `json:"speed"`
	Moved     bool     `json:"moved"`
	Attacked  bool     `json:"attacked"`
	Layer     int      `json:"layer"`
	Coord     HexCoord `json:"coord"`
}

type UnitInfo struct {
	Type         UnitType  `json:"type"`
	Name         string    `json:"name"`
	Cost         Resources `json:"cost"`
	HP           int       `json:"hp"`
	Attack       int       `json:"attack"`
	Defense      int       `json:"defense"`
	Range        int       `json:"range"`
	Speed        int       `json:"speed"`
	Description  string    `json:"description"`
	RaceSpecific string    `json:"race_specific"`
}

var UnitDefs = map[UnitType]UnitInfo{
	UnitSapper: {
		Type:        UnitSapper,
		Name:        "工兵",
		Cost:        Resources{Stone: 10, GlowMushroom: 5},
		HP:          30,
		Attack:      5,
		Defense:     3,
		Range:       1,
		Speed:       2,
		Description: "可挖掘，简单战斗",
	},
	UnitInfantry: {
		Type:        UnitInfantry,
		Name:        "步兵",
		Cost:        Resources{Metal: 15, GlowMushroom: 10},
		HP:          60,
		Attack:      15,
		Defense:     10,
		Range:       1,
		Speed:       2,
		Description: "标准战斗单位",
	},
	UnitArcher: {
		Type:        UnitArcher,
		Name:        "弓手",
		Cost:        Resources{Metal: 20, GlowMushroom: 10},
		HP:          40,
		Attack:      20,
		Defense:     5,
		Range:       2,
		Speed:       2,
		Description: "射程2格，隧道弯折挡视线",
	},
	UnitSiegeRam: {
		Type:        UnitSiegeRam,
		Name:        "攻城锤",
		Cost:        Resources{Metal: 40, Stone: 30, FossilFuel: 10},
		HP:          100,
		Attack:      40,
		Defense:     15,
		Range:       1,
		Speed:       1,
		Description: "对建筑伤害×3",
	},
	UnitIronGuard: {
		Type:         UnitIronGuard,
		Name:         "铁甲卫士",
		Cost:         Resources{Metal: 50, GlowMushroom: 20},
		HP:           120,
		Attack:       25,
		Defense:      25,
		Range:        1,
		Speed:        1,
		Description:  "矮人独有，高防御重甲单位",
		RaceSpecific: "dwarf",
	},
	UnitSporeGrenadier: {
		Type:         UnitSporeGrenadier,
		Name:         "孢子投弹兵",
		Cost:         Resources{GlowMushroom: 40, MagicCrystal: 10},
		HP:           50,
		Attack:       30,
		Defense:      5,
		Range:        2,
		Speed:        2,
		Description:  "蘑菇人独有，范围孢子攻击",
		RaceSpecific: "mushroom",
	},
	UnitShadowArcher: {
		Type:         UnitShadowArcher,
		Name:         "暗影射手",
		Cost:         Resources{Metal: 30, MagicCrystal: 20},
		HP:           45,
		Attack:       35,
		Defense:      8,
		Range:        3,
		Speed:        3,
		Description:  "洞穴精灵独有，超远射程",
		RaceSpecific: "elf",
	},
	UnitRockGiant: {
		Type:         UnitRockGiant,
		Name:         "岩石巨人",
		Cost:         Resources{Stone: 80, Metal: 40},
		HP:           200,
		Attack:       35,
		Defense:      30,
		Range:        1,
		Speed:        1,
		Description:  "石像族独有，超高生命值",
		RaceSpecific: "golem",
	},
	UnitBurrowBomber: {
		Type:         UnitBurrowBomber,
		Name:         "钻地爆虫",
		Cost:         Resources{FossilFuel: 30, GlowMushroom: 20},
		HP:           35,
		Attack:       50,
		Defense:      2,
		Range:        1,
		Speed:        3,
		Description:  "深渊虫族独有，自爆单位，可钻地",
		RaceSpecific: "zerg",
	},
}
