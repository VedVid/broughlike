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
	blt "bearlibterminal"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var KeyboardLayout int
var CustomControls bool

func main() {
	var cells = new(Board)
	var actors = new(Creatures)
	StartGame(cells, actors)
	for {
		RenderAll(*cells, *actors)
		if (*actors)[0].HPCurrent <= 0 {
			DeleteSaves()
			blt.Read()
			break
		}
		key := ReadInput()
		if (key == blt.TK_S && blt.Check(blt.TK_SHIFT) != 0) ||
			key == blt.TK_CLOSE {
			err := SaveGame(*cells, *actors)
			if err != nil {
				fmt.Println(err)
			}
			break
		} else if key == blt.TK_Q && blt.Check(blt.TK_SHIFT) != 0 {
			DeleteSaves()
			break
		} else {
			turnSpent := Controls(key, (*actors)[0], cells, actors)
			if turnSpent == true {
				CreaturesTakeTurn(*cells, *actors)
			}
		}
	}
	blt.Close()
}

func NewGame(b *Board, c *Creatures) {
	/* Function NewGame initializes game state - creates player, monsters,
	   and game map. */
	*b = InitializeEmptyMap()
	player, err := NewPlayer(1, 1)
	if err != nil {
		fmt.Println(err)
	}
	*c = append(*c, player)
}

func StartGame(b *Board, c *Creatures) {
	/* Function StartGame determines if game save is present (and valid), then
	   loads data, or initializes new game.
	   Panics if some-but-not-all save files are missing. */
	_, errBoard := os.Stat(MapPathGob)
	_, errCreatures := os.Stat(CreaturesPathGob)
	if errBoard == nil && errCreatures == nil {
		LoadGame(b, c)
	} else if errBoard != nil && errCreatures != nil {
		NewGame(b, c)
	} else {
		txt := CorruptedSaveError(errBoard, errCreatures)
		fmt.Println("Error: save files are corrupted: " + txt)
		panic(-1)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	InitializeBLT()
	InitializeKeyboardLayouts()
	ReadOptionsControls()
	ChooseKeyboardLayout()
}
