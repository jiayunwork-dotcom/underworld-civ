package game

type Tech struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        int       `json:"cost"`
	Category    string    `json:"category"`
	Prerequisites []string `json:"prerequisites"`
	Effects     map[string]float64 `json:"effects"`
	Researched  bool      `json:"researched"`
	Progress    int       `json:"progress"`
}

var TechDefs = []Tech{
	{ID: "basic_mining", Name: "基础采矿", Description: "挖掘速度+20%", Cost: 10, Category: "mining", Effects: map[string]float64{"mining_speed": 0.2}},
	{ID: "advanced_mining", Name: "高级采矿", Description: "挖掘速度+30%", Cost: 25, Category: "mining", Prerequisites: []string{"basic_mining"}, Effects: map[string]float64{"mining_speed": 0.3}},
	{ID: "support_beams", Name: "支撑梁技术", Description: "塌方概率减半", Cost: 20, Category: "mining", Effects: map[string]float64{"collapse_chance": 0.5}},
	{ID: "deep_drilling", Name: "深层钻探", Description: "电梯建造速度+50%", Cost: 30, Category: "mining", Prerequisites: []string{"advanced_mining"}, Effects: map[string]float64{"elevator_speed": 0.5}},

	{ID: "basic_smelting", Name: "基础冶炼", Description: "金属产出+20%", Cost: 15, Category: "industry", Effects: map[string]float64{"metal_production": 0.2}},
	{ID: "advanced_smelting", Name: "高级冶炼", Description: "金属产出+30%", Cost: 35, Category: "industry", Prerequisites: []string{"basic_smelting"}, Effects: map[string]float64{"metal_production": 0.3}},
	{ID: "stonecutting", Name: "石材加工", Description: "建筑建造速度+25%", Cost: 20, Category: "industry", Effects: map[string]float64{"build_speed": 0.25}},
	{ID: "automation", Name: "自动化", Description: "所有资源产出+15%", Cost: 50, Category: "industry", Prerequisites: []string{"advanced_smelting", "stonecutting"}, Effects: map[string]float64{"all_production": 0.15}},

	{ID: "mushroom_farming", Name: "真菌培育", Description: "发光菌产出+30%", Cost: 12, Category: "biology", Effects: map[string]float64{"fungus_production": 0.3}},
	{ID: "water_recycling", Name: "水循环系统", Description: "水消耗-20%", Cost: 18, Category: "biology", Effects: map[string]float64{"water_efficiency": 0.2}},
	{ID: "bio_luminescence", Name: "生物发光", Description: "视野范围+1格", Cost: 25, Category: "biology", Prerequisites: []string{"mushroom_farming"}, Effects: map[string]float64{"vision_range": 1.0}},

	{ID: "basic_weapons", Name: "基础武器", Description: "单位攻击力+10%", Cost: 15, Category: "military", Effects: map[string]float64{"attack": 0.1}},
	{ID: "advanced_weapons", Name: "高级武器", Description: "单位攻击力+20%", Cost: 35, Category: "military", Prerequisites: []string{"basic_weapons"}, Effects: map[string]float64{"attack": 0.2}},
	{ID: "armor", Name: "装甲技术", Description: "单位防御力+15%", Cost: 25, Category: "military", Effects: map[string]float64{"defense": 0.15}},
	{ID: "tactics", Name: "战术训练", Description: "部队士气+20%", Cost: 30, Category: "military", Prerequisites: []string{"basic_weapons"}, Effects: map[string]float64{"morale": 0.2}},
	{ID: "siege_craft", Name: "攻城工程", Description: "攻城单位伤害+30%", Cost: 40, Category: "military", Prerequisites: []string{"advanced_weapons"}, Effects: map[string]float64{"siege_damage": 0.3}},

	{ID: "magic_crystal_theory", Name: "魔晶理论", Description: "魔晶产出+25%", Cost: 20, Category: "magic", Effects: map[string]float64{"crystal_production": 0.25}},
	{ID: "crystal_resonance", Name: "水晶共振", Description: "研究速度+20%", Cost: 30, Category: "magic", Prerequisites: []string{"magic_crystal_theory"}, Effects: map[string]float64{"research_speed": 0.2}},
	{ID: "arcane_defense", Name: "奥术防御", Description: "建筑耐久+20%", Cost: 35, Category: "magic", Prerequisites: []string{"crystal_resonance"}, Effects: map[string]float64{"building_hp": 0.2}},
	{ID: "teleportation", Name: "传送术", Description: "单位移动+1格", Cost: 60, Category: "magic", Prerequisites: []string{"arcane_defense"}, Effects: map[string]float64{"move_speed": 1.0}},

	{ID: "fossil_extraction", Name: "化石提取", Description: "化石燃料产出+30%", Cost: 20, Category: "energy", Effects: map[string]float64{"fossil_production": 0.3}},
	{ID: "steam_power", Name: "蒸汽动力", Description: "电梯消耗-30%", Cost: 30, Category: "energy", Prerequisites: []string{"fossil_extraction"}, Effects: map[string]float64{"elevator_efficiency": 0.3}},
	{ID: "geothermal", Name: "地热能源", Description: "熔岩环境伤害-50%", Cost: 45, Category: "energy", Prerequisites: []string{"steam_power"}, Effects: map[string]float64{"lava_resistance": 0.5}},

	{ID: "diplomacy", Name: "外交学", Description: "贸易效率+20%", Cost: 15, Category: "civil", Effects: map[string]float64{"trade_efficiency": 0.2}},
	{ID: "logistics", Name: "后勤学", Description: "资源传输容量+50%", Cost: 25, Category: "civil", Effects: map[string]float64{"pipe_capacity": 0.5}},
	{ID: "urban_planning", Name: "城市规划", Description: "人口上限+20%", Cost: 35, Category: "civil", Prerequisites: []string{"logistics"}, Effects: map[string]float64{"population_cap": 0.2}},

	{ID: "scouting", Name: "侦察术", Description: "视野范围+1格", Cost: 12, Category: "exploration", Effects: map[string]float64{"vision_range": 1.0}},
	{ID: "cartography", Name: "制图学", Description: "探索速度+20%", Cost: 20, Category: "exploration", Prerequisites: []string{"scouting"}, Effects: map[string]float64{"explore_speed": 0.2}},
	{ID: "ancient_lore", Name: "远古知识", Description: "遗迹奖励+50%", Cost: 40, Category: "exploration", Prerequisites: []string{"cartography"}, Effects: map[string]float64{"relic_bonus": 0.5}},

	{ID: "storage_tech", Name: "存储技术", Description: "仓储容量+50%", Cost: 18, Category: "economy", Effects: map[string]float64{"storage_cap": 0.5}},
	{ID: "trade_routes", Name: "贸易路线", Description: "贸易距离不限", Cost: 28, Category: "economy", Prerequisites: []string{"storage_tech"}, Effects: map[string]float64{"trade_range": 1.0}},
	{ID: "banking", Name: "地下银行", Description: "每回合利息5%", Cost: 50, Category: "economy", Prerequisites: []string{"trade_routes"}, Effects: map[string]float64{"interest": 0.05}},

	{ID: "basic_research", Name: "基础研究", Description: "研究点产出+1", Cost: 8, Category: "research", Effects: map[string]float64{"research_point": 1.0}},
	{ID: "advanced_research", Name: "高级研究", Description: "研究点产出+2", Cost: 25, Category: "research", Prerequisites: []string{"basic_research"}, Effects: map[string]float64{"research_point": 2.0}},
	{ID: "scientific_method", Name: "科学方法", Description: "研究速度+30%", Cost: 40, Category: "research", Prerequisites: []string{"advanced_research"}, Effects: map[string]float64{"research_speed": 0.3}},

	{ID: "fortification", Name: "要塞化", Description: "城墙防御+50%", Cost: 30, Category: "defense", Effects: map[string]float64{"wall_defense": 0.5}},
	{ID: "trap_design", Name: "陷阱设计", Description: "入侵方伤害-15%", Cost: 25, Category: "defense", Effects: map[string]float64{"invader_penalty": 0.15}},
	{ID: "bunker", Name: "地堡系统", Description: "建筑被攻击时有30%概率免伤", Cost: 45, Category: "defense", Prerequisites: []string{"fortification", "trap_design"}, Effects: map[string]float64{"dodge_chance": 0.3}},
}
