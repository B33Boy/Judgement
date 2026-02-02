package game

type PlayerHandChangePayload struct {
	Cards Hand `json:"cards"`
}

type MakeBid struct {
	Bid Bid `json:"bid"`
}

type InvalidActionPayload struct {
	Message string `json:"message"`
}
