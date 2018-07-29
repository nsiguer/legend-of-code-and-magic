package main


import (
//        "errors"
        "fmt"
        "math/rand"
        "reflect"
        "strings"
        "time"
        "os"
        "strconv"
        "sort"

///        copier "github.com/jinzhu/copier"
)

const (
        CARD_TYPE_MONSTER           = 0
        CARD_ABILITY_BREAKTHROUGH   = 0
        CARD_ABILITY_GUARD          = 3
        CARD_ABILITY_CHARGE         = 1
)

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

type Card interface {
        Id()    int
        Cost()  int
        Type()  int
}

type Monster struct {
        id      int
        cost    int
        type_   int

        Attack  int
        Defense int
        Celerity int
        
        Abilities []string
}


type Player struct {
        Id      int
        Life    int
        Mana    int
        Board   []Card
        Hand    []Card

        s_life  int
        s_mana  int
}

func NewMonster(id int,
                type_ int,
                cost int,
                attack int,
                defense int,
                abilities string,
) *Monster {
        return &Monster{
                id:             id,
                type_:          type_,
                cost:           cost,
                Attack:         attack,
                Defense:        defense,
                Celerity:       0,
                Abilities:      strings.Split(abilities, ""),
        }
}

func (m *Monster) Id() int {
        return m.id
}


func (m *Monster) Cost() int {
        return m.cost
}

func (m *Monster) Type() int {
        return m.type_
}


func NewPlayer(id, life, mana int) *Player {
        return &Player{
                Id:     id,
                Mana:   mana,
                Life:   life,
                s_life: life,
                s_mana: mana,
                Board:  make([]Card, 0),
                Hand:   make([]Card, 0),
        }

}

func (p *Player) Clear() {
        p.Board = make([]Card, 0)
        p.Hand  = make([]Card, 0)
        p.Life  = p.s_life
        p.Mana  = p.s_mana
}


