package game

type MakeBid struct {
	Bid Bid `json:"bid"`
}

type InvalidActionPayload struct {
	Message string `json:"message"`
}
