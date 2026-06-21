package game

type SpyStatus string

const (
	SpyStatusIdle       SpyStatus = "idle"
	SpyStatusLurking    SpyStatus = "lurking"
	SpyStatusSuccessful SpyStatus = "successful"
	SpyStatusCaught     SpyStatus = "caught"
)

type TechBlockade struct {
	TargetPlayerID string       `json:"target_player_id"`
	Category       TechCategory `json:"category"`
	TurnsRemaining int          `json:"turns_remaining"`
	IssuerPlayerID string       `json:"issuer_player_id"`
}

type PlayerState struct {
	PlayerID        string              `json:"player_id"`
	Username        string              `json:"username"`
	Race            Race                `json:"race"`
	Color           string              `json:"color"`
	IsHost          bool                `json:"is_host"`
	Resources       ResourceStorage     `json:"resources"`
	Production      ResourceProduction  `json:"production"`
	Population      int                 `json:"population"`
	PopulationCap   int                 `json:"population_cap"`
	ResearchPoints  int                 `json:"research_points"`
	KnowledgeReserve int                `json:"knowledge_reserve"`
	Techs           map[string]bool     `json:"techs"`
	TechProgresses  map[string]int      `json:"tech_progresses"`
	ResearchQueue   []string            `json:"research_queue"`
	CurrentResearch string              `json:"current_research"`
	ResearchProgress int                `json:"research_progress"`
	TechSynergies   map[TechCategory]*TechSynergy `json:"tech_synergies"`
	Units           []Unit              `json:"units"`
	Buildings       map[string]*Building `json:"buildings"`
	Score           int                 `json:"score"`
	Eliminated      bool                `json:"eliminated"`
	VisionRange     int                 `json:"vision_range"`
	Diplomacy       map[string]*DiplomaticRelation `json:"diplomacy"`
	TradeOffers     []TradeOffer        `json:"trade_offers"`
	Embargoes       []Embargo           `json:"embargoes"`
	IncomingBlockades []TechBlockade    `json:"incoming_blockades"`
	OutgoingBlockades []TechBlockade    `json:"outgoing_blockades"`
	Ready           bool                `json:"ready"`
}
