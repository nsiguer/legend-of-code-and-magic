package game

import (
	//"fmt"
	ag "agents"
)

func main() {

	ai := ag.NewAI()
	ai.LoadAgentMCTS() 

	fmt.Println(ai)

}