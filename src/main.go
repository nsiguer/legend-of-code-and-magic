package main

import "fmt"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 *

	*/
func main() {
	game := NewGame()
	p1 := NewPlayer(1)
	p2 := NewPlayer(2)

	e := game.AddPlayer(p1)
	if e != nil {
		fmt.Println(e)
	}
	e = game.AddPlayer(p2)
	if e != nil {
		fmt.Println(e)
	}

	for ;; {
		winner, _ := game.Start()
		if winner != nil {
			winner.PrintDeck()
		}
		game.Clear()
		fmt.Println(game)
	}
}
/*
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
        }

        // fmt.Fprintln(os.Stderr, "Debug messages...")
        fmt.Println("PASS")// Write action to stdout
    }
}
*/