func (p *Player) Draw(c Card) {
    if p.HandGetCard(c.Id()) == nil {
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

func (p *Player) HandGetCard(id int) Card {
        for _, c := range(p.Hand) {
                if c.Id() == id {
                        return c
                }
        }
        return nil
}

func (p *Player) HandRemoveCard(id int) (error) {
        idx := -1
        for i, c := range(p.Hand) {
                if c.Id() == id {
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
            p.Hand = make([]Card, 0)
        }
        return nil
}
func (p *Player) HandPlayCard(id int) (error) {
        c := p.HandGetCard(id)
        if c == nil {
                return fmt.Errorf("Card with id %d is not present in hand of player %d", id, p.Id)
        }
        switch c.Type() {
        case CARD_TYPE_MONSTER:
             
                fmt.Fprintln(os.Stderr, p.Id, "Board", p.Board)
                fmt.Fprintln(os.Stderr, p.Id, "Hand", p.Hand)
                p.Board = append(p.Board, c)
                p.HandRemoveCard(id)
                fmt.Fprintln(os.Stderr, p.Id, "Board", p.Board)
                fmt.Fprintln(os.Stderr, p.Id, "Hand", p.Hand) 
                  
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
                if c.Id() == id {
                    switch c.Type() {
                    case CARD_TYPE_MONSTER:
                        cm := c.(*Monster)
                        fmt.Fprintln(os.Stderr, "Reduce defense of", c.Id(), cm.Defense, " -> ", cm.Defense - damage)
                        cm.Defense -= damage
                       
                        if cm.Defense <= 0 {
                            p.BoardRemoveCard(c.Id())
                            return true
                        }
                        
                    }    
                    return false
                }
        }
        
        return false
}
func (p *Player) UpdateBoard(card Card) bool {
        
        for _, c := range(p.Board) {
                if c.Id() == card.Id() {
                    switch c.Type() {
                    case CARD_TYPE_MONSTER:
                        cm := c.(*Monster)
                        cardm := card.(*Monster)
                        cm.Celerity = 1
                        fmt.Fprintln(os.Stderr, "Reduce defense of", c.Id(), cm.Defense, " -> ", cardm.Defense)
                        cm.Defense = cardm.Defense
                       
                        if cm.Defense <= 0 {
                            p.BoardRemoveCard(c.Id())
                        }
                    }    
                    return true
                }
        }
        fmt.Fprintln(os.Stderr, p.Id, "Board add monster", card)
        p.Board = append(p.Board, card)
        return false
}
func (p *Player) BoardGetCard(id int) Card {
        for _, c := range(p.Board) {
                if c.Id() == id {
                        return c
                }
        }
        return nil
}

func (p *Player) BoardRemoveCard(id int) error {
        idx := -1
        for i, b := range p.Board {
                if b.Id() == id {
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
            p.Board = make([]Card, 0)
        }

        fmt.Fprintln(os.Stderr, "Kill monster", id)
        return nil

}

func (p *Player) Pick(cards []Card) {
        source := rand.NewSource(time.Now().UnixNano())
        random := rand.New(source)

        fmt.Fprintln(os.Stderr, "PICK BETWEEN", cards)        
        if len(cards) >= 0 {
                r := random.Intn(len(cards))
                fmt.Println("PICK", r)
        }
}


type ByCost []Card

func (a ByCost) Len() int           { return len(a) }
func (a ByCost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCost) Less(i, j int) bool { return a[i].Cost() < a[j].Cost() }


func (p *Player) SortHand() {
       sort.Sort(ByCost(p.Hand))
}
func (p *Player) ActionSummon(opponent *Player, round int) []string {
        var summon_cost int

        picks := make([]string, 0)
        p.SortHand()
        for i, c := range p.Hand {
                summon_cost = 0
                tmp_picks := make([]string, 0)
                
                summon_cost += c.Cost()
                if summon_cost >= p.Mana {
                    continue
                }   
                action := fmt.Sprintf("SUMMON %d", c.Id())
                tmp_picks = append(tmp_picks, action)
                
                for j := 0 ; j < len(p.Hand) ; j ++ { 
                    if i == j {
                        continue
                    }
                    if summon_cost + p.Hand[j].Cost() > p.Mana {
                        continue
                    }
                  
                    switch c.Type() {
                    case CARD_TYPE_MONSTER:
                        //fmt.Println("Action", action, c)
                            action := fmt.Sprintf("SUMMON %d", p.Hand[j].Id())
                            tmp_picks = append(tmp_picks, action)
                            summon_cost = summon_cost + p.Hand[j].Cost()
                    }
                }
                picks = append(picks, strings.Join(tmp_picks, ";"))
        }
        fmt.Fprintln(os.Stderr, "PICKS", picks)
        return picks
}

func (p *Player) AvailableMovesSummon(opponent *Player) []string {
                
        var summon_cost int

        picks := make([]string, 0)
        p.SortHand()
        for i, c := range p.Hand {
                summon_cost = 0
                tmp_picks := make([]string, 0)
                
                summon_cost += c.Cost()
                if summon_cost > p.Mana {
                    continue
                }   
                fmt.Fprintf(os.Stderr, "Monster %d good for summon", c.Id())
                action := fmt.Sprintf("SUMMON %d", c.Id())
                tmp_picks = append(tmp_picks, action)
                
                for j := 0 ; j < len(p.Hand) ; j ++ { 
                    if i == j {
                        continue
                    }
                    if summon_cost + p.Hand[j].Cost() > p.Mana {
                        continue
                    }
                  
                    switch c.Type() {
                    case CARD_TYPE_MONSTER:
                        //fmt.Println("Action", action, c)
                            action := fmt.Sprintf("SUMMON %d", p.Hand[j].Id())
                            tmp_picks = append(tmp_picks, action)
                            summon_cost = summon_cost + p.Hand[j].Cost()
                    }
                }
                picks = append(picks, strings.Join(tmp_picks, ";"))
        }
        fmt.Fprintln(os.Stderr, "PICKS", picks)
        return picks
}

func (p *Player) BoardGetGuards() []int {
    guards := make([]int, 0)
    for _, c := range(p.Board) {
        cm := c.(*Monster)
        if cm.Abilities[CARD_ABILITY_GUARD] != "-" {
            guards = append(guards, cm.Id())
        }
    }
    return guards
}
func (p *Player) AvailableMovesAttack(opponent *Player) []string {

        oBoardId := []int{-1}

        source := rand.NewSource(time.Now().UnixNano())
        random := rand.New(source)

        for _, c := range opponent.Board {
                oBoardId = append(oBoardId, c.Id())
        }

     
        oBoardGuardId := opponent.BoardGetGuards()
        
        fmt.Fprintln(os.Stderr, "Opponent Guard", oBoardGuardId)
        tmp_actions := make([]string, 0)
        for _, c := range p.Board {
                var i int
                switch c.Type() {
                case CARD_TYPE_MONSTER:
                        cm := c.(*Monster)
                        fmt.Fprintln(os.Stderr, "CELERITY OF", cm.Id(), "IS", cm.Celerity)
                        len_guard_id := len(oBoardGuardId)
                        if cm.Celerity > 0 {           
                                if len_guard_id > 0 {
                                    i = oBoardGuardId[0]
                                } else {
                                    i = oBoardId[random.Intn(len(oBoardId))]
                                }
                                if i != -1 {
                                    // Attack kill the opponent monster
                                    if opponent.BoardAttackCard(i, cm.Attack) && len_guard_id > 0 {
                                        if len_guard_id > 2 {
                                            oBoardGuardId = oBoardGuardId[1:]
                                        } else {
                                            oBoardGuardId = make([]int, 0)
                                        }
                                        
                                        
                                    }
                                    
                                }
                                action_str := fmt.Sprintf("ATTACK %d %d", c.Id(), i)
                                //fmt.Println("Action", action, c)
                                tmp_actions = append(tmp_actions, action_str)
                        } else {
                                //fmt.Println("Monster", cm, "can't attack. Summon on turn", round)
                        }
                }
        }
        return []string{ strings.Join(tmp_actions, ";") }
}
func (p *Player) Action(opponent *Player) {
        var avms, avma []string
        
        actions := make([]string, 0)
        
        source := rand.NewSource(time.Now().UnixNano())
        random := rand.New(source)
        
        avms = p.AvailableMovesSummon(opponent)
        if len(avms) > 0 {
            actions = append(actions, avms[random.Intn(len(avms))])
            for _, astr := range(actions) {
                   a := strings.Split(astr, " ")
                   n, _ := strconv.Atoi(a[1])
                   p.HandPlayCard(n)
            }
        }
        avma = p.AvailableMovesAttack(opponent)
        if len(avma) > 0 {
            actions = append(actions, avma[random.Intn(len(avma))])
        }
        fmt.Fprintln(os.Stderr, actions)
        if len(actions) == 0 {
            fmt.Println("PASS")
        } else {
            fmt.Println(strings.Join(actions, ";"))
        }
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/


func main() {
 

    var previous_hero, previous_vilain *Player
    
    for {

    var players []*Player
    var hero, vilain *Player
    var step_draft bool

    players = make([]*Player, 2)
    players[0] = NewPlayer(0, 0, 0)
    players[1] = NewPlayer(1, 0, 0)

    hero    = players[0]
    vilain  = players[1]
    step_draft = false

        for i := 0; i < 2; i++ {
            var playerHealth, playerMana, playerDeck, playerRune int
            fmt.Scan(&playerHealth, &playerMana, &playerDeck, &playerRune)

            
            players[i].Life = playerHealth
            players[i].Mana = playerMana

        }



        var opponentHand int
        fmt.Scan(&opponentHand)

        var cardCount int
        fmt.Scan(&cardCount)

        var draft []Card = make([]Card, 0)

        for i := 0; i < cardCount; i++ {
            var cardNumber, instanceId, location, cardType, cost, attack, defense int
            var abilities string
            var myHealthChange, opponentHealthChange, cardDraw int
            fmt.Scan(&cardNumber, &instanceId, &location, &cardType, &cost, &attack, &defense, &abilities, &myHealthChange, &opponentHealthChange, &cardDraw)
         
            fmt.Fprintln(os.Stderr, cardNumber, instanceId, location, cardType, cost, attack, defense, abilities, myHealthChange, opponentHealthChange, cardDraw)
            var card Card
            switch cardType {

            case CARD_TYPE_MONSTER:
                card = NewMonster(instanceId, cardType, cost, attack, defense, abilities)
                if (card.(*Monster)).Abilities[CARD_ABILITY_CHARGE] != "-" {
                    (card.(*Monster)).Celerity = 1
                }
            default:
                card = nil

            }
            fmt.Fprintln(os.Stderr, card)

            // Draft step
            if instanceId == -1 {
                step_draft = true
                draft = append(draft, card)

            } else {
                step_draft = false
            
                switch location {
                case -1:
                        if previous_vilain != nil && previous_vilain.BoardGetCard(card.Id()) != nil {
                            (card.(*Monster)).Celerity = 1
                        }
                        vilain.UpdateBoard(card)
                case 0:
                        hero.Draw(card)
                case 1:
                        if previous_hero != nil && previous_hero.BoardGetCard(card.Id()) != nil {
                            (card.(*Monster)).Celerity = 1
                        }
                        hero.UpdateBoard(card)
                }
            }
        }

        if step_draft {
                hero.Pick(draft)
        } else {
                hero.Action(vilain)
                fmt.Fprintln(os.Stderr, "Board after action", hero)
        }
        // fmt.Fprintln(os.Stderr, "Debug messages...")
        previous_hero   = hero
        previous_vilain = vilain
        fmt.Fprintln(os.Stderr, "Previous vilain", previous_vilain)
        
    }
}

