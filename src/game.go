package main

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"os/exec"
	"os"
	"io"
	"bufio"



	//copier "github.com/jinzhu/copier"
)

const (

	MAX_MANA		= 12
	MAX_PLAYERS		= 2
	MIN_PLAYERS		= 2
	STARTING_MANA		= 1
	STARTING_LIFE		= 30
	STARTING_CARD		= 30
	DECK_CARDS		= 30
	DRAFT_PICK		= 3

        CARD_TYPE_MONSTER           = 0

        CARD_ABILITY_BREAKTHROUGH   = 0
        CARD_ABILITY_GUARD          = 3
        CARD_ABILITY_CHARGE         = 1
)

var CARDS = []*Card{
	NewCard(1, -1, 0, CARD_TYPE_MONSTER, 1, 2, 2, "------", 0, 0, 0),
	NewCard(2, -1, 0, CARD_TYPE_MONSTER, 2, 4, 1, "------", 0, 0, 0),
	NewCard(3, -1, 0, CARD_TYPE_MONSTER, 2, 1, 5, "------", 0, 0, 0),
	NewCard(4, -1, 0, CARD_TYPE_MONSTER, 2, 2, 3, "------", 0, 0, 0),
	NewCard(5, -1, 0, CARD_TYPE_MONSTER, 4, 4, 5, "------", 0, 0, 0),
	NewCard(6, -1, 0, CARD_TYPE_MONSTER, 4, 1, 8, "------", 0, 0, 0),
	NewCard(7, -1, 0, CARD_TYPE_MONSTER, 5, 8, 2, "------", 0, 0, 0),
	NewCard(8, -1, 0, CARD_TYPE_MONSTER, 5, 6, 5, "------", 0, 0, 0),
	NewCard(9, -1, 0, CARD_TYPE_MONSTER, 7, 8, 8, "------", 0, 0, 0),
	NewCard(10, -1, 0, CARD_TYPE_MONSTER, 9, 10, 10, "------", 0, 0, 0),
}


func PickRandomCard() (*Card) {
	source	:= rand.NewSource(time.Now().UnixNano())
	random	:= rand.New(source)
	idx	:= random.Intn(len(CARDS))
	return CARDS[idx]
}

func in_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}

/************ MONTE CARLO TREE ***********/



/*****************************************/

type Deck struct {
	Cards []*Card
}

type Card struct {
	CardNumber 		int
        Id    			int
	Location		int
        Type  			int
        Cost  			int

	Attack			int
	Defense			int
	Abilities		[]string
	HealthChange		int
	OpponentHealthChange 	int
	CardDraw		int
	

}

type IAPlayer struct {
	id	int
	binary	string
	cmd	*exec.Cmd
	stdin	io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
}


