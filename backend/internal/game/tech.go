package game

type TechCategory string

const (
	TechMilitary TechCategory = "military"
	TechEconomy  TechCategory = "economy"
	TechMining   TechCategory = "mining"
	TechSpecial  TechCategory = "special"
)

type Tech struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Cost          int                `json:"cost"`
	Category      TechCategory       `json:"category"`
	Tier          int                `json:"tier"`
	Prerequisites []string           `json:"prerequisites"`
	Effects       map[string]float64 `json:"effects"`
	RaceSpecific  Race               `json:"race_specific,omitempty"`
}

var TechDefs = []Tech{
	{ID: "mil_1_sharp_blades", Name: "锋利刀刃", Description: "所有单位攻击力+5", Cost: 5, Category: TechMilitary, Tier: 1, Prerequisites: []string{}, Effects: map[string]float64{"attack_flat": 5}},
	{ID: "mil_2_leather_armor", Name: "皮甲工艺", Description: "所有单位防御力+3", Cost: 8, Category: TechMilitary, Tier: 2, Prerequisites: []string{"mil_1_sharp_blades"}, Effects: map[string]float64{"defense_flat": 3}},
	{ID: "mil_3_reinforced_armor", Name: "强化护甲", Description: "步兵单位防御力额外+8", Cost: 12, Category: TechMilitary, Tier: 3, Prerequisites: []string{"mil_2_leather_armor"}, Effects: map[string]float64{"infantry_defense": 8}},
	{ID: "mil_4_combat_training", Name: "战斗训练", Description: "所有单位攻击力+10%", Cost: 18, Category: TechMilitary, Tier: 4, Prerequisites: []string{"mil_3_reinforced_armor"}, Effects: map[string]float64{"attack_pct": 0.10}},
	{ID: "mil_5_tactical_formation", Name: "战术阵型", Description: "所有单位防御力+15%", Cost: 25, Category: TechMilitary, Tier: 5, Prerequisites: []string{"mil_4_combat_training"}, Effects: map[string]float64{"defense_pct": 0.15}},
	{ID: "mil_6_steel_weapons", Name: "精钢武器", Description: "所有单位攻击力+10", Cost: 35, Category: TechMilitary, Tier: 6, Prerequisites: []string{"mil_5_tactical_formation"}, Effects: map[string]float64{"attack_flat": 10}},
	{ID: "mil_7_morale_training", Name: "士气训练", Description: "部队不会因50%伤亡以下而溃退", Cost: 45, Category: TechMilitary, Tier: 7, Prerequisites: []string{"mil_6_steel_weapons"}, Effects: map[string]float64{"morale_threshold": 0.1}},
	{ID: "mil_8_siege_mastery", Name: "攻城专精", Description: "攻城单位伤害+50%，建筑伤害+20%", Cost: 55, Category: TechMilitary, Tier: 8, Prerequisites: []string{"mil_7_morale_training"}, Effects: map[string]float64{"siege_damage": 0.50, "building_damage": 0.20}},
	{ID: "mil_9_elite_training", Name: "精锐训练", Description: "所有单位生命值+25%", Cost: 70, Category: TechMilitary, Tier: 9, Prerequisites: []string{"mil_8_siege_mastery"}, Effects: map[string]float64{"hp_pct": 0.25}},
	{ID: "mil_10_total_war", Name: "全面战争", Description: "所有单位攻击+20%，防御+15%，生命+15%", Cost: 80, Category: TechMilitary, Tier: 10, Prerequisites: []string{"mil_9_elite_training"}, Effects: map[string]float64{"attack_pct": 0.20, "defense_pct": 0.15, "hp_pct": 0.15}},

	{ID: "eco_1_basic_tools", Name: "基础工具", Description: "石材产出+20%", Cost: 5, Category: TechEconomy, Tier: 1, Prerequisites: []string{}, Effects: map[string]float64{"stone_production": 0.20}},
	{ID: "eco_2_smelting", Name: "冶炼技术", Description: "金属产出+25%", Cost: 8, Category: TechEconomy, Tier: 2, Prerequisites: []string{"eco_1_basic_tools"}, Effects: map[string]float64{"metal_production": 0.25}},
	{ID: "eco_3_agriculture", Name: "地下农业", Description: "发光菌产出+30%", Cost: 12, Category: TechEconomy, Tier: 3, Prerequisites: []string{"eco_2_smelting"}, Effects: map[string]float64{"fungus_production": 0.30}},
	{ID: "eco_4_water_management", Name: "水源管理", Description: "水资源产出+25%，消耗-15%", Cost: 18, Category: TechEconomy, Tier: 4, Prerequisites: []string{"eco_3_agriculture"}, Effects: map[string]float64{"water_production": 0.25, "water_efficiency": 0.15}},
	{ID: "eco_5_crystal_refining", Name: "魔晶精炼", Description: "魔晶产出+40%", Cost: 25, Category: TechEconomy, Tier: 5, Prerequisites: []string{"eco_4_water_management"}, Effects: map[string]float64{"crystal_production": 0.40}},
	{ID: "eco_6_fuel_extraction", Name: "燃料萃取", Description: "化石燃料产出+35%", Cost: 35, Category: TechEconomy, Tier: 6, Prerequisites: []string{"eco_5_crystal_refining"}, Effects: map[string]float64{"fossil_production": 0.35}},
	{ID: "eco_7_automation", Name: "自动化生产", Description: "所有资源产出+15%", Cost: 45, Category: TechEconomy, Tier: 7, Prerequisites: []string{"eco_6_fuel_extraction"}, Effects: map[string]float64{"all_production": 0.15}},
	{ID: "eco_8_storage_expansion", Name: "仓储扩容", Description: "仓储容量+100%", Cost: 55, Category: TechEconomy, Tier: 8, Prerequisites: []string{"eco_7_automation"}, Effects: map[string]float64{"storage_cap": 1.0}},
	{ID: "eco_9_urban_planning", Name: "城市规划", Description: "人口上限+50%，建筑建造速度+25%", Cost: 70, Category: TechEconomy, Tier: 9, Prerequisites: []string{"eco_8_storage_expansion"}, Effects: map[string]float64{"population_cap": 0.50, "build_speed": 0.25}},
	{ID: "eco_10_economic_dominance", Name: "经济霸权", Description: "所有资源产出+30%，仓储+50%", Cost: 80, Category: TechEconomy, Tier: 10, Prerequisites: []string{"eco_9_urban_planning"}, Effects: map[string]float64{"all_production": 0.30, "storage_cap": 0.50}},

	{ID: "min_1_pickaxe_upgrade", Name: "镐头升级", Description: "挖掘速度+25%", Cost: 5, Category: TechMining, Tier: 1, Prerequisites: []string{}, Effects: map[string]float64{"mining_speed": 0.25}},
	{ID: "min_2_support_timbers", Name: "支撑木架", Description: "塌方概率-20%", Cost: 8, Category: TechMining, Tier: 2, Prerequisites: []string{"min_1_pickaxe_upgrade"}, Effects: map[string]float64{"collapse_chance": 0.20}},
	{ID: "min_3_drilling_tech", Name: "钻探技术", Description: "挖掘速度+25%，硬岩挖掘效率+30%", Cost: 12, Category: TechMining, Tier: 3, Prerequisites: []string{"min_2_support_timbers"}, Effects: map[string]float64{"mining_speed": 0.25, "hard_rock_bonus": 0.30}},
	{ID: "min_4_steel_beams", Name: "钢梁支护", Description: "塌方概率-40%", Cost: 18, Category: TechMining, Tier: 4, Prerequisites: []string{"min_3_drilling_tech"}, Effects: map[string]float64{"collapse_chance": 0.40}},
	{ID: "min_5_surveying", Name: "地质勘探", Description: "挖掘时有30%概率额外获得矿物", Cost: 25, Category: TechMining, Tier: 5, Prerequisites: []string{"min_4_steel_beams"}, Effects: map[string]float64{"bonus_mineral_chance": 0.30}},
	{ID: "min_6_ventilation", Name: "通风系统", Description: "深层挖掘不减速，塌方概率再-20%", Cost: 35, Category: TechMining, Tier: 6, Prerequisites: []string{"min_5_surveying"}, Effects: map[string]float64{"collapse_chance": 0.20, "deep_layer_bonus": 1.0}},
	{ID: "min_7_heavy_machinery", Name: "重型机械", Description: "挖掘速度+50%", Cost: 45, Category: TechMining, Tier: 7, Prerequisites: []string{"min_6_ventilation"}, Effects: map[string]float64{"mining_speed": 0.50}},
	{ID: "min_8_earthquake_prediction", Name: "地震预警", Description: "塌方时50%概率免于损失，事件伤害-50%", Cost: 55, Category: TechMining, Tier: 8, Prerequisites: []string{"min_7_heavy_machinery"}, Effects: map[string]float64{"collapse_save": 0.50, "event_damage_reduction": 0.50}},
	{ID: "min_9_layer_transport", Name: "电梯提速", Description: "跨层移动不消耗额外回合，电梯建造速度+50%", Cost: 70, Category: TechMining, Tier: 9, Prerequisites: []string{"min_8_earthquake_prediction"}, Effects: map[string]float64{"elevator_speed": 0.50, "cross_layer_move": 1.0}},
	{ID: "min_10_mine_emperor", Name: "挖矿帝国", Description: "挖掘速度+100%，塌方率-80%，额外矿物概率+50%", Cost: 80, Category: TechMining, Tier: 10, Prerequisites: []string{"min_9_layer_transport"}, Effects: map[string]float64{"mining_speed": 1.0, "collapse_chance": 0.80, "bonus_mineral_chance": 0.50}},

	{ID: "spc_1_racial_heritage", Name: "种族传承", Description: "激活种族基础天赋加成+20%", Cost: 5, Category: TechSpecial, Tier: 1, Prerequisites: []string{}, Effects: map[string]float64{"racial_bonus_boost": 0.20}},
	{ID: "spc_2_sacred_knowledge", Name: "神圣知识", Description: "研究点数产出+2/回合", Cost: 8, Category: TechSpecial, Tier: 2, Prerequisites: []string{"spc_1_racial_heritage"}, Effects: map[string]float64{"research_point": 2.0}},
	{ID: "spc_3_mystic_vision", Name: "神秘视野", Description: "视野范围+1格", Cost: 12, Category: TechSpecial, Tier: 3, Prerequisites: []string{"spc_2_sacred_knowledge"}, Effects: map[string]float64{"vision_range": 1.0}},
	{ID: "spc_4_race_tech_1", Name: "种族专精I", Description: "解锁种族专属能力第一阶", Cost: 18, Category: TechSpecial, Tier: 4, Prerequisites: []string{"spc_3_mystic_vision"}, Effects: map[string]float64{"race_tier_1": 1.0}},
	{ID: "spc_5_arcane_research", Name: "奥术研究", Description: "研究速度+30%，魔晶产出+25%", Cost: 25, Category: TechSpecial, Tier: 5, Prerequisites: []string{"spc_4_race_tech_1"}, Effects: map[string]float64{"research_speed": 0.30, "crystal_production": 0.25}},
	{ID: "spc_6_building_durability", Name: "建筑耐久", Description: "所有建筑生命值+30%", Cost: 35, Category: TechSpecial, Tier: 6, Prerequisites: []string{"spc_5_arcane_research"}, Effects: map[string]float64{"building_hp": 0.30}},
	{ID: "spc_7_race_tech_2", Name: "种族专精II", Description: "解锁种族专属能力第二阶", Cost: 45, Category: TechSpecial, Tier: 7, Prerequisites: []string{"spc_6_building_durability"}, Effects: map[string]float64{"race_tier_2": 1.0}},
	{ID: "spc_8_teleport_theory", Name: "传送理论", Description: "单位移动力+1格", Cost: 55, Category: TechSpecial, Tier: 8, Prerequisites: []string{"spc_7_race_tech_2"}, Effects: map[string]float64{"move_speed": 1.0}},
	{ID: "spc_9_race_tech_3", Name: "种族专精III", Description: "解锁种族终极能力", Cost: 70, Category: TechSpecial, Tier: 9, Prerequisites: []string{"spc_8_teleport_theory"}, Effects: map[string]float64{"race_tier_3": 1.0}},
	{ID: "spc_10_ascension", Name: "文明升华", Description: "全属性+10%，解锁科技胜利条件", Cost: 80, Category: TechSpecial, Tier: 10, Prerequisites: []string{"spc_9_race_tech_3"}, Effects: map[string]float64{"all_stats": 0.10, "tech_victory": 1.0}},
}

type TechState struct {
	Researched bool `json:"researched"`
	Progress   int  `json:"progress"`
}

func GetTechByID(id string) *Tech {
	for i := range TechDefs {
		if TechDefs[i].ID == id {
			return &TechDefs[i]
		}
	}
	return nil
}

func GetTechsByCategory(cat TechCategory) []Tech {
	result := make([]Tech, 0)
	for _, t := range TechDefs {
		if t.Category == cat {
			result = append(result, t)
		}
	}
	return result
}

func CheckPrerequisites(techID string, researched map[string]bool) bool {
	tech := GetTechByID(techID)
	if tech == nil {
		return false
	}
	for _, prereq := range tech.Prerequisites {
		if !researched[prereq] {
			return false
		}
	}
	return true
}

func (p *PlayerState) GetTechEffect(key string) float64 {
	var total float64 = 0
	for techID := range p.Techs {
		tech := GetTechByID(techID)
		if tech != nil {
			if val, ok := tech.Effects[key]; ok {
				total += val
			}
		}
	}
	return total
}
