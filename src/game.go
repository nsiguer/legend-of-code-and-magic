package main

import (
	"errors"
	"math/rand"
	"fmt"
	"time"
	"strings"
	"strconv"
	"reflect"

	copier "github.com/jinzhu/copier"
)

const (
	MAX_MANA		= 12
	MAX_PLAYERS		= 2
	MIN_PLAYERS		= 2
	STARTING_MANA	= 1
	STARTING_LIFE	= 30
	DECK_CARDS		= 30
	DRAFT_PICK		= 3
)

var CARDS = []Card {
	NewMonster(1, "A", 1, 2, 2),
	NewMonster(2, "B", 2, 4, 1),
	NewMonster(3, "C", 2, 1, 5),
	NewMonster(4, "D", 2, 2, 3),
	NewMonster(5, "E", 2, 3, 2),
	NewMonster(6, "F", 3, 2, 5),
	NewMonster(7, "G", 3, 3, 4),
	NewMonster(8, "H", 3, 5, 2),
	NewMonster(9, "I", 4, 4, 5),
	NewMonster(10, "J", 4, 1, 8),
	NewMonster(11, "K", 4, 2, 7),
	NewMonster(12, "L", 4, 9, 1),
	NewMonster(13, "M", 4, 6, 2),
	NewMonster(14, "N", 4, 7, 4),
	NewMonster(15, "O", 5, 8, 2),
	NewMonster(16, "P", 5, 5, 6),
	NewMonster(17, "Q", 5, 6, 5),
	NewMonster(18, "R", 6, 7, 5),
	NewMonster(19, "S", 7, 8, 8),
	NewMonster(20, "T", 7, 4, 8),
	NewMonster(21, "U", 9, 10, 10),
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

type Card interface {
	Id ()		uint32
	Name()		string
	Cost()		uint32
}

type Monster struct {
	id		uint32
	name	string
	cost	uint32

	Attack	uint32
	Defense int
	Round	uint32
}

type Deck struct {
	Cards	[]Card
}

type Player struct {
	Id		uint32
	Deck	*Deck
	CDeck	*Deck
	Life	int
	Mana	uint32
	Board	[]Card
	Hand	[]Card
}


type Game struct {
	players			[]*Player
	current_player	*Player
	opponent		*Player
}

func NewMonster(id	uint32,
				name string,
				cost uint32,
				attack uint32,
				defense int,
				) (*Monster)  {
	return &Monster{
		name: name,
		id: id,
		cost: cost,
		Attack: attack,
		Defense: defense,
		Round: 10000,
	}
}

func (m *Monster) Id() (uint32) {
	return m.id
}

func (m *Monster) Name() (string) {
	return m.name
}

func (m *Monster) Cost() (uint32) {
	return m.cost
}


func NewDeck() (*Deck) {
	return &Deck{
		Cards: make([]Card, 0),
	}
}

func (d *Deck) AddCard(c Card) {
	d.Cards = append(d.Cards, c)
}

func (d *Deck) Draw() (Card, error) {
	if len(d.Cards) > 0 {
		c := d.Cards[0]
		d.Cards = d.Cards[1:]
		return c, nil
	} else {
		return nil, errors.New("There is no more card in the deck")
	}
}

func (d *Deck) Shuffle() () {
	source := rand.NewSource(time.Now().UnixNano())
    random := rand.New(source)
    for i := len(d.Cards) - 1; i > 0; i-- {
        j := random.Intn(i + 1)
        d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
    }
}

func NewPlayer(id uint32) (*Player) {
	return &Player{
		Id: id,
		Deck: NewDeck(),
		Life: STARTING_LIFE,
		Board: make([]Card, 0),
		Hand: make([]Card, 0),
	}

}

func (p *Player) Clear() {
	p.Deck	= NewDeck()
	p.Board = make([]Card, 0)
	p.Hand	= make([]Card, 0)
	p.Life	= STARTING_LIFE
}

func (p *Player) SaveDeck() {
	p.CDeck = NewDeck()
	copier.Copy(p.CDeck, p.Deck)
}

func (p *Player) PrintDeck() {
	for _, c := range(p.CDeck.Cards) {
		fmt.Print(c, ";")
	}
	fmt.Println("")
}
func (p *Player) Draw(n uint32) (error) {
	var i uint32

	if p.Deck == nil {
		return errors.New("Player have no Deck")
	}

	for i = 0 ; i < n ; i++ {
		c, e := p.Deck.Draw()
		if e != nil {
			//fmt.Println(e)
			return nil
		}
		p.Hand = append(p.Hand, c)
	}
	return nil
}

func (p *Player) IncreaseMana() {
	if p.Mana < MAX_MANA {
		p.Mana = p.Mana + 1
		//fmt.Println("Increase Mana")
	}
}

func (p *Player) Summon(id uint32) (error) {
	index := -1
    for i, b := range p.Hand {
        if b.Id() == id {
			index = i
			break
        }
    }

	if index == -1 {
		err := fmt.Sprintf("Summon: Hand with id %i not in hand", index)
		return errors.New(err)
	}

	c := p.Hand[index]

	if c.Cost() > p.Mana {
		err := fmt.Sprintf("Summon: Card cost %i and mana is %i", c.Cost(), p.Mana)
		return errors.New(err)
	}

	p.Hand[index] = p.Hand[len(p.Hand)-1]
	p.Hand[len(p.Hand)-1] = nil
	p.Hand = p.Hand[:len(p.Hand)-1]

	p.Board = append(p.Board, c)

	return nil
}

func (p *Player) ReduceLife(damage int) {
	p.Life = p.Life - damage
}

func (p *Player) RemoveCardBoard(id uint32) (error) {
	index := -1
    for i, b := range p.Board {
        if b.Id() == id {
			index = i
			break
        }
    }

	if index == -1 {
		err := fmt.Sprintf("RemoveCard: Hand with id %i not in hand", index)
		return errors.New(err)
	}

	//c := p.Board[index]
	//fmt.Printf("RemoveCardBoard: Player %d remove card %d from board\n", p.Id, c.Id())

	p.Board[index] = p.Board[len(p.Board)-1]
	p.Board[len(p.Board)-1] = nil
	p.Board = p.Board[:len(p.Board)-1]

	return nil

}

func (p *Player) Pick(cards []Card) {
	var new_copy Card
	source := rand.NewSource(time.Now().UnixNano())
    random := rand.New(source)
	if len(cards) >= 0 {
		c := cards[random.Intn(len(cards))]

		switch reflect.TypeOf(c).String() {
		case "*main.Monster":
			new_copy = &Monster{}
		default:
			//fmt.Println("Pick: Unknow type", reflect.TypeOf(c).String())
			new_copy = nil
		}
		if new_copy != nil {
			copier.Copy(new_copy, c)
			p.Deck.AddCard(new_copy)
		}
	}
}

func (p *Player) ActionSummon(opponent *Player, round uint32) []string {
	var summon_cost uint32

	picks		:= make([]string, 0)
	summon_cost = 0

	for _, c := range(p.Hand) {
		if summon_cost >= p.Mana {
			break
		}
		if c.Cost() <= p.Mana {
			type1 := reflect.TypeOf(c)
			if type1.String() == "*main.Monster" {
				cm := c.(*Monster)
				cm.Round = round
			}
			action		:= fmt.Sprintf("SUMMON %d", c.Id())
			//fmt.Println("Action", action, c)
			picks		= append(picks, action)
			summon_cost = summon_cost + c.Cost()
		}
	}
	return picks
}

func (p *Player) ActionAttack(opponent *Player, round uint32) []string {

	attacks		:= make([]string, 0)
	oBoardId	:= []int{-1}

	source := rand.NewSource(time.Now().UnixNano())
    random := rand.New(source)

	for _, c := range(opponent.Board) {
		oBoardId = append(oBoardId, int(c.Id()))
	}

	for _, c := range(p.Board) {
		type1 := reflect.TypeOf(c)
		if type1.String() == "*main.Monster" {
			cm := c.(*Monster)
			if cm.Round < round {
				i := oBoardId[random.Intn(len(oBoardId))]
				action		:= fmt.Sprintf("ATTACK %d %d", c.Id(), i)
				//fmt.Println("Action", action, c)
				attacks		= append(attacks, action)
			} else {
				//fmt.Println("Monster", cm, "can't attack. Summon on turn", round)
			}
		}
	}
	return attacks
}
func (p *Player) Action(opponent *Player, round uint32) string {
	actions := make([]string, 0)
	actions = append(actions, p.ActionSummon(opponent, round)...)
	actions = append(actions, p.ActionAttack(opponent, round)...)
	return strings.Join(actions, ";")
}

func (p *Player) BoardCard(id uint32) (Card) {
    for _, b := range p.Board {
        if b.Id() == id {
            return b
        }
    }
    return nil
}

func (p *Player) HandCard(id uint32) (Card) {
    for _, b := range p.Hand {
        if b.Id() == id {
            return b
        }
    }
    return nil
}


func NewGame() (*Game) {
	return &Game{
		players: make([]*Player, 0),
	}
}

func (g *Game) AddPlayer(p *Player) (error) {
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

func (g *Game) OrderPlayer(head_player *Player) (error) {
	index := -1
    for i, b := range g.players {
        if b == head_player {
			index = i
			break
        }
    }

	if index == -1 {
		err := fmt.Sprintf("OrderPlayer: Player not found")
		return errors.New(err)
	}

	if index != 0 {
		p := g.players[index]
		f := g.players[0]
		g.players[0] = p
		g.players[index] = f
	}

	g.current_player, _		= g.NextPlayer()
	g.opponent				= g.players[0]

	return nil
}

func (g *Game) NextPlayer() (*Player, error) {
	l := len(g.players)
	if l == 0 {
		return nil, errors.New("NextPlayer: There is 0 players in the game")
	} else if l == 1 {
		return g.players[0], nil
	}
	p := g.players[0]
	g.players = append(g.players[1:], p)

	return p, nil
}

func (g *Game) ParseAction(actions string) (err error) {
	data := strings.Split(actions, ";")
	for _, a := range(data) {
		switch s := strings.Split(a, " ") ; s[0] {
		case "ATTACK":
			err = g.ParseActionAttack(s)
		case "SUMMON":
			err = g.ParseActionSummon(s)
		default:
			err_str := fmt.Sprintf("ParseAction: Unknow action %s", s[0])
			err = errors.New(err_str)
		}
		if err != nil {
			//fmt.Println(err)
		}
	}
	return nil
}
func (g *Game) ParseActionSummon(params []string) (error) {

	if len(params) != 2 {
		return errors.New("ParseAttack: Format should be SUMMON id1 id2")
	}

	id1, err1 := strconv.ParseInt(params[1], 10, 32)
	if err1 != nil {
		return err1
	}

	return g.ActionSummon(int(id1))
}
func (g *Game) ParseActionAttack(params []string) (error) {
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
		return g.ActionAttack(int(id1), int(id2))
	default:
		err := fmt.Sprintf("ParseAttack: Unknow command %s", params[0])
		return errors.New(err)
	}
}

func (g *Game) ActionSummon(id1 int) (error) {
	err := g.current_player.Summon(uint32(id1))
	if err == nil {
		return err
	}
	//fmt.Println("Summon card", id1, "for player", g.current_player.Id)
	return nil
}
func (g *Game) ActionAttack(id1, id2 int) (error) {
	c1 := g.current_player.BoardCard(uint32(id1))
	if c1 == nil {
		err := fmt.Sprintf("ActionAttack: Current player %i don't have card %i", g.current_player.Id, id1)
		return errors.New(err)
	}

	type1 := reflect.TypeOf(c1)
	if type1.String() != "*main.Monster" {
		err := fmt.Sprintf("ActionAttack: Card %i is not a Monster (%s)", id1, type1)
		return errors.New(err)
	}

	c1m := c1.(*Monster)
	c1a := c1m.Attack
	c1d := c1m.Defense

	if id2 == -1 {
		//pl := g.opponent.Life
		g.opponent.ReduceLife(int(c1a))
		//fmt.Println("ActionAttack: Reducing life", pl, "->", g.opponent.Life)
	} else {
		c2 := g.opponent.BoardCard(uint32(id2))
		if c2 == nil {
			err := fmt.Sprintf("ActionAttack: Current player %d don't have card %d", g.opponent.Id, id2)
			return errors.New(err)
		}

		type2 := reflect.TypeOf(c2)
		if type2.String() != "*main.Monster" {
			err := fmt.Sprintf("ActionAttack: Card %d is not a Monster (%s)", id2, type2)
			return errors.New(err)
		}
		c2m := c2.(*Monster)
		c2a := c2m.Attack
		c2d := c2m.Defense

		c2m.Defense = c2d - int(c1a)
		//fmt.Println("ActionAttack:", c1m, "attack", c2m, ". Reducing defense ", c2d, "->", c2m.Defense, "for Monster", c2m.Id())
		if c2m.Defense <= 0 {
			//fmt.Println("ActionAttack:", c1.Name, "kill", c2.Name)
			g.opponent.RemoveCardBoard(uint32(id2))
		}

		c1m.Defense = c1d - int(c2a)
		//fmt.Println("ActionAttack:", c2m, "Repost", c1m, ". Reducing defense ", c2d, "->", c2m.Defense, "for Monster", c2m.Id())
		if c1m.Defense <= 0 {
			//fmt.Println("ActionAttack:", c2.Name, "kill", c1.Name)
			g.current_player.RemoveCardBoard(uint32(id1))
		}
	}

	return nil
}

func (g *Game) Draft() {
	source := rand.NewSource(time.Now().UnixNano())
    random := rand.New(source)
    for i := 0 ; i < DECK_CARDS ; i++ {
		draft := make([]Card, DRAFT_PICK)
		numbers := make([]int, 0)
		for j := 0 ; j < DRAFT_PICK ; j++ {
			num := random.Intn(len(CARDS))
			for exist, _ := in_array(numbers, num) ; exist ; {
				num = random.Intn(len(CARDS))
			}
			c := CARDS[num]
			draft[j] = c
			numbers = append(numbers, num)
		}
		g.current_player.Pick(draft)
		g.opponent.Pick(draft)
    }

	g.current_player.SaveDeck()
	g.opponent.SaveDeck()

}

func (g *Game) CheckWinner() (*Player) {
	if g.current_player.Life <= 0 {
		return g.opponent
	} else if g.opponent.Life <= 0 {
		return g.current_player
	}
	return nil
}

func (g *Game) CheckDraw() (bool) {
	if len(g.current_player.Board) == 0 &&
	   len(g.current_player.Deck.Cards) == 0 &&
	   len(g.opponent.Board) == 0 &&
	   len(g.opponent.Deck.Cards) == 0 {
		return true
	}
	return false
}

func (g *Game) Clear() {
	if g.current_player != nil {
		g.current_player.Clear()
	}
	if g.opponent != nil {
		g.opponent.Clear()
	}
}

func (g *Game) Start() (winner *Player, err error) {
	var round uint32
	//fmt.Println("Starting Game")
	if len(g.players) != MIN_PLAYERS {
		return nil, errors.New("AddPlayer: There should be players")
	}

	start_player, _ := g.GetPlayerRandom()
	//fmt.Println("Starting Player", start_player)
	g.OrderPlayer(start_player)

	g.Draft()

	p, _ := g.NextPlayer() ; p.Draw(4)
	//fmt.Println(p)
	p, _ = g.NextPlayer() ; p.Draw(5)
	//fmt.Println(p)

	winner = nil
	round  = 2

	for winner == nil {
		//fmt.Println("Round:", round / 2)
		g.current_player, err	= g.NextPlayer()
		g.opponent				= g.players[0]

		if err != nil {
			return nil, err
		}


		g.current_player.IncreaseMana()
		g.current_player.Draw(1)

		//fmt.Println("Current Player", g.current_player)

		actions := g.current_player.Action(g.opponent, round / 2)
		//fmt.Println("Action for Player", g.current_player.Id, actions)
		err	= g.ParseAction(actions)

		winner = g.CheckWinner()
		if winner != nil {
			break
		} else if g.CheckDraw() {
			break
		}

		round = round + 1
//		time.Sleep(1 * time.Second)
	}
	return winner, nil
}