func NewIAPlayer(id int, binary string) *IAPlayer{
	cmd := exec.Command(binary)
	if cmd == nil {
		fmt.Println("Unable to load", binary)
		return nil
	}

	stdin, err1 := cmd.StdinPipe()
	if err1 != nil {
		fmt.Println(err1)
		return nil
	}

	stdout, err2 := cmd.StdoutPipe()
	if err2 != nil {
		fmt.Println(err2)
		return nil
	}

	stderr, err3 := cmd.StderrPipe()
	if err3 != nil {
		fmt.Println(err3)
		return nil
	}
	ia := &IAPlayer {
		id: id,
		binary: binary,
		cmd: cmd,
		stdin: stdin, 
		stdout: stdout,
		stderr: stderr,  
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return ia
}

func (ia *IAPlayer) ReadMove(ms int) (string, error) {
	reader := bufio.NewReader(ia.stdout)
	scanner := bufio.NewScanner(reader)
	done_reading := make(chan bool)
	go func() {
		scanner.Scan()
		done_reading <- true
	} ()

	readerErr := bufio.NewReader(ia.stderr)
	scannerErr := bufio.NewScanner(readerErr)

	go func() {
		for scannerErr.Scan() {
		  fmt.Println("STDERR:", scannerErr.Text())
		}
	} ()


	select {
	case <- time.After(time.Nanosecond * time.Duration(ms) * (1000000)):
		fmt.Println("Timeout reading")
		return "", fmt.Errorf("Timeout after %d ms", ms)
	case <- done_reading:
		break
	}
	return scanner.Text(), nil
}

func (ia *IAPlayer) WriteData(data string) (error) {
	fmt.Println("STDIN:", data)
	ia.stdin.Write([]byte(data))
	ia.stdin.Write([]byte("\n"))

	return nil
}
func (ia *IAPlayer) Move(players, cards [][]interface{}, opponentHand int, timeout int) (string, error) {
	var str string

	for i := 0; i < len(players) ; i++ {
		str = fmt.Sprintf("%v", players[i])	
		ia.WriteData(str[1:len(str)-1])
	}

	ia.WriteData(fmt.Sprintf("%d", opponentHand))
	ia.WriteData(fmt.Sprintf("%d", len(cards)))

	
	for i := 0; i < len(cards) ; i++ {
		str = fmt.Sprintf("%v", cards[i])	
		ia.WriteData(str[1:len(str)-1])
	}

	return ia.ReadMove(timeout)
}



func NewCard(cardNumber int,
	     id int,
	     location int,
             type_ int,
             cost int,
             attack int,
             defense int,
             abilities string,
	     heroHealthChange int,
	     opponentHealthChange int,
	     cardDraw int,
		
) *Card {
        return &Card{
		CardNumber: 	cardNumber,
                Id:             id,
		Location:	location,
                Type:           type_,
                Cost:           cost,
                Attack:         attack,
                Defense:        defense,
                Abilities:      strings.Split(abilities, ""),
		HealthChange: 	heroHealthChange,
		OpponentHealthChange: opponentHealthChange,
		CardDraw:	cardDraw,
        }
}

func (c *Card) Raw() []interface{} {
	return []interface{} {
			c.CardNumber,
			c.Id,
			c.Location,
			c.Type,
			c.Cost,
			c.Attack,
			c.Defense,
			strings.Join(c.Abilities, ""),
			c.HealthChange,
			c.OpponentHealthChange,
			c.CardDraw,
		}
}


type Player struct {
	Id    int
	Deck  *Deck
	Life  int
	Mana  int
	Board []*Card
	Hand  []*Card

	IA	*IAPlayer
}

type Game struct {
	players        []*Player
	Turn		int
}

func NewGame() *Game {
	return &Game{
		players: make([]*Player, 0),
		Turn: 1,
	}
}


func NewDeck() *Deck {
	return &Deck{
		Cards: make([]*Card, 0),
	}
}

func (d *Deck) Count() int {
	return len(d.Cards)
}
func (d *Deck) AddCard(c *Card) {
	d.Cards = append(d.Cards, c)
}

func (d *Deck) Draw() (*Card, error) {
	if len(d.Cards) > 0 {
		c := d.Cards[0]
		d.Cards = d.Cards[1:]
		return c, nil
	} else {
		return nil, errors.New("There is no more card in the deck")
	}
}

func (d *Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

func NewPlayer(id, life, mana int, ia *IAPlayer) *Player {
        return &Player{
                Id:     id,
                Mana:   mana,
                Life:   life,
		Deck: 	NewDeck(),
                Board:  make([]*Card, 0),
                Hand:   make([]*Card, 0),
		IA: ia,
        }

}


func (p *Player) Raw() []interface{} {
	return []interface{}{
		p.Life,
		p.Mana,
		p.Deck.Count(),
		0,
	}
}
func (p *Player) DrawCard() {
    c, err := p.Deck.Draw()
    if err == nil {
	p.Draw(c)
    }
}

func (p *Player) DrawCardN(n int) {
    for i := 0 ; i < n ; i++ {
    	c, err := p.Deck.Draw()
    	if err == nil {
		p.Draw(c)
    	} else {
		break
	}
    }
}

func (p *Player) Draw(c *Card) {
    if p.HandGetCard(c.Id) == nil {
        p.Hand = append(p.Hand, c)
    }
}

func (p *Player) SetMana(mana int) {
        if mana <= 12 && mana >= 0 {
                p.Mana = mana
        }
}

func (p *Player) SetLife(life int) {
        p.Life = life
}

func (p *Player) HandGetCard(id int) *Card {
        for _, c := range(p.Hand) {
                if c.Id == id {
                        return c
                }
        }
        return nil
}

func (p *Player) PlayerToInt() []int {	
	return []int{
		p.Life,
		p.Mana,
		p.Deck.Count(),
		0,
	}
}

func (p *Player) CardsToInt(location int) [][]interface{} {	
	cards := make([][]interface{}, 0)

	for _, c := range(append(p.Hand, p.Board...)) {
		cards = append(cards, c.Raw())
	}

	return cards
}

func (p *Player) HandRemoveCard(id int) (error) {
        idx := -1
        for i, c := range(p.Hand) {
                if c.Id == id {
                        idx = i
                        break
                }
        }
        if idx == -1 {
                return fmt.Errorf("RemoveCard: Hand with id %i not in hand", idx)
        }
        
        l := len(p.Hand)
        p.Hand[idx] = p.Hand[l - 1]
        if l >= 2 {
            p.Hand = p.Hand[:l - 1]
        } else if l == 1 {
            p.Hand = make([]*Card, 0)
        }
        return nil
}
func (p *Player) HandPlayCard(id int) (error) {
        c := p.HandGetCard(id)
        if c == nil {
                return fmt.Errorf("Card with id %d is not present in hand of player %d", id, p.Id)
        }
        switch c.Type {
        case CARD_TYPE_MONSTER:
                p.Board = append(p.Board, c)
                p.HandRemoveCard(id)
                  
        default:
                return fmt.Errorf("Unkow card type %d for player %d", c.Type, p.Id)
        }

        return nil
}

func (p *Player) ReduceLife(damage int) {
        p.Life = p.Life - damage
}

func (p* Player) BoardAttackCard(id, damage int) bool {
        for _, c := range(p.Board) {
                if c.Id == id {
                    switch c.Type {
                    case CARD_TYPE_MONSTER:
                        c.Defense -= damage
                        if c.Defense <= 0 {
                            p.BoardRemoveCard(c.Id)
                            return true
                        }
                        
                    }    
                    return false
                }
        }
        
        return false
}
func (p *Player) UpdateBoard(card *Card) bool {
        
        for _, c := range(p.Board) {
                if c.Id == card.Id {
                    switch c.Type {
                    case CARD_TYPE_MONSTER:
                        //c.Celerity = 1
                        fmt.Fprintln(os.Stderr, "Reduce defense of", c.Id, c.Defense, " -> ", card.Defense)
                        c.Defense = card.Defense
                       
                        if c.Defense <= 0 {
                            p.BoardRemoveCard(c.Id)
                        }
                    }    
                    return true
                }
        }
        fmt.Fprintln(os.Stderr, p.Id, "Board add monster", card)
        p.Board = append(p.Board, card)
        return false
}
func (p *Player) BoardGetCard(id int) *Card {
        for _, c := range(p.Board) {
                if c.Id == id {
                        return c
                }
        }
        return nil
}

func (p *Player) BoardGetGuardsId() []int {
	ids := make([]int, 0)
        for _, c := range(p.Board) {
                if c.Abilities[CARD_ABILITY_GUARD] != "-" {
                        ids = append(ids, c.Id)
                }
        }
        return ids
}


func (p *Player) BoardRemoveCard(id int) error {
        idx := -1
        for i, b := range p.Board {
                if b.Id == id {
                        idx = i
                        break
                }
        }

        if idx == -1 {
                return fmt.Errorf("RemoveCard: Hand with id %i not in hand", idx)
        }

        l := len(p.Board)
        if l >= 2 {
            p.Board[idx] = p.Board[l - 1]
            p.Board = p.Board[:l - 2]
        } else if l == 1 {
            p.Board = make([]*Card, 0)
        }

        fmt.Fprintln(os.Stderr, "Kill monster", id)
        return nil

}

func (p *Player) IncreaseMana() {
	if p.Mana < MAX_MANA {
		p.Mana = p.Mana + 1
		//fmt.Println("Increase Mana")
	}
}

func (g *Game) AddPlayer(p *Player) error {
	if len(g.players) >= MAX_PLAYERS {
		return errors.New("AddPlayer: There is already 2 players")
	}

	g.players = append(g.players, p)
	return nil
}

func (g *Game) GetPlayerRandom() (*Player, error) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	i := random.Intn(len(g.players))
	return g.players[i], nil
}

func (g *Game) OrderPlayer(head_player *Player) error {
	index := -1
	for i, b := range g.players {
		if b == head_player {
			index = i
			break
		}
	}

	g.players[0], g.players[1] = g.players[index], g.players[len(g.players) - index - 1]

	return nil
}

func (g *Game) Hero() (*Player) {
	return g.players[0]
}

func (g *Game) Vilain() (*Player) {
	return g.players[1]
}

func (g *Game) NextPlayer() (*Player, error) {
	p := g.players[0]
	g.players[0], g.players[1] = g.players[1], g.players[0]

	g.Turn++

	return p, nil
}

func (g *Game) ParseAction(actions string) (n int, err error) {
	data := strings.Split(actions, ";")
	for _, a := range data {
		switch s := strings.Split(a, " "); s[0] {
		case "PICK":
			n, err = g.ParseMovePick(s)
		case "ATTACK":
			_ = g.ParseMoveAttack(s)
		case "SUMMON":
			_ = g.ParseMoveSummon(s)
		case "PASS":
		default:
		}
	}
	return -1, err
}

func (g *Game) ParseMovePick(params []string) (int, error) {

	if len(params) != 2 {
		return -1, errors.New("ParseAttack: Format should be SUMMON id1 id2")
	}
	id1, err1 := strconv.ParseInt(params[1], 10, 32)
	if err1 != nil {
		return int(id1), err1
	}

	if id1 >= DRAFT_PICK {
		return int(id1), fmt.Errorf("Wrong pick number")
	}

	return int(id1), nil
}

func (g *Game) ParseMoveSummon(params []string) error {

	if len(params) != 2 {
		return errors.New("ParseAttack: Format should be SUMMON id1 id2")
	}

	id1, err1 := strconv.ParseInt(params[1], 10, 32)
	if err1 != nil {
		return err1
	}

	return g.MoveSummon(int(id1))
}

func (g *Game) ParseMoveAttack(params []string) error {
	if len(params) != 3 {
		return errors.New("ParseAttack: Format should be ACTION id1 id2")
	}

	id1, err1 := strconv.ParseInt(params[1], 10, 32)
	if err1 != nil {
		return err1
	}
	id2, err2 := strconv.ParseInt(params[2], 10, 32)
	if err2 != nil {
		return err2
	}

	switch params[0] {
	case "ATTACK":
		return g.MoveAttack(int(id1), int(id2))
	default:
		err := fmt.Sprintf("ParseAttack: Unknow command %s", params[0])
		return errors.New(err)
	}
}

func (g *Game) MoveSummon(id1 int) error {
	err := g.Hero().HandPlayCard(int(id1))
	if err == nil {
		return err
	}
	//fmt.Println("Summon card", id1, "for player", g.current_player.Id)
	return nil
}

func (g *Game) MoveAttackPolicy(id1, id2 int) (error) {
	guards := g.Vilain().BoardGetGuardsId()
	len_guards := len(guards)

	exist, _ := in_array(guards, id2)
	if len_guards > 0 && (id2 == -1 || ! exist) {
		return fmt.Errorf("Move ATTACK %d %d not permitted", id1, id2)
	} 
	return nil
}
func (g *Game) MoveAttack(id1, id2 int) error {
	c1 := g.Hero().BoardGetCard(int(id1))
	if c1 == nil {
		err := fmt.Sprintf("MoveAttack: Current player %i don't have card %i", g.Hero().Id, id1)
		return errors.New(err)
	}

	err := g.MoveAttackPolicy(id1, id2)
	if err != nil {
		return err
	}

	c1a := c1.Attack
	c1d := c1.Defense

	if id2 == -1 {
		g.Vilain().SetLife(g.Vilain().Life - c1a)
	} else {
		c2 := g.Vilain().BoardGetCard(id2)
		if c2 == nil {
			err := fmt.Sprintf("MoveAttack: Current player %d don't have card %d", g.Vilain().Id, id2)
			return errors.New(err)
		}

		c2a := c2.Attack
		c2d := c2.Defense

		c2.Defense = c2d - c1a
		if c2.Defense <= 0 {
			g.Vilain().BoardRemoveCard(id2)
			if c1.Abilities[CARD_ABILITY_BREAKTHROUGH] != "-" {
				g.Vilain().SetLife(g.Vilain().Life - c2.Defense)
			}
			
		}

		c1.Defense = c1d - c2a
		//fmt.Println("MoveAttack:", c2m, "Repost", c1m, ". Reducing defense ", c2d, "->", c2m.Defense, "for Monster", c2m.Id())
		if c1.Defense <= 0 {
			//fmt.Println("MoveAttack:", c2.Name, "kill", c1.Name)
			g.Hero().BoardRemoveCard(id1)
		}
	}

	return nil
}

func (g *Game) RawPlayers() [][]interface{} {

	raw_data := make([][]interface{}, 0)
	for i := 0 ; i < len(g.players) ; i++ {
		p := g.players[i].Raw()
		raw_data = append(raw_data, p)
	}
	return raw_data
}
func (g *Game) Draft() (error) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := 0; i < DECK_CARDS; i++ {
		draft_raw 	:= make([][]interface{}, DRAFT_PICK)
		draft 		:= make([]*Card, DRAFT_PICK)
		numbers := make([]int, 0)
		for j := 0; j < DRAFT_PICK; j++ {
			num := random.Intn(len(CARDS))
			for exist, _ := in_array(numbers, num); exist; {
				num = random.Intn(len(CARDS))
			}
			c := CARDS[num]
			draft_raw[j] 	= c.Raw()
			draft[j] 	= c
			numbers 	= append(numbers, num)
		}

		for h := 0; h < MAX_PLAYERS ; h++ {
			pick, err := g.Hero().IA.Move(g.RawPlayers(), draft_raw, g.Vilain().Deck.Count(), 100)
			if err != nil {
				return err
			}
			params := strings.Split(pick, " ")
			num, err := g.ParseAction(pick)
			if err != nil {
				return fmt.Errorf("Wrong pick at draft %d", num)
			}
			if num >= 0 {	
				g.Hero().Deck.AddCard(draft[num])
			} else if len(params) == 1 && params[0] == "PASS" {
				num = random.Intn(len(draft))
				g.Hero().Deck.AddCard(draft[num])
			} else {
				return fmt.Errorf("Wrong Action %s", pick)
			}
			_, _ = g.NextPlayer()	

			
		}
	}

	return nil
}

