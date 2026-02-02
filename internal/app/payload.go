package app

import t "github.com/B33Boy/Judgement/internal/types"

type PlayerPublic struct {
	ID   t.PlayerID `json:"id"`
	Name string     `json:"name"`
}
