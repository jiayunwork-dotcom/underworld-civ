package game

type PlayerState struct {
	PlayerID       string              `json:"player_id"`
	Username       string              `json:"username"`
	Race           Race                `json:"race"`
	Color          string              `json:"color"`
	IsHost         bool                `json:"is_host"`
	Resources      ResourceStorage     `json:"resources"`
	Production     ResourceProduction  `json:"production"`
	Population     int                 `json:"population"`
	PopulationCap  int                 `json:"population_cap"`
	ResearchPoints int                 `json:"research_points"`
	Techs          map[string]bool     `json:"techs"`
	ResearchQueue  []string            `json:"research_queue"`
	CurrentResearch string             `json:"current_research"`
	ResearchProgress int               `json:"research_progress"`
	Units          []Unit              `json:"units"`
	Buildings      map[string]*Building `json:"buildings"`
	Score          int                 `json:"score"`
	Eliminated     bool                `json:"eliminated"`
	VisionRange    int                 `json:"vision_range"`
	Diplomacy      map[string]*DiplomaticRelation `json:"diplomacy"`
	TradeOffers    []TradeOffer        `json:"trade_offers"`
	Embargoes      []Embargo           `json:"embargoes"`
	Ready          bool                `json:"ready"`
}
