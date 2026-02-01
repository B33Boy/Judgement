package game

import t "github.com/B33Boy/Judgement/internal/types"

type PlayerScore map[t.PlayerID][]Score

type ScoreBoard struct {
	Trumps       []string
	PlayerScores map[t.PlayerID][]Score `json:"playerscores"`
}

func NewScoreboard(playerCnt int, gamePlayers PlayerMap, maxRounds Round) PlayerScore {
	scores := make(PlayerScore, playerCnt)
	for playerId := range gamePlayers {
		scores[playerId] = make([]Score, maxRounds)
	}
	return scores
}
