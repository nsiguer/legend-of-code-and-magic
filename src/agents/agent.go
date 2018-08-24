package agent

import (
	g "game"
)

/* STARTING */
type Agent interface {
	Name()			string
	Think(s *g.State) []*g.Move
}

type AI struct {
	Agents	[]Agent
}

func NewAI() *AI {
	return &AI{
		Agents: make([]Agent, 0),
	}
}
func (ai *AI) LoadAgentRandom() {
	ai.Agents = append(ai.Agents, NewAgentRandom())
}
func (ai *AI) LoadAgentMCTS() {
	ai.Agents = append(ai.Agents, NewAgentMCTS())
}

func (ai *AI) LoadAgentDraft() {
	ai.Agents = append(ai.Agents, NewAgentDraft())
}

func (ai *AI) GetAgent(name_agent string) Agent {
	for _, a := range(ai.Agents) {
		if a.Name() == name_agent { return a }
	}
	return nil
}
func (ai *AI) Think(name_agent string, s *g.State) []*g.Move {
	a := ai.GetAgent(name_agent)
	if a != nil {
		return a.Think(s)
	}
	return []*g.Move{}
}
