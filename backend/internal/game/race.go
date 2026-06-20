package game

type Race string

const (
	RaceDwarf    Race = "dwarf"
	RaceMushroom Race = "mushroom"
	RaceElf      Race = "elf"
	RaceGolem    Race = "golem"
	RaceZerg     Race = "zerg"
)

type RaceInfo struct {
	Race        Race   `json:"race"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UniqueUnit  UnitType `json:"unique_unit"`
	UniqueBuilding BuildingType `json:"unique_building"`
	Bonuses     map[string]float64 `json:"bonuses"`
}

var RaceDefs = map[Race]RaceInfo{
	RaceDwarf: {
		Race:           RaceDwarf,
		Name:           "矮人",
		Description:    "挖掘速度+30%，独有熔炉圣殿解锁高级冶炼",
		UniqueUnit:     UnitIronGuard,
		UniqueBuilding: BuildingForgeShrine,
		Bonuses: map[string]float64{
			"mining_speed": 1.3,
			"metal_production": 1.0,
		},
	},
	RaceMushroom: {
		Race:           RaceMushroom,
		Name:           "蘑菇人",
		Description:    "农菌产出+50%，独有孢子网络实现建筑间即时资源传送",
		UniqueUnit:     UnitSporeGrenadier,
		UniqueBuilding: BuildingSporeNetwork,
		Bonuses: map[string]float64{
			"fungus_production": 1.5,
		},
	},
	RaceElf: {
		Race:           RaceElf,
		Name:           "洞穴精灵",
		Description:    "视野范围+2格，独有水晶共鸣塔提供跨层感知",
		UniqueUnit:     UnitShadowArcher,
		UniqueBuilding: BuildingCrystalTower,
		Bonuses: map[string]float64{
			"vision_range": 2.0,
		},
	},
	RaceGolem: {
		Race:           RaceGolem,
		Name:           "石像族",
		Description:    "建筑耐久+40%，独有石化长城可将整行格子变为不可挖掘",
		UniqueUnit:     UnitRockGiant,
		UniqueBuilding: BuildingPetrifiedWall,
		Bonuses: map[string]float64{
			"building_hp": 1.4,
		},
	},
	RaceZerg: {
		Race:           RaceZerg,
		Name:           "深渊虫族",
		Description:    "可在熔岩层无伤移动，独有虫巢裂隙允许部队地下穿行绕过障碍",
		UniqueUnit:     UnitBurrowBomber,
		UniqueBuilding: BuildingHiveRift,
		Bonuses: map[string]float64{
			"lava_immune": 1.0,
			"burrow_move": 1.0,
		},
	},
}
