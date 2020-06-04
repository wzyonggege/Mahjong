package internal

import (
	"bufio"
	"fmt"
	"github.com/wzyonggege/goutils/convert"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

/**
* @Author: Jam Wong
* @Date: 2020/6/3
 */

type Game struct {
	Tiles       chan *Tile
	Players     []*Player
	DealerIndex int
}

func NewGame() *Game {
	tiles := shuffle()

	g := &Game{
		Players:     make([]*Player, 0),
		Tiles:       make(chan *Tile, len(tiles)),
		DealerIndex: 0,
	}

	for _, i := range tiles {
		g.Tiles <- i
	}

	// TODO Dice 骰子确定index
	fmt.Println("init players")
	// init player
	for i := 0; i < 4; i++ {
		p := InitPlayer(i, g.DealerIndex == 0)
		g.Players = append(g.Players, p)
		fmt.Printf("player: %s\n", p.Name)
	}

	return g
}

// shuffle 洗牌
// TODO shuffle 洗牌 Fisher-Yates 高纳德置乱算法
// 从1~8中随机抽取一个数，例如随机数是3，那么交换第8位和第三位的数字。
// 此时数组顺序为12456783，重复第一步，从1~7中抽取一个数字，假设为4，那么交换第7位和第4位的数字
// 依次类推，直到第一个位置也被替代。
func shuffle() []*Tile {
	tiles := AllTiles()
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})

	return tiles
}

func (game *Game) DealNToPlayer(index int, n int) {
	pl := game.Players[index]
	if n < 1 || n > 4 {
		log.Fatal("n require [1, 4]")
	}
	for i := 1; i <= n; i++ {
		pl.Draw(game.Tiles)
	}
}

// Deal 发牌
func (game *Game) Deal() {
	fmt.Println("start deal")

	// 三轮：一人四张， 从index 开始
	for i := 0; i < 4; i++ {
		for _, j := range game.Players {
			if i == 3 {
				game.DealNToPlayer(j.Index, 1)
			} else {
				game.DealNToPlayer(j.Index, 4)
			}
			//fmt.Printf("%s: %s\n", j.Name, j.Show())
		}
	}

	// 庄 1 张
	game.DealNToPlayer(game.DealerIndex, 1)

	// debug
	for _, j := range game.Players {
		j.SortTiles()
		fmt.Printf("%s: %s\n\n", j.Name, j.Show())
	}

	i := 1
	for {
		if len(game.Players[game.DealerIndex].HoldTiles) != 14 {
			log.Fatal("not 14")
		}

		if game.Players[game.DealerIndex].Win() {
			fmt.Printf("winwin %s\n", game.Players[game.DealerIndex].Show())
			break
		}

		fmt.Printf("round %d: %s \n", i, game.Players[game.DealerIndex].Show())

		discardIndex := readIndex()
		_dis := game.Players[game.DealerIndex].HoldTiles[discardIndex]
		game.Players[game.DealerIndex].Discard(discardIndex)
		fmt.Printf("%s ====> %s\n", game.Players[game.DealerIndex].Show(), _dis.Print())
		i++

		fmt.Println("==========================================================")

		// draw
		draw := game.Players[game.DealerIndex].Draw(game.Tiles)

		fmt.Printf("<<<======= %s \n", draw.Print())

	}
}

func readIndex() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	i, _ := convert.StringToInt(text)
	return i - 1
}