func (g *Game) CheckWinner() *Player {
	if g.Hero().Life <= 0 {
		return g.Vilain()
	} else if g.Vilain().Life <= 0 {
		return g.Hero()
	}
	return nil
}

func (g *Game) CheckDraw() bool {
	return false
}

func (g *Game) Start() (winner *Player, err error) {
	//var round uint32

	winner = nil
	//fmt.Println("Starting Game")
	if len(g.players) != MIN_PLAYERS {
		return nil, errors.New("AddPlayer: There should be players")
	}

	hero, _ := g.GetPlayerRandom()
	//fmt.Println("Starting Player", start_player)
	g.OrderPlayer(hero)

	err = g.Draft()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
/*

	g.Hero().Deck.Shuffle()
	g.Vilain().Deck.Shuffle()

	g.Hero().DrawCardN(4)
	//fmt.Println(p)
	g.Vilain().DrawCardN(5)
	//fmt.Println(p)
/*
	winner = nil
	round = 2

	for winner == nil {
		//fmt.Println("Round:", round / 2)
		g.current_player, err = g.NextPlayer()
		g.opponent = g.players[0]

		if err != nil {
			return nil, err
		}

		g.current_player.IncreaseMana()
		g.current_player.Draw(1)

		//fmt.Println("Current Player", g.current_player)

		actions := g.current_player.Action(g.opponent, round/2)
		//fmt.Println("Action for Player", g.current_player.Id, actions)
		err = g.ParseAction(actions)

		winner = g.CheckWinner()
		if winner != nil {
			break
		} else if g.CheckDraw() {
			break
		}

		round = round + 1
		//		time.Sleep(1 * time.Second)
	}
*/
	return winner, nil
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "ia-binary-1 ia-binary-2")
		os.Exit(1)
	}
	p1 := NewPlayer(1, STARTING_LIFE, STARTING_MANA, NewIAPlayer(1, os.Args[1]))
	p2 := NewPlayer(2, STARTING_LIFE, STARTING_MANA, NewIAPlayer(2, os.Args[2]))
	gm := NewGame()

	gm.AddPlayer(p1)
	gm.AddPlayer(p2)

	gm.Start()
}
