package game

type PlayerHandChangePayload struct {
	Cards Hand `json:"cards"`
}

type RoundInfoPayload struct {
	Round      Round  `json:"round"`
	TurnPlayer string `json:"turnPlayer"`
	State      State  `json:"state"`
}

type MakeBid struct {
	Bid Bid `json:"bid"`
}

type InvalidActionPayload struct {
	Message string `json:"message"`
}
