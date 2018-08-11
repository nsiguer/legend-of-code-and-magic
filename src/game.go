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

	MOVE_TIMEOUT	= 100000
	MAX_MANA		= 12
	MAX_PLAYERS		= 2
	MIN_PLAYERS		= 2
	MAX_HAND_CARD	= 8
	MAX_BOARD_CARD	= 6

	STARTING_MANA	= 0
	STARTING_LIFE	= 30
	STARTING_CARD	= 30
	STARTING_RUNES	= 5
	DECK_CARDS		= 30

	DRAFT_PICK		= 3

	CARD_TYPE_CREATURE	        = 0
	CARD_TYPE_ITEM_GREEN	    = 1
	CARD_TYPE_ITEM_RED	        = 2
	CARD_TYPE_ITEM_BLUE	        = 3

    CARD_ABILITY_BREAKTHROUGH   = 0
	CARD_ABILITY_CHARGE         = 1
	CARD_ABILITY_DRAIN       	= 2
	CARD_ABILITY_GUARD          = 3
	CARD_ABILITY_LETHAL         = 4
	CARD_ABILITY_WARD           = 5

)

var CARDS = []*Card{
	NewCard(1, -1, 1, 2, 1, "------", 1, 0, 0, CARD_TYPE_CREATURE),
	NewCard(2, -1, 1, 1, 2, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(3, -1, 1, 2, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(4, -1, 2, 1, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(5, -1, 2, 4, 1, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(6, -1, 2, 3, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(7, -1, 2, 2, 2, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(8, -1, 2, 2, 3, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(9, -1, 3, 3, 4, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(10, -1, 3, 3, 1, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(11, -1, 3, 5, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(12, -1, 3, 2, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(13, -1, 4, 5, 3, "------", 1, -1, 0, CARD_TYPE_CREATURE),
	NewCard(14, -1, 4, 9, 1, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(15, -1, 4, 4, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(16, -1, 4, 6, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(17, -1, 4, 4, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(18, -1, 4, 7, 4, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(19, -1, 5, 5, 6, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(20, -1, 5, 8, 2, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(21, -1, 5, 6, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(22, -1, 6, 7, 5, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(23, -1, 7, 8, 8, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(24, -1, 1, 1, 1, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(25, -1, 2, 3, 1, "------", -2, -2, 0, CARD_TYPE_CREATURE),
	NewCard(26, -1, 2, 3, 2, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(27, -1, 2, 2, 2, "------", 2, 0, 0, CARD_TYPE_CREATURE),
	NewCard(28, -1, 2, 1, 2, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(29, -1, 2, 2, 1, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(30, -1, 3, 4, 2, "------", 0, -2, 0, CARD_TYPE_CREATURE),
	NewCard(31, -1, 3, 3, 1, "------", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(32, -1, 3, 3, 2, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(33, -1, 4, 4, 3, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(34, -1, 5, 3, 5, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(35, -1, 6, 5, 2, "B-----", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(36, -1, 6, 4, 4, "------", 0, 0, 2, CARD_TYPE_CREATURE),
	NewCard(37, -1, 6, 5, 7, "------", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(38, -1, 1, 1, 3, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(39, -1, 1, 2, 1, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(40, -1, 3, 2, 3, "--DG--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(41, -1, 3, 2, 2, "-CD---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(42, -1, 4, 4, 2, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(43, -1, 6, 5, 5, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(44, -1, 6, 3, 7, "--D-L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(45, -1, 6, 6, 5, "B-D---", -3, 0, 0, CARD_TYPE_CREATURE),
	NewCard(46, -1, 9, 7, 7, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(47, -1, 2, 1, 5, "--D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(48, -1, 1, 1, 1, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(49, -1, 2, 1, 2, "---GL-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(50, -1, 3, 3, 2, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(51, -1, 4, 3, 5, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(52, -1, 4, 2, 4, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(53, -1, 4, 1, 1, "-C--L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(54, -1, 3, 2, 2, "----L-", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(55, -1, 2, 0, 5, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(56, -1, 4, 2, 7, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(57, -1, 4, 1, 8, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(58, -1, 6, 5, 6, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(59, -1, 7, 7, 7, "------", 1, -1, 0, CARD_TYPE_CREATURE),
	NewCard(60, -1, 7, 4, 8, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(61, -1, 9, 10, 10, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(62, -1, 12, 12, 12, "B--G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(63, -1, 2, 0, 4, "---G-W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(64, -1, 2, 1, 1, "---G-W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(65, -1, 2, 2, 2, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(66, -1, 5, 5, 1, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(67, -1, 6, 5, 5, "-----W", 0, -2, 0, CARD_TYPE_CREATURE),
	NewCard(68, -1, 6, 7, 5, "-----W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(69, -1, 3, 4, 4, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(70, -1, 4, 6, 3, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(71, -1, 4, 3, 2, "BC----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(72, -1, 4, 5, 3, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(73, -1, 4, 4, 4, "B-----", 4, 0, 0, CARD_TYPE_CREATURE),
	NewCard(74, -1, 5, 5, 4, "B--G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(75, -1, 5, 6, 5, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(76, -1, 6, 5, 5, "B-D---", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(77, -1, 7, 7, 7, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(78, -1, 8, 5, 5, "B-----", 0, -5, 0, CARD_TYPE_CREATURE),
	NewCard(79, -1, 8, 8, 8, "B-----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(80, -1, 8, 8, 8, "B--G--", 0, 0, 1, CARD_TYPE_CREATURE),
	NewCard(81, -1, 9, 6, 6, "BC----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(82, -1, 7, 5, 5, "B-D--W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(83, -1, 0, 1, 1, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(84, -1, 2, 1, 1, "-CD--W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(85, -1, 3, 2, 3, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(86, -1, 3, 1, 5, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(87, -1, 4, 2, 5, "-C-G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(88, -1, 5, 4, 4, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(89, -1, 5, 4, 1, "-C----", 2, 0, 0, CARD_TYPE_CREATURE),
	NewCard(90, -1, 8, 5, 5, "-C----", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(91, -1, 0, 1, 2, "---G--", 0, 1, 0, CARD_TYPE_CREATURE),
	NewCard(92, -1, 1, 0, 1, "---G--", 2, 0, 0, CARD_TYPE_CREATURE),
	NewCard(93, -1, 1, 2, 1, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(94, -1, 2, 1, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(95, -1, 2, 2, 3, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(96, -1, 2, 3, 2, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(97, -1, 3, 3, 3, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(98, -1, 3, 2, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(99, -1, 3, 2, 5, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(100, -1, 3, 1, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(101, -1, 4, 3, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(102, -1, 4, 3, 3, "---G--", 0, -1, 0, CARD_TYPE_CREATURE),
	NewCard(103, -1, 4, 3, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(104, -1, 4, 4, 4, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(105, -1, 5, 4, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(106, -1, 5, 5, 5, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(107, -1, 5, 3, 3, "---G--", 3, 0, 0, CARD_TYPE_CREATURE),
	NewCard(108, -1, 5, 2, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(109, -1, 5, 5, 6, "------", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(110, -1, 5, 0, 9, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(111, -1, 6, 6, 6, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(112, -1, 6, 4, 7, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(113, -1, 6, 2, 4, "---G--", 4, 0, 0, CARD_TYPE_CREATURE),
	NewCard(114, -1, 7, 7, 7, "---G--", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(115, -1, 8, 5, 5, "---G-W", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(116, -1, 12, 8, 8, "BCDGLW", 0, 0, 0, CARD_TYPE_CREATURE),
	NewCard(117, -1, 1, 1, 1, "B-----", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(118, -1, 0, 0, 3, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(119, -1, 1, 1, 2, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(120, -1, 2, 1, 0, "----L-", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(121, -1, 2, 0, 3, "------", 0, 0, 1, CARD_TYPE_ITEM_GREEN),
	NewCard(122, -1, 2, 1, 3, "---G--", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(123, -1, 2, 4, 0, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(124, -1, 3, 2, 1, "--D---", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(125, -1, 3, 1, 4, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(126, -1, 3, 2, 3, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(127, -1, 3, 0, 6, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(128, -1, 4, 4, 3, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(129, -1, 4, 2, 5, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(130, -1, 4, 0, 6, "------", 4, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(131, -1, 4, 4, 1, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(132, -1, 5, 3, 3, "B-----", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(133, -1, 5, 4, 0, "-----W", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(134, -1, 4, 2, 2, "------", 0, 0, 1, CARD_TYPE_ITEM_GREEN),
	NewCard(135, -1, 6, 5, 5, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(136, -1, 0, 1, 1, "------", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(137, -1, 2, 0, 0, "-----W", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(138, -1, 2, 0, 0, "---G--", 0, 0, 1, CARD_TYPE_ITEM_GREEN),
	NewCard(139, -1, 4, 0, 0, "----LW", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(140, -1, 2, 0, 0, "-C----", 0, 0, 0, CARD_TYPE_ITEM_GREEN),
	NewCard(141, -1, 0, -1, -1, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(142, -1, 0, 0, 0, "BCDGLW", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(143, -1, 0, 0, 0, "---G--", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(144, -1, 1, 0, -2, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(145, -1, 3, -2, -2, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(146, -1, 4, -2, -2, "------", 0, -2, 0, CARD_TYPE_ITEM_RED),
	NewCard(147, -1, 2, 0, -1, "------", 0, 0, 1, CARD_TYPE_ITEM_RED),
	NewCard(148, -1, 2, 0, -2, "BCDGLW", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(149, -1, 3, 0, 0, "BCDGLW", 0, 0, 1, CARD_TYPE_ITEM_RED),
	NewCard(150, -1, 2, 0, -3, "------", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(151, -1, 5, 0, -99, "BCDGLW", 0, 0, 0, CARD_TYPE_ITEM_RED),
	NewCard(152, -1, 7, 0, -7, "------", 0, 0, 1, CARD_TYPE_ITEM_RED),
	NewCard(153, -1, 2, 0, 0, "------", 5, 0, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(154, -1, 2, 0, 0, "------", 0, -2, 1, CARD_TYPE_ITEM_BLUE),
	NewCard(155, -1, 3, 0, -3, "------", 0, -1, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(156, -1, 3, 0, 0, "------", 3, -3, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(157, -1, 3, 0, -1, "------", 1, 0, 1, CARD_TYPE_ITEM_BLUE),
	NewCard(158, -1, 3, 0, -4, "------", 0, 0, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(159, -1, 4, 0, -3, "------", 3, 0, 0, CARD_TYPE_ITEM_BLUE),
	NewCard(160, -1, 2, 0, 0, "------", 2, -2, 0, CARD_TYPE_ITEM_BLUE),
}

var CARDS_COUNT = len(CARDS)

func PickRandomCard() (*Card) {
	source	:= rand.NewSource(time.Now().UnixNano())
	random	:= rand.New(source)
	idx	:= random.Intn(CARDS_COUNT)
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
	CardNumber 				int
    Id    					int
	Location				int
    Type  					int
    Cost  					int
	Attack					int
	Defense					int
	Abilities				[]string
	HealthChange			int
	OpponentHealthChange 	int
	CardDraw				int
	

}
type IAPlayer struct {
	id	int
	binary	string
	cmd	*exec.Cmd
	stdin	io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser

	action  chan string
}
type Player struct {
	Id    int
	Deck  *Deck
	Life  int
	Mana  int
	Board []*Card
	Hand  []*Card
	Runes	int

	Graveyard []*Card
	stack_draw int
	current_mana	int
	IA	*IAPlayer
}
type Game struct {
	players        []*Player
	Turn		int
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
		action: make(chan string),
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(ia.stdout)
	scanner := bufio.NewScanner(reader)
	
	readerErr := bufio.NewReader(ia.stderr)
	scannerErr := bufio.NewScanner(readerErr)

	go func() {
		for scanner.Scan() {
			//txt := scanner.Text()
			fmt.Println("[GAME] Player (", ia.id, ") move:", scanner.Text())
			ia.action<- scanner.Text()
		}
	} ()

	go func() {
		for scannerErr.Scan() {

		  fmt.Println("STDERR: (", ia.id, ")", scannerErr.Text())
		}
	} ()
	return ia
}

func (ia *IAPlayer) ReadMove(ms int) (string, error) {
	var a string
    start := time.Now()

    
	select {
	case <- time.After(time.Nanosecond * time.Duration(ms) * (1000000)):
		return "", fmt.Errorf("Timeout after %d ms", ms)
	case a = <-ia.action:
		break
	}
	
	elapsed := time.Since(start)
	fmt.Println("[GAME] Action took", elapsed)

	return a, nil
}

func (ia *IAPlayer) WriteData(data string) (error) {
	//fmt.Println("STDIN : (", ia.id, ")", data)
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
	    	cost int,
        	attack int,
			defense int,
			abilities string,
	     	heroHealthChange int,
	     	opponentHealthChange int,
			cardDraw int,
			type_ int, 
) *Card {
        return &Card{
				CardNumber: 	cardNumber,
                Id:             id,
				Location:		0,
                Type:           type_,
                Cost:           cost,
                Attack:         attack,
                Defense:        defense,
				Abilities:      strings.Split(abilities, ""),
				HealthChange: 	heroHealthChange,
				OpponentHealthChange: opponentHealthChange,
				CardDraw:		cardDraw,
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

func (c *Card) Copy() *Card {
	new_c := &Card{
		CardNumber: 	c.CardNumber,
		Id:             c.Id,
		Location:		c.Location,
		Type:           c.Type,
		Cost:           c.Cost,
		Attack:         c.Attack,
		Defense:        c.Defense,
		HealthChange: 	c.HealthChange,
		Abilities:		make([]string, len(c.Abilities)),
		OpponentHealthChange: c.OpponentHealthChange,
		CardDraw:		c.CardDraw,
	}

	copy(new_c.Abilities, c.Abilities)
	return new_c
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

func NewPlayer(id, life, mana, runes int, ia *IAPlayer) *Player {
        return &Player{
                Id:     id,
                Mana:   mana,
                Life:   life,
		Deck: 	NewDeck(),
                Board:  make([]*Card, 0),
				Hand:   make([]*Card, 0),
				Graveyard: make([]*Card, 0),
				Runes: runes,
				stack_draw: 0,
				IA: ia,
        }

}


func (p* Player) ReloadMana () {
	p.current_mana = p.Mana
}
func (p *Player) Raw() []interface{} {
	return []interface{}{
		p.Life,
		p.Mana,
		p.Deck.Count(),
		p.Runes,
	}
}
func (p *Player) LoseLifeToNextRune() {
	if p.Runes > 0 && p.Life >= STARTING_RUNES * p.Runes {
		var damage int
		damage = p.Life - (STARTING_RUNES  * (p.Runes - 1))
		p.Runes -= 1
		fmt.Println("[GAME][RUNE] Player", p.Id, "can't draw card. Losing", damage, "damage")
		p.ReceiveDamage(damage)
	} else {
		fmt.Println("[GAME][RUNE] Player", p.Id, "can't draw card and have no more Rune")
		p.ReceiveDamage(p.Life)
	}
}
func (p *Player) DrawCard() (error) {
	if len(p.Hand) >= MAX_HAND_CARD && p.Deck.Count() > 0 {
		fmt.Println("[GAME][DECK] Maximum card hand reach", MAX_HAND_CARD)
		return fmt.Errorf("[GAME][DECK] Maximum card hand reach %d", MAX_HAND_CARD)
	}
    c, err := p.Deck.Draw()
    if err == nil {
		fmt.Println("[GAME][DECK]: Player", p.Id, "draw card", c)
		p.Draw(c)
    } else {
		p.LoseLifeToNextRune()
		fmt.Println("[GAME][INFO]:", err)
		return err
	}

	return nil
}

func (p *Player) DrawCardN(n int) (err error) {
    for i := 0 ; i < n ; i++ {
		err = p.DrawCard()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Player) Draw(c *Card) {
    if p.HandGetCard(c.Id) == nil {
		if len(p.Hand) < MAX_HAND_CARD {
			p.Hand = append(p.Hand, c)
		} else {
			fmt.Println("GAME: Max hand card reach", MAX_HAND_CARD)
		}

    } else {
		fmt.Println("GAME: Card", c, "already exist in", p.Hand)
	}
}

func (p *Player) SetMana(mana int) {
        if mana <= 12 && mana >= 0 {
                p.Mana = mana
        }
}

func (p *Player) SetLife(life int) {
		p.Life = life
		p.CheckRune()
}

func (p *Player) StackDraw() {
	fmt.Println("[GAME][RUNE] Player", p.Id, "stacking draw card")
	p.stack_draw++
}

func (p *Player) DrawStackCards() (err error) {
	max := p.stack_draw
	for i := 0 ; i < max ; i++ {
		err = p.DrawCard()
		if err != nil && p.Life <= 0 {
			return err
		}
		p.stack_draw--
	}
	return nil

}
func (p *Player) CheckRune() {
	if p.Life < STARTING_LIFE && p.Life <= (STARTING_RUNES * p.Runes) {		
		for ; p.Life <= (STARTING_RUNES * p.Runes) && p.Runes > 0 ; {
			p.Runes -= 1
			fmt.Println("[GAME][RUNE] Losing a rune. There are", p.Runes, "left")
			p.StackDraw()
		}
	}
}

func (p *Player) HandGetCard(id int) *Card {
        for _, c := range(p.Hand) {
                if c.Id == id {
                        return c
                }
        }
        return nil
}
func (p* Player) AddCardGrayeyard(c *Card) {
	p.Graveyard = append(p.Graveyard, c)
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
		c := p.Hand[idx]
        p.Hand[idx] = p.Hand[l - 1]
        if l >= 2 {
            p.Hand = p.Hand[:l - 1]
        } else if l == 1 {
            p.Hand = make([]*Card, 0)
		}
		p.AddCardGrayeyard(c)
        return nil
}
func (p *Player) HandPlayCard(id int) (error) {
        c := p.HandGetCard(id)
        if c == nil {
                return fmt.Errorf("[PLAYER] Card with id %d is not present in hand of player %d", id, p.Id)
		}
		
		if c.Cost > p.current_mana {
			return fmt.Errorf("[PLAYER] No more mana (%d) for playing card %d with cost %d", p.current_mana, id, c.Cost)
		}

		p.current_mana -= c.Cost

        switch c.Type {
        case CARD_TYPE_CREATURE:
                p.Board = append(p.Board, c)
		case CARD_TYPE_ITEM_BLUE:
		case CARD_TYPE_ITEM_GREEN:
		case CARD_TYPE_ITEM_RED:
                  
        default:
                return fmt.Errorf("Unkow card type %d for player %d", c.Type, p.Id)
        }

		if c.CardDraw > 0 {
			p.stack_draw += c.CardDraw
		}
		p.HandRemoveCard(id)
        return nil
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
		c := p.Board[idx]
        if l >= 2 {
            p.Board[idx] = p.Board[l - 1]
            p.Board = p.Board[:l - 1]
        } else if l == 1 {
            p.Board = make([]*Card, 0)
        }
		p.AddCardGrayeyard(c)
        fmt.Println("[GAME][DAMAGE] Monster", id, "has been killed")
        return nil

}

func (p *Player) GainLife(life int) {
	if life > 0 {
		fmt.Println("[GAME][HEALTH] Player", p.Id, "Gain", life, "life")
		p.SetLife(p.Life + life)
	} else if life < 0 {
		p.ReceiveDamage(-life)
	}
}

func (p *Player) ReceiveDamage(damage int) {
	if damage > 0 {
		fmt.Println("[GAME][HEALTH] Player", p.Id, "Receive", damage, "damage")
		p.SetLife(p.Life - damage)
	} else if damage < 0 {
		p.GainLife(-damage)
	}
}

func (p *Player) IncreaseMana() {
	if p.Mana < MAX_MANA {
		p.Mana = p.Mana + 1
		//fmt.Println("Increase Mana")
	}
}

func (g *Game) AddPlayer(p *Player) error {
	if len( g.players) >= MAX_PLAYERS {
		return errors.New("AddPlayer: There is already 2 players")
	}

	 g.players = append( g.players, p)
	return nil
}

func (g *Game) GetPlayerRandom() (*Player, error) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	i := random.Intn(len( g.players))
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

	 g.players[0], g.players[1] = g.players[index], g.players[len( g.players) - index - 1]

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

func (g *Game) PrintHand(p *Player) {
	if p != nil {
		for _, c := range(p.Hand) {
			fmt.Println("[GAME] Hand", c.Cost, c)
		}
	}
}
func (g *Game) PrintBoard(p *Player) {
	if p != nil {
		for _, c := range(p.Board) {
			fmt.Println("[GAME] Board", c.Attack, c.Defense, c)
		}
	}
}
func (g *Game) PrintDeck(p *Player) {
	if p != nil {
		for _, c := range(p.Deck.Cards) {
			fmt.Println("[GAME](", p.Id, ") Deck", c)
		}
	}
}

func (g *Game) PrintHeroDeck() { g.PrintDeck(g.Hero()) }
func (g *Game) PrintVilainDeck() { g.PrintDeck(g.Vilain()) }
func (g *Game) PrintHeroHand() { g.PrintHand(g.Hero()) }
func (g *Game) PrintHeroBoard() { g.PrintBoard(g.Hero()) }
func (g *Game) PrintVilainHand() { g.PrintHand(g.Vilain()) }
func (g *Game) PrintVilainBoard() { g.PrintBoard(g.Vilain()) }

func (g *Game) CheckRune() {

}
func (g *Game) CreatureFight(id1, id2 int) (error) {
	var dead1, dead2 bool

	dead1, dead2 = false, false

	c1 := g.Hero().BoardGetCard(id1)
	c2 := g.Vilain().BoardGetCard(id2)

	dmg1, err1 := g.DealDamageMonster(c2, c1.Attack)
	if err1 != nil { return err1 }
	if dmg1 > 0 && c1.Abilities[CARD_ABILITY_DRAIN] != "-" { g.Hero().GainLife(dmg1) }
	if dmg1 > 0 && c1.Abilities[CARD_ABILITY_LETHAL] != "-" { dead1 = true }
	dead1 = dead1 || (c2.Defense <= 0)

	dmg2, err2 := g.DealDamageMonster(c1, c2.Attack)
	if err2 != nil { return err2 }
	if dmg2 > 0 && c2.Abilities[CARD_ABILITY_LETHAL] != "-" { dead2 = true }
	dead2 = dead2 || c1.Defense <= 0

	if dead1 { g.Vilain().BoardRemoveCard(id2) }
	if dead2 { g.Hero().BoardRemoveCard(id1) }

	return nil
}
func (g *Game) DealDamageMonster(m1 *Card, dmg int) (int, error) {
	if m1 == nil || m1.Type != CARD_TYPE_CREATURE {
		return 0, fmt.Errorf("[GAME][ERROR] Card", m1, "is not a creature")
	}

	if dmg < 0 {
		return 0, fmt.Errorf("[GAME][ERROR] Cannot deal negative damage")
	}
	if dmg > 0 && m1.Abilities[CARD_ABILITY_WARD] != "-" {
		fmt.Println("[GAME][DEFENSE] Creature", m1.Id, "Use Ward protection")
		m1.Abilities[CARD_ABILITY_WARD] = "-"
		return 0, nil
	}
	m1.Defense -= dmg

	return dmg, nil
}
func (g *Game) DealDamage(c1, c2 *Card) (dead bool, err error) {
	var c1a, dmg, id2 int

	dead = false
	if c1 == nil {
		return false, fmt.Errorf("[GAME][ERROR] Card doesn't exist in Hand")	
	}

	switch c1.Type {
	case CARD_TYPE_CREATURE:
		if c1.Attack < 0 {
			c1a = -(c1.Attack)
		} else {
			c1a = c1.Attack
		}
	default:
		if c1.Defense < 0 {
			c1a = -(c1.Defense)
		} else {
			c1a = c1.Defense
		}
	}
	
	switch c2 {
	case nil:
		g.Vilain().ReceiveDamage(c1a)
		dmg 	= c1a
		id2 	= -1
	default:
		ward 		:= c2.Abilities[CARD_ABILITY_WARD]
		dmg, err 	= g.DealDamageMonster(c2, c1a)

		if c2.Defense < 0 && c1.Abilities[CARD_ABILITY_BREAKTHROUGH] != "-" {
			g.Hero().ReceiveDamage(-(c2.Defense))
		}

		if c2.Defense <= 0 ||
		   c1.Abilities[CARD_ABILITY_LETHAL] != "-" && ward == "-" {
			dead = true
		}
		id2 	= c2.Id
	} 
	if dmg > 0 {
		fmt.Println("[GAME][DAMAGE] Card", c1.Id, "deal", dmg, "to", id2)
	}
	return dead, nil
}
func (g *Game) ParseAction(actions string) (n int, err error) {
	var err1 error
	data := strings.Split(actions, ";")
	
	for _, a := range data {
		switch s := strings.Split(a, " "); s[0] {
		case "PICK":
			return g.ParseMovePick(s)
		case "ATTACK":
			err1 = g.ParseMoveAttack(s)
		case "SUMMON":
			err1 = g.ParseMoveSummon(s)
		case "USE":
			err1 = g.ParseMoveUse(s)
		case "PASS":
		default:
			fmt.Println("Unknow action", s[0])
		}
		if err1 != nil {
			fmt.Println(err1)
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
func (g *Game) ParseMoveUse(params []string) error {
	if len(params) != 3 {
		return errors.New("ParseUse: Format should be ACTION id1 id2")
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
	case "USE":
		return g.MoveUse(int(id1), int(id2))
	default:
		err := fmt.Sprintf("ParseUse: Unknow command %s", params[0])
		return errors.New(err)
	}
}
func (g *Game) CreatureBoostAttack(c *Card, bonus int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][ERROR] Can't boost attack")
	}
	if bonus == 0 {
		return nil
	}
	c.Attack += bonus
	fmt.Println("[GAME][CREATURE] Boost Attack by", bonus, "for", c.Id)
	return nil
}
func (g *Game) CreatureBoostDefense(c *Card, bonus int) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][ERROR] Can't boost defense")
	}
	if bonus == 0 {
		return nil
	}
	c.Defense += bonus
	fmt.Println("[GAME][CREATURE] Boost Defense by", bonus, "for", c.Id)
	return nil
}
func (g *Game) CreatureBoostAbilities(c *Card, bonus []string) (error) {
	if c == nil || c.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][ERROR] Can't boost abilities")
	}
	if len(bonus) != len(c.Abilities) {
		return fmt.Errorf("[GAME][ERROR] Wrong number of abilities")
	}
	for i := 0 ; i < len(bonus) ; i++ {
		if bonus[i] != "-" {
			c.Abilities[i] = bonus[i]
			fmt.Println("[GAME][CREATURE] Boost abilities ", bonus[i], "for", c.Id)
		}
	}
	return nil
}
func (g *Game) MoveUse(id1, id2 int) error {
	c1 := g.Hero().HandGetCard(id1)
	if c1 == nil {
		return fmt.Errorf("GAME: Use %d. Card doesn't exist in Hand", id1)
	}

	err := g.Hero().HandPlayCard(int(id1))
	if err != nil {
		return err
	}

	fmt.Println("[GAME][USE] Player", g.Hero().Id, "use card", id1, "on", id2, "for cost", c1.Cost, "(", g.Hero().current_mana, ")")

	switch c1.Type {
	case CARD_TYPE_ITEM_BLUE:
		g.Hero().GainLife(c1.HealthChange)
		g.Vilain().ReceiveDamage(-c1.OpponentHealthChange)
		if id2 != -1 {
			g.DealDamage(c1, nil)
		}

	case CARD_TYPE_ITEM_GREEN:
		c2 := g.Hero().BoardGetCard(id2)
		if c2 == nil {
			////////fmt.Fprintln(og.Stderr, "[GAME][ERROR] Player", g.Hero().Id, "doesn't have card", id2, "on board")
			return fmt.Errorf("[GAME][ERROR]: Use %d. Card doesn't exist in Board on Player %d", id2, g.Hero().Id)	
		}
		g.CreatureBoostAttack(c2, c1.Attack)
		g.CreatureBoostDefense(c2, c1.Defense)
		g.CreatureBoostAbilities(c2, c1.Abilities)
		g.Hero().GainLife(c1.HealthChange)
		g.Vilain().ReceiveDamage(-c1.OpponentHealthChange)

	case CARD_TYPE_ITEM_RED:
		c2 := g.Vilain().BoardGetCard(id2)
		if c2 == nil {
			return fmt.Errorf("[GAME][ERROR]: Use %d. Card doesn't exist in Board on Player %d", id2, g.Vilain().Id)	
		}
		g.CreatureBoostAttack(c2, c1.Attack)
		g.CreatureBoostDefense(c2, c1.Defense)
		g.Hero().GainLife(c1.HealthChange)
		g.Vilain().ReceiveDamage(-c1.OpponentHealthChange)

		for i := 0 ; i < len(c1.Abilities) ; i++ {
			if c2.Abilities[i] != "-" {
				c1.Abilities[i] = "-"
			}
		}
		if c2.Defense <= 0 { g.Vilain().BoardRemoveCard(id2) }
	}	
	return nil
}
func (g *Game) MoveSummon(id1 int) error {
	c1 := g.Hero().HandGetCard(id1)
	if c1 == nil {
		return fmt.Errorf("GAME: Use %d. Card doesn't exist in Hand", id1)
		
	}
	err := g.Hero().HandPlayCard(id1)
	if err != nil {
		return err
	}
	if c1.Type != CARD_TYPE_CREATURE {
		return fmt.Errorf("[GAME][SUMMON]: Can't summon card type %d", c1.Type)
	}

	fmt.Println("[GAME][SUMMON] Player", g.Hero().Id, "summon creature", id1, "for cost", c1.Cost, "(", g.Hero().current_mana, ")")

	g.Hero().GainLife(c1.HealthChange)
	g.Vilain().ReceiveDamage(-c1.OpponentHealthChange)
	

	return nil
}
func (g *Game) MoveAttackPolicy(id1, id2 int) (error) {
	guards := g.Vilain().BoardGetGuardsId()
	//len_guards := len(guards)

	exist, _ := in_array(id2, guards)
	if len(guards) > 0 && ! exist {
		return fmt.Errorf("[GAME][ATTACK] Move ATTACK %d %d not permitted", id1, id2)
	} 
	return nil
}
func (g *Game) MoveAttack(id1, id2 int) (err error) {
	var err_str string
	c1 := g.Hero().BoardGetCard(id1)
	if c1 == nil {
		fmt.Println("[GAME][ATTACK] Create", id1, "not present in Hero board")
	g.PrintHeroBoard()
		err_str = fmt.Sprintf("MoveAttack: Current player %d don't have card %d", g.Hero().Id, id1)
	g.PrintBoard(g.Hero())
		return errors.New(err_str)
	}

	err = g.MoveAttackPolicy(id1, id2)
	if err != nil {
		return err
	}

	fmt.Println("[GAME][ATTACK] Player", g.Hero().Id, "attack", id2, "with", id1)
	c1a := c1.Attack

	if id2 == -1 {
		g.Vilain().ReceiveDamage(c1a)
	} else {
		c2 := g.Vilain().BoardGetCard(id2)
		if c2 == nil {
			err_str = fmt.Sprintf("MoveAttack: Current oppoent %d don't have card %d", g.Vilain().Id, id2)
			return errors.New(err_str)
		}
		fmt.Println("[GAME][FIGHT]", c1.Id, "attack", c2.Id, ". May the force be with them")
		err = g.CreatureFight(c1.Id, c2.Id)
	}
	if c1.Abilities[CARD_ABILITY_DRAIN] != "-" {
		g.Hero().GainLife(c1.Attack)
	}

	return nil
}
func (g *Game) RawPlayers() [][]interface{} {

	raw_data := make([][]interface{}, 0)
	for i := 0 ; i < len( g.players) ; i++ {
		p := g.players[i].Raw()
		raw_data = append(raw_data, p)
	}
	return raw_data
}
func (g *Game) RawCards() [][]interface{} {
	cards := make([][]interface{}, 0)
	for _, c := range(g.Hero().Hand) {
		c.Location = 0
		cards = append(cards, c.Raw())
	}
	for _, c := range(g.Hero().Board) {
		c.Location = 1
		cards = append(cards, c.Raw())
	}
	for _, c := range(g.Vilain().Board) {
		c.Location = -1
		cards = append(cards, c.Raw())
	}
	//fmt.Println("Card for player", g.Hero().Hand, g.Hero().Id, cards)
	return cards
}
func (g *Game) Draft() (error) {

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	
	for i := 0; i < DECK_CARDS; i++ {
		draft_raw 	:= make([][]interface{}, DRAFT_PICK)
		draft 		:= make([]*Card, DRAFT_PICK)
		numbers := make([]int, 0)
		for j := 0; j < DRAFT_PICK; j++ {

			num := random.Intn(CARDS_COUNT)
			exist := true
			for ; exist ; exist, _ = in_array(num, numbers) {
				num = random.Intn(CARDS_COUNT)
			}
			new_card		:= CARDS[num].Copy()
			draft_raw[j] 	= new_card.Raw()
			draft[j] 		= new_card
			numbers 		= append(numbers, num)
		}

		fmt.Println("DRAFT(", i + 1, "):")
		for _, k := range(draft_raw) {
			fmt.Println("\t", k)
		}
		for h := 0; h < MAX_PLAYERS ; h++ {
			pick, err := g.Hero().IA.Move(g.RawPlayers(), draft_raw, len(g.Vilain().Hand), MOVE_TIMEOUT)
			if err != nil {
				return err
			}
			params := strings.Split(pick, " ")
			num, err := g.ParseAction(pick)
			if err != nil {
				return fmt.Errorf("Wrong pick at draft %d", num)
			}
			
			if len(params) == 1 && params[0] == "PASS" {
				num = 0
			} else if num < 0 {
				return fmt.Errorf("Wrong Action %s", pick)
			}


			new_id := g.Hero().Deck.Count() + g.Vilain().Deck.Count() + 1
			fmt.Println("Set ID", new_id, "for", draft[num].CardNumber)
			draft[num].Id = new_id
			g.Hero().Deck.AddCard(draft[num].Copy())
			g.NextPlayer()	

			
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
func (g *Game) Start() (winner *Player, err error) {
	var round int
    

	

	winner = nil
	//fmt.Println("Starting Game")
	if len( g.players) != MIN_PLAYERS {
		return nil, errors.New("AddPlayer: There should be players")
	}

	hero, _ := g.GetPlayerRandom()
	//fmt.Println("Starting Player", start_player)
	g.OrderPlayer(hero)

	err = g.Draft()
	if err != nil {
		fmt.Println("Nil")
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}


	g.PrintHeroDeck()
	g.PrintVilainDeck()

	g.Hero().Deck.Shuffle()
	g.Vilain().Deck.Shuffle()

	g.Hero().DrawCardN(4)
	g.Vilain().DrawCardN(5)


	winner 	= nil
	round 	= 2

	timeout := 1000
	for winner == nil {
		g.Hero().IncreaseMana()
		g.Hero().ReloadMana()

		fmt.Println("================================")
		fmt.Println("[GAME] Turn:", round / 2,
					"Hero:", g.Hero().Id,
					"Life:", g.Hero().Life,
					"Mana:", g.Hero().Mana,
					"Deck:", g.Hero().Deck.Count(),
					"Runes:", g.Hero().Runes)
		fmt.Println("[GAME] Turn:", round / 2,
					"Vilain:", g.Vilain().Id,
					"Life:", g.Vilain().Life,
					"Mana:", g.Vilain().Mana,
					"Deck:", g.Vilain().Deck.Count(),
					"Runes:", g.Vilain().Runes)

		g.PrintHeroHand()
		g.PrintHeroBoard()
		err = g.Hero().DrawStackCards()

		if err != nil {
			winner = g.Vilain()
			break
		}

		err = g.Hero().DrawCard()
		if err != nil && g.Hero().Life <= 0 {
			winner = g.Vilain()
			break
		}


		if (round / 2) > 1 {
			timeout = MOVE_TIMEOUT
		}

		move, err := g.Hero().IA.Move(g.RawPlayers(), g.RawCards(), len(g.Vilain().Hand), timeout)
		if err != nil {
			return g.Vilain(), err
		}

		_, err = g.ParseAction(move)
		if err != nil {
			return g.Vilain(), err
		}

		winner = g.CheckWinner()
		if winner != nil {
			break
		} 

		//time.Sleep(1 * time.Second)
		g.NextPlayer()
		round = round + 1
	}

	return winner, nil
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "ia-binary-1 ia-binary-2")
		os.Exit(1)
	}
	p1 := NewPlayer(1, STARTING_LIFE, STARTING_MANA, STARTING_RUNES, NewIAPlayer(1, os.Args[1]))
	p2 := NewPlayer(2, STARTING_LIFE, STARTING_MANA, STARTING_RUNES, NewIAPlayer(2, os.Args[2]))
	gm := NewGame()

	gm.AddPlayer(p1)
	gm.AddPlayer(p2)

	start := time.Now()
	winner, err := gm.Start()
	elapsed := time.Since(start)
	fmt.Println("[GAME] Game duration:", elapsed)
	if err != nil {
		fmt.Println(err)
	} 
	if winner != nil {
		fmt.Println("[GAME] The Winner is Player", winner.Id)
	}
		
}
