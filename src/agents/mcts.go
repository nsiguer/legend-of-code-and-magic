package agent

import (
	g "game"

	mcts "github.com/nsiguer/go-mcts"
	
)

const (
	BIAS_PARAMETER 	= 0.9
	MCTS_ITERATION	= 50
	MCTS_SIMULATION	= 50
	MCTS_TIMEOUT	= 10 
)
type AgentMCTS struct {

}

func NewAgentMCTS() *AgentMCTS {
	return &AgentMCTS{}
}

func (a *AgentMCTS) Name() string { return "MCTS" }
func (a *AgentMCTS) Think(s *g.State) []*g.Move {
	moves := mcts.MonteCarloTimeout(s, BIAS_PARAMETER, MCTS_ITERATION, MCTS_SIMULATION, MCTS_TIMEOUT)
	return moves
}