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
	    fmt.Fprintln(os.Stderr, "Receive", playerHealth, playerMana, playerDeck, playerRune)
        }



        var opponentHand int
        fmt.Scan(&opponentHand)

	fmt.Fprintln(os.Stderr, "Hand:", opponentHand)

        var cardCount int
        fmt.Scan(&cardCount)

	fmt.Fprintln(os.Stderr, "CardsCount:", cardCount)




        for i := 0; i < cardCount; i++ {
            var cardNumber, instanceId, location, cardType, cost, attack, defense int
            var abilities string
            var myHealthChange, opponentHealthChange, cardDraw int
            fmt.Scan(&cardNumber, &instanceId, &location, &cardType, &cost, &attack, &defense, &abilities, &myHealthChange, &opponentHealthChange, &cardDraw)
         
            fmt.Fprintln(os.Stderr, "Card:", cardNumber, instanceId, location, cardType, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw)
        }
            fmt.Fprintln(os.Stderr, "Done reading data")


	fmt.Println("PASS")
        
    }
}


