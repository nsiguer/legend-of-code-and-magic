package main

import (
	"fmt"
	"os"
)
type IADummy struct {}
func NewIADummy() *IADummy { return &IADummy{} }
func (ia *IADummy) Move(players, cards [][]interface{}) string {
	return "PASS"
}


func main() {
 
    
    for {

        for i := 0; i < 2; i++ {
            var playerHealth, playerMana, playerDeck, playerRune int
            fmt.Scan(&playerHealth, &playerMana, &playerDeck, &playerRune)
	   }



        var opponentHand int
        fmt.Scan(&opponentHand)

        var cardCount int
        fmt.Scan(&cardCount)

        for i := 0; i < cardCount; i++ {
            var cardNumber, instanceId, location, cardType, cost, attack, defense int
            var abilities string
            var myHealthChange, opponentHealthChange, cardDraw int
            fmt.Scan(&cardNumber, &instanceId, &location, &cardType, &cost, &attack, &defense, &abilities, &myHealthChange, &opponentHealthChange, &cardDraw)
         
            fmt.Fprintln(os.Stderr, "Card:", cardNumber, instanceId, location, cardType, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw)
        }

	    fmt.Println("PASS")
        
    }
}


