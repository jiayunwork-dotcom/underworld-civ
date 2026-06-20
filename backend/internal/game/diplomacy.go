package game

type DiplomaticStatus string

const (
	DiplomacyNeutral  DiplomaticStatus = "neutral"
	DiplomacyAlliance DiplomaticStatus = "alliance"
	DiplomacyWar      DiplomaticStatus = "war"
	DiplomacyTrade    DiplomaticStatus = "trade"
)

type DiplomaticRelation struct {
	PlayerA string           `json:"player_a"`
	PlayerB string           `json:"player_b"`
	Status  DiplomaticStatus `json:"status"`
	TradeAccess bool         `json:"trade_access"`
	PassageRights bool       `json:"passage_rights"`
	SharedVision bool        `json:"shared_vision"`
	BorderRecognized bool    `json:"border_recognized"`
	CooldownTurns int        `json:"cooldown_turns"`
	TreatyViolationPenalty int `json:"treaty_violation_penalty"`
}

type TradeOffer struct {
	ID         string    `json:"id"`
	FromPlayer string    `json:"from_player"`
	ToPlayer   string    `json:"to_player"`
	Offer      Resources `json:"offer"`
	Demand     Resources `json:"demand"`
	TurnsValid int       `json:"turns_valid"`
}

type Embargo struct {
	TargetPlayer string   `json:"target_player"`
	ResourceType string   `json:"resource_type"`
	ProposedBy   string   `json:"proposed_by"`
	VotesYes     []string `json:"votes_yes"`
	VotesNo      []string `json:"votes_no"`
	Passed       bool     `json:"passed"`
	TurnsLeft    int      `json:"turns_left"`
}
