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
	STARTING_MANA	= 0
	STARTING_LIFE	= 30
	STARTING_CARD	= 30
	DECK_CARDS		= 30

	DRAFT_PICK		= 3

    CARD_TYPE_MONSTER           = 0
	CARD_TYPE_ITEMDEAL			= 1
	CARD_TYPE_ITEMGAIN			= 2
	CARD_TYPE_ITEMGIVE			= 3
	CARD_TYPE_ITEMREMOVE		= 4
	CARD_TYPE_CREATURE			= 0
	
	CARD_COLOR_BLUE				= 0
	CARD_COLOR_GREEN			= 1
	CARD_COLOR_RED				= 2
	CARD_COLOR_NONE				= 3

    CARD_ABILITY_BREAKTHROUGH   = 0
    CARD_ABILITY_GUARD          = 3
	CARD_ABILITY_CHARGE         = 1
	

)

var CARDS = []*Card{
	NewCard(1, 1, 2, 1, "------", 1, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(2, 1, 1, 2, "------", 0, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(3, 1, 2, 2, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(4, 2, 1, 5, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(5, 2, 4, 1, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(6, 2, 3, 2, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(7, 2, 2, 2, "-----W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(8, 2, 2, 3, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(9, 3, 3, 4, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(10, 3, 3, 1, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(11, 3, 5, 2, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(12, 3, 2, 5, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(13, 4, 5, 3, "------", 1, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(14, 4, 9, 1, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(15, 4, 4, 5, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(16, 4, 6, 2, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(17, 4, 4, 5, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(18, 4, 7, 4, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(19, 5, 5, 6, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(20, 5, 8, 2, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(21, 5, 6, 5, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(22, 6, 7, 5, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(23, 7, 8, 8, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(24, 1, 1, 1, "------", 0, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(25, 2, 3, 1, "------", -2, -2, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(26, 2, 3, 2, "------", 0, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(27, 2, 2, 2, "------", 2, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(28, 2, 1, 2, "------", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(29, 2, 2, 1, "------", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(30, 3, 4, 2, "------", 0, -2, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(31, 3, 3, 1, "------", 0, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(32, 3, 3, 2, "------", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(33, 4, 4, 3, "------", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(34, 5, 3, 5, "------", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(35, 6, 5, 2, "B-----", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(36, 6, 4, 4, "------", 0, 0, 2, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(37, 6, 5, 7, "------", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(38, 1, 1, 3, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(39, 1, 2, 1, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(40, 3, 2, 3, "--DG--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(41, 3, 2, 2, "-CD---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(42, 4, 4, 2, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(43, 6, 5, 5, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(44, 6, 3, 7, "--D-L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(45, 6, 6, 5, "B-D---", -3, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(46, 9, 7, 7, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(47, 2, 1, 5, "--D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(48, 1, 1, 1, "----L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(49, 2, 1, 2, "---GL-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(50, 3, 3, 2, "----L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(51, 4, 3, 5, "----L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(52, 4, 2, 4, "----L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(53, 4, 1, 1, "-C--L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(54, 3, 2, 2, "----L-", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(55, 2, 0, 5, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(56, 4, 2, 7, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(57, 4, 1, 8, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(58, 6, 5, 6, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(59, 7, 7, 7, "------", 1, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(60, 7, 4, 8, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(61, 9, 10, 10, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(62, 12, 12, 12, "B--G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(63, 2, 0, 4, "---G-W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(64, 2, 1, 1, "---G-W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(65, 2, 2, 2, "-----W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(66, 5, 5, 1, "-----W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(67, 6, 5, 5, "-----W", 0, -2, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(68, 6, 7, 5, "-----W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(69, 3, 4, 4, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(70, 4, 6, 3, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(71, 4, 3, 2, "BC----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(72, 4, 5, 3, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(73, 4, 4, 4, "B-----", 4, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(74, 5, 5, 4, "B--G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(75, 5, 6, 5, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(76, 6, 5, 5, "B-D---", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(77, 7, 7, 7, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(78, 8, 5, 5, "B-----", 0, -5, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(79, 8, 8, 8, "B-----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(80, 8, 8, 8, "B--G--", 0, 0, 1, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(81, 9, 6, 6, "BC----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(82, 7, 5, 5, "B-D--W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(83, 0, 1, 1, "-C----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(84, 2, 1, 1, "-CD--W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(85, 3, 2, 3, "-C----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(86, 3, 1, 5, "-C----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(87, 4, 2, 5, "-C-G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(88, 5, 4, 4, "-C----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(89, 5, 4, 1, "-C----", 2, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(90, 8, 5, 5, "-C----", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(91, 0, 1, 2, "---G--", 0, 1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(92, 1, 0, 1, "---G--", 2, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(93, 1, 2, 1, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(94, 2, 1, 4, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(95, 2, 2, 3, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(96, 2, 3, 2, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(97, 3, 3, 3, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(98, 3, 2, 4, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(99, 3, 2, 5, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(100, 3, 1, 6, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(101, 4, 3, 4, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(102, 4, 3, 3, "---G--", 0, -1, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(103, 4, 3, 6, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(104, 4, 4, 4, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(105, 5, 4, 6, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(106, 5, 5, 5, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(107, 5, 3, 3, "---G--", 3, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(108, 5, 2, 6, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(109, 5, 5, 6, "------", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(110, 5, 0, 9, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(111, 6, 6, 6, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(112, 6, 4, 7, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(113, 6, 2, 4, "---G--", 4, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(114, 7, 7, 7, "---G--", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(115, 8, 5, 5, "---G-W", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(116, 12, 8, 8, "BCDGLW", 0, 0, 0, CARD_COLOR_NONE, CARD_TYPE_CREATURE),
	NewCard(117, 1, 1, 1, "B-----", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(118, 0, 0, 3, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(119, 1, 1, 2, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(120, 2, 1, 0, "----L-", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(121, 2, 0, 3, "------", 0, 0, 1, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(122, 2, 1, 3, "---G--", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(123, 2, 4, 0, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(124, 3, 2, 1, "--D---", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(125, 3, 1, 4, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(126, 3, 2, 3, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(127, 3, 0, 6, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(128, 4, 4, 3, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(129, 4, 2, 5, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(130, 4, 0, 6, "------", 4, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(131, 4, 4, 1, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(132, 5, 3, 3, "B-----", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(133, 5, 4, 0, "-----W", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(134, 4, 2, 2, "------", 0, 0, 1, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(135, 6, 5, 5, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(136, 0, 1, 1, "------", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(137, 2, 0, 0, "-----W", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(138, 2, 0, 0, "---G--", 0, 0, 1, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(139, 4, 0, 0, "----LW", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(140, 2, 0, 0, "-C----", 0, 0, 0, CARD_COLOR_GREEN, CARD_TYPE_ITEMGIVE),
	NewCard(141, 0, -1, -1, "------", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMGIVE),
	NewCard(142, 0, 0, 0, "BCDGLW", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMREMOVE),
	NewCard(143, 0, 0, 0, "---G--", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMREMOVE),
	NewCard(144, 1, 0, -2, "------", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMDEAL),
	NewCard(145, 3, -2, -2, "------", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMGIVE),
	NewCard(146, 4, -2, -2, "------", 0, -2, 0, CARD_COLOR_RED, CARD_TYPE_ITEMGIVE),
	NewCard(147, 2, 0, -1, "------", 0, 0, 1, CARD_COLOR_RED, CARD_TYPE_ITEMDEAL),
	NewCard(148, 2, 0, -2, "BCDGLW", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMREMOVE),
	NewCard(149, 3, 0, 0, "BCDGLW", 0, 0, 1, CARD_COLOR_RED, CARD_TYPE_ITEMREMOVE),
	NewCard(150, 2, 0, -3, "------", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMDEAL),
	NewCard(151, 5, 0, -99, "BCDGLW", 0, 0, 0, CARD_COLOR_RED, CARD_TYPE_ITEMREMOVE),
	NewCard(152, 7, 0, -7, "------", 0, 0, 1, CARD_COLOR_RED, CARD_TYPE_ITEMDEAL),
	NewCard(153, 2, 0, 0, "------", 5, 0, 0, CARD_COLOR_BLUE, CARD_TYPE_ITEMGAIN),
	NewCard(154, 2, 0, 0, "------", 0, -2, 1, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	NewCard(155, 3, 0, -3, "------", 0, -1, 0, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	NewCard(156, 3, 0, 0, "------", 3, -3, 0, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	NewCard(157, 3, 0, -1, "------", 1, 0, 1, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	NewCard(158, 3, 0, -4, "------", 0, 0, 0, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	NewCard(159, 4, 0, -3, "------", 3, 0, 0, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	NewCard(160, 2, 0, 0, "------", 2, -2, 0, CARD_COLOR_BLUE, CARD_TYPE_ITEMDEAL),
	
	/*
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
*/
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
	Color			int
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
		fmt.Println("STDOUT: (", ia.id, ")", scanner.Text())
		done_reading <- true
	} ()

	readerErr := bufio.NewReader(ia.stderr)
	scannerErr := bufio.NewScanner(readerErr)

	go func() {
		for scannerErr.Scan() {
		  	fmt.Println("STDERR: (", ia.id, ")", scannerErr.Text())
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
	fmt.Println("STDIN: (", ia.id, ")", data)
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
	    	cost int,
        	attack int,
			defense int,
			abilities string,
	     	heroHealthChange int,
	     	opponentHealthChange int,
			cardDraw int,
			color int,
			type_ int, 
) *Card {
        return &Card{
				CardNumber: 	cardNumber,
                Id:             -1,
				Location:		0,
                Type:           type_,
                Cost:           cost,
                Attack:         attack,
                Defense:        defense,
				Abilities:      strings.Split(abilities, ""),
				Color:			color,
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
		fmt.Println("GAME: Player", p.Id, "draw card", c)
		p.Draw(c)
    } else {
		fmt.Println("GAME:", err)
	}
}

func (p *Player) DrawCardN(n int) {
    for i := 0 ; i < n ; i++ {
		p.DrawCard()
    }
}

func (p *Player) Draw(c *Card) {
    if p.HandGetCard(c.Id) == nil {
        p.Hand = append(p.Hand, c)
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
	fmt.Println("Card for player", g.Hero().Hand, g.Hero().Id, cards)
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
			num := random.Intn(len(CARDS))
			for exist, _ := in_array(numbers, num); exist; {
				num = random.Intn(len(CARDS))
			}
			c 				:= CARDS[num]
			draft_raw[j] 	= c.Raw()
			draft[j] 		= c
			numbers 		= append(numbers, num)
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
			
			if len(params) == 1 && params[0] == "PASS" {
				num = 0
			} else if num < 0 {
				return fmt.Errorf("Wrong Action %s", pick)
			}

			draft[num].Id = g.Hero().Deck.Count() + g.Vilain().Deck.Count() + 1

			g.Hero().Deck.AddCard(draft[num])
	
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

func (g *Game) CheckDraw() bool {
	return false
}

func (g *Game) Start() (winner *Player, err error) {
	var round int

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
		fmt.Println("Nil")
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}


	g.Hero().Deck.Shuffle()
	g.Vilain().Deck.Shuffle()

	g.Hero().DrawCardN(4)
	g.Vilain().DrawCardN(5)


	winner 	= nil
	round 	= 2

	timeout := 1000
	for winner == nil {
		fmt.Println("Turn:", round / 2)
		
		g.Hero().IncreaseMana()
		g.Hero().DrawCard()

		if (round / 2) > 1 {
			timeout = 100
		}
		move, err := g.Hero().IA.Move(g.RawPlayers(), g.RawCards(), g.Vilain().Deck.Count(), timeout)
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
		} else if g.CheckDraw() {
			break
		}

		
		time.Sleep(1 * time.Second)
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
	p1 := NewPlayer(1, STARTING_LIFE, STARTING_MANA, NewIAPlayer(1, os.Args[1]))
	p2 := NewPlayer(2, STARTING_LIFE, STARTING_MANA, NewIAPlayer(2, os.Args[2]))
	gm := NewGame()

	gm.AddPlayer(p1)
	gm.AddPlayer(p2)

	gm.Start()
}
