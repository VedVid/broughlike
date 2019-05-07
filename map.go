/*
Copyright (c) 2018, Tomasz "VedVid" Nowakowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"unicode/utf8"
)

type Tile struct {
	// Tiles are map cells - floors, walls, doors.
	BasicProperties
	VisibilityProperties
	Explored bool
	CollisionProperties
}

/* Board is map representation, that uses 2d slice
   to hold data of its every cell. */
type Board [][]*Tile

func NewTile(layer, x, y int, character, name, color, colorDark string,
	alwaysVisible, explored, blocked, blocksSight bool) (*Tile, error) {
	/* Function NewTile takes all values necessary by its struct,
	   and creates then returns Tile. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Tile layer is smaller than 0." + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Tile coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Tile character string length is not equal to 1." + txt)
	}
	tileBasicProperties := BasicProperties{x, y, character, name, color,
		colorDark}
	tileVisibilityProperties := VisibilityProperties{layer, alwaysVisible}
	tileCollisionProperties := CollisionProperties{blocked, blocksSight}
	tileNew := &Tile{tileBasicProperties, tileVisibilityProperties,
		explored, tileCollisionProperties}
	return tileNew, err
}

func InitializeEmptyMap() Board {
	/* Function InitializeEmptyMap returns new Board, filled with
	   generic (ie "empty") tiles.
	   It starts by declaring 2d slice of *Tile - unfortunately, Go seems to
	   lack simple way to do it, therefore it's necessary to use
	   the first for loop.
	   The second, nested loop initializes specific Tiles within Board bounds.
	   In this game, all map is explored from the start. */
	b := make([][]*Tile, MapSizeX)
	for i := range b {
		b[i] = make([]*Tile, MapSizeY)
	}
	for x := 0; x < MapSizeX; x++ {
		for y := 0; y < MapSizeY; y++ {
			var err error
			b[x][y], err = NewTile(BoardLayer, x, y, "#", "floor", "dark gray",
				"darkest gray", true, true, true, false)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return b
}

func MakeDrunkardsMap(b Board) {
	percent := float64(MapSizeX*MapSizeY) / float64(100)
	digMin := RoundFloatToInt(percent * float64(60))
	digMax := RoundFloatToInt(percent * float64(85))
	fmt.Println(digMin, digMax)
	diggedPercent := RandRange(digMin, digMax)
	var directions = [][]int{{0, 1}, {-1, 0}, {1, 0}, {0, -1}}
	x, y := MapSizeX/2, MapSizeY/2
	for {
		if b[x][y].Blocked == true {
			b[x][y].Char = "."
			b[x][y].Blocked = false
			b[x][y].Color = "light gray"
			b[x][y].ColorDark = "dark gray"
			diggedPercent--
		}
		if diggedPercent <= 0 {
			break
		}
		dir := directions[rand.Intn(len(directions))]
		newX := x + dir[0]
		newY := y + dir[1]
		if newX >= 0 && newX < MapSizeX &&
			newY >= 0 && newY < MapSizeY {
			x = newX
			y = newY
		}
	}
}

func MapCheck(b Board) bool {
	valid := true
	var q1 = []int{0, MapSizeX / 2, 0, MapSizeY / 2}
	var q2 = []int{MapSizeX / 2, MapSizeX, 0, MapSizeY / 2}
	var q3 = []int{MapSizeX / 2, MapSizeX, MapSizeY / 2, MapSizeY}
	var q4 = []int{0, MapSizeX / 2, MapSizeY / 2, MapSizeY}
	var count = []int{}
	sum := 0
	var qs = [][]int{q1, q2, q3, q4}
	for _, q := range qs {
		counter := 0
		for x := q[0]; x < q[1]; x++ {
			for y := q[2]; y < q[3]; y++ {
				if b[x][y].Blocked == true {
					counter++
				}
			}
		}
		count = append(count, counter)
		sum += counter
	}
	average := float64(sum) / float64(len(qs))
	percent := average / float64(100)
	min := RoundFloatToInt(average - (percent * 20))
	max := RoundFloatToInt(average + (percent * 20))
	for _, v := range count {
		if v < min || v > max {
			valid = false
		}
	}
	return valid
}

func MakeNewLevel() Board {
	var b Board
	for {
		b = InitializeEmptyMap()
		MakeDrunkardsMap(b)
		if MapCheck(b) == true {
			break
		}
	}
	return b
}
