package game

type BuildingType string

const (
	BuildingLivingQuarters BuildingType = "living_quarters"
	BuildingWorkshop       BuildingType = "workshop"
	BuildingFungusFarm     BuildingType = "fungus_farm"
	BuildingSmelter        BuildingType = "smelter"
	BuildingAcademy        BuildingType = "academy"
	BuildingAltar          BuildingType = "altar"
	BuildingWarehouse      BuildingType = "warehouse"
	BuildingWatchtower     BuildingType = "watchtower"
	BuildingWall           BuildingType = "wall"
	BuildingElevator       BuildingType = "elevator"
	BuildingForgeShrine    BuildingType = "forge_shrine"
	BuildingSporeNetwork   BuildingType = "spore_network"
	BuildingCrystalTower   BuildingType = "crystal_tower"
	BuildingPetrifiedWall  BuildingType = "petrified_wall"
	BuildingHiveRift       BuildingType = "hive_rift"
	BuildingMainBase       BuildingType = "main_base"
)

type Building struct {
	Type      BuildingType `json:"type"`
	HP        int          `json:"hp"`
	MaxHP     int          `json:"max_hp"`
	Level     int          `json:"level"`
	Owner     string       `json:"owner"`
	Completed bool         `json:"completed"`
	Progress  int          `json:"progress"`
	BuildTime int          `json:"build_time"`
}

type BuildingInfo struct {
	Type         BuildingType `json:"type"`
	Name         string       `json:"name"`
	Cost         Resources    `json:"cost"`
	BuildTime    int          `json:"build_time"`
	HP           int          `json:"hp"`
	Description  string       `json:"description"`
	RaceSpecific string       `json:"race_specific"`
}

var BuildingDefs = map[BuildingType]BuildingInfo{
	BuildingMainBase: {
		Type:        BuildingMainBase,
		Name:        "主基地",
		Cost:        Resources{},
		BuildTime:   0,
		HP:          500,
		Description: "玩家的主基地，被摧毁则失败",
	},
	BuildingLivingQuarters: {
		Type:        BuildingLivingQuarters,
		Name:        "居住窟",
		Cost:        Resources{Stone: 20},
		BuildTime:   2,
		HP:          100,
		Description: "+5人口上限",
	},
	BuildingWorkshop: {
		Type:        BuildingWorkshop,
		Name:        "工坊",
		Cost:        Resources{Stone: 30, Metal: 10},
		BuildTime:   3,
		HP:          80,
		Description: "+2金属产出",
	},
	BuildingFungusFarm: {
		Type:        BuildingFungusFarm,
		Name:        "农菌场",
		Cost:        Resources{Stone: 15},
		BuildTime:   2,
		HP:          60,
		Description: "+3发光菌产出",
	},
	BuildingSmelter: {
		Type:        BuildingSmelter,
		Name:        "冶炼炉",
		Cost:        Resources{Stone: 40, Metal: 20},
		BuildTime:   4,
		HP:          120,
		Description: "解锁金属加工，金属产出效率+50%",
	},
	BuildingAcademy: {
		Type:        BuildingAcademy,
		Name:        "学堂",
		Cost:        Resources{Stone: 50, Metal: 30},
		BuildTime:   5,
		HP:          80,
		Description: "+1研究点/回合",
	},
	BuildingAltar: {
		Type:        BuildingAltar,
		Name:        "祭坛",
		Cost:        Resources{Stone: 60, MagicCrystal: 10},
		BuildTime:   6,
		HP:          100,
		Description: "每10回合触发随机正面事件",
	},
	BuildingWarehouse: {
		Type:        BuildingWarehouse,
		Name:        "仓储洞",
		Cost:        Resources{Stone: 25},
		BuildTime:   2,
		HP:          150,
		Description: "+100资源存储上限",
	},
	BuildingWatchtower: {
		Type:        BuildingWatchtower,
		Name:        "哨塔",
		Cost:        Resources{Stone: 35, Metal: 15},
		BuildTime:   3,
		HP:          70,
		Description: "+3视野范围",
	},
	BuildingWall: {
		Type:        BuildingWall,
		Name:        "城墙",
		Cost:        Resources{Stone: 50},
		BuildTime:   4,
		HP:          200,
		Description: "格子防御+50%",
	},
	BuildingElevator: {
		Type:        BuildingElevator,
		Name:        "电梯",
		Cost:        Resources{Stone: 50, Metal: 30, FossilFuel: 20},
		BuildTime:   6,
		HP:          100,
		Description: "连接上下层，允许单位和资源跨层移动",
	},
	BuildingForgeShrine: {
		Type:         BuildingForgeShrine,
		Name:         "熔炉圣殿",
		Cost:         Resources{Stone: 80, Metal: 60, MagicCrystal: 20},
		BuildTime:    8,
		HP:           180,
		Description:  "矮人独有，解锁高级冶炼，金属产出+100%",
		RaceSpecific: "dwarf",
	},
	BuildingSporeNetwork: {
		Type:         BuildingSporeNetwork,
		Name:         "孢子网络",
		Cost:         Resources{Stone: 40, GlowMushroom: 30},
		BuildTime:    5,
		HP:           80,
		Description:  "蘑菇人独有，建筑间即时资源传送",
		RaceSpecific: "mushroom",
	},
	BuildingCrystalTower: {
		Type:         BuildingCrystalTower,
		Name:         "水晶共鸣塔",
		Cost:         Resources{Stone: 60, MagicCrystal: 40},
		BuildTime:    7,
		HP:           100,
		Description:  "洞穴精灵独有，提供跨层感知",
		RaceSpecific: "elf",
	},
	BuildingPetrifiedWall: {
		Type:         BuildingPetrifiedWall,
		Name:         "石化长城",
		Cost:         Resources{Stone: 100, MagicCrystal: 30},
		BuildTime:    10,
		HP:           500,
		Description:  "石像族独有，将整行格子变为不可挖掘",
		RaceSpecific: "golem",
	},
	BuildingHiveRift: {
		Type:         BuildingHiveRift,
		Name:         "虫巢裂隙",
		Cost:         Resources{Stone: 70, FossilFuel: 40},
		BuildTime:    6,
		HP:           120,
		Description:  "深渊虫族独有，允许部队地下穿行绕过障碍",
		RaceSpecific: "zerg",
	},
}
