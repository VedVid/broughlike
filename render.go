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
	"unicode/utf8"
)

const (
	/* Constant values for layers. Their usage is optional,
	   but (for now, at leas) recommended, because default
	   rendering functions depends on these values.
	   They are important for proper clearing characters
	   that should not be displayed, as, for example,
	   bracelet under the monster. */
	_ = iota
	UILayer
	BoardLayer
	DeadLayer
	ObjectsLayer
	CreaturesLayer
	PlayerLayer
	LookLayer
)

const (
	BallisticIcon            = "☉"
	BallisticColorGood       = "crimson"
	BallisticColorBad        = "darker crimson"
	ExplosiveIcon            = "☄"
	ExplosiveColorGood       = "flame"
	ExplosiveColorBad        = "darker flame"
	KineticIcon              = "☀"
	KineticColorGood         = "amber"
	KineticColorBad          = "darker amber"
	ElectromagneticIcon      = "☇"
	ElectromagneticColorGood = "cyan"
	ElectromagneticColorBad  = "darker cyan"
)

func PrintBoard(b Board, c Creatures) {
	/* Function PrintBoard is used in RenderAll function.
	   Takes level map and list of monsters as arguments
	   and iterates through Board.
	   It has to check for "]" and "[" characters, because
	   BearLibTerminal uses these symbols for config.
	   Instead of checking it here, one could just remember to
	   always pass "]]" instead of "]".
	   Prints every tile on its coords if certain conditions are met:
	   is Explored already, and:
	   - is in player's field of view (prints "normal" color) or
	   - is AlwaysVisible (prints dark color). */
	for x := 0; x < MapSizeX; x++ {
		for y := 0; y < MapSizeY; y++ {
			// Technically, "t" is new variable with own memory address...
			t := b[x][y] // Should it be *b[x][y]?
			blt.Layer(t.Layer)
			if t.Explored == true {
				color := blt.ColorFromName(t.Color)
				SimplePutExt(t.X, t.Y, 0, 0, t.Char, color, color, color, color)
			}
		}
	}
}

func PrintCreatures(b Board, c Creatures) {
	/* Function PrintCreatures is used in RenderAll function.
	   Takes map of level and slice of Creatures as arguments.
	   Iterates through Creatures.
	   It has to check for "]" and "[" characters, because
	   BearLibTerminal uses these symbols for config.
	   Instead of checking it here, one could just remember to
	   always pass "]]" instead of "]".
	   Checks for every creature on its coords if certain conditions are met:
	   AlwaysVisible bool is set to true, or is in player fov. */
	for _, v := range c {
		blt.Layer(v.Layer)
		baseColor := blt.ColorFromName(v.Color)
		badColor := blt.ColorFromName("darkest gray")
		var colors = []uint32{baseColor, baseColor, baseColor, baseColor}
		hppc := Percents(v.HPCurrent, v.HPMax)
		switch {
		case hppc <= 0:
			colors = []uint32{badColor, badColor, badColor, badColor}
		case hppc < 25:
			colors[0], colors[1], colors[2] = badColor, badColor, badColor
		case hppc < 50:
			colors[0], colors[1] = badColor, badColor
		case hppc < 75:
			colors[0] = badColor
		default:
			colors = []uint32{baseColor, baseColor, baseColor, baseColor}
		}
		SimplePutExt(v.X, v.Y, 0, 0, v.Char,
			colors[0], colors[1], colors[2], colors[3])
	}
}

func PrintUI(c *Creature) {
	/* Function PrintUI takes *Creature (it's supposed to be player) as argument.
	   It prints UI infos on the right side of screen.
	   For now its functionality is very modest, but it will expand when
	   new elements of game mechanics will be introduced. So, for now, it
	   provides only one basic, yet essential information: player's HP. */
	blt.Layer(UILayer)
	const hpIconFull = "♦"
	const hpIconEmpty = "♢"
	hp := "[color=light blue]"
	for i := 1; i <= c.HPMax; i++ {
		if i > c.HPCurrent {
			hp = hp + hpIconEmpty
		} else {
			hp = hp + hpIconFull
		}
	}
	hp = hp + "[/color]"
	blt.Print(UIPosX+1, UIPosY, hp)
	const levelIcon = "■"
	const levelColor = "darkest green"
	const levelCurrentColor = "dark green"
	for i := 1; i <= NoOfLevels; i++ {
		levelStr := ""
		if i != CurrentLevel {
			levelStr =
				"[color=" + levelColor + "]" + levelIcon + "[/color]"
		} else {
			levelStr =
				"[color=" + levelCurrentColor + "]" + levelIcon + "[/color]"
		}
		blt.Print(UIPosX+i-1+3, UIPosY+1, levelStr)
	}
	for y := 0; y < AmmoMax; y++ {
		ballisticStr := ""
		if y < c.Ballistic {
			ballisticStr =
				"[color=" + BallisticColorGood + "]" + BallisticIcon + "[/color]"
		} else {
			ballisticStr =
				"[color=" + BallisticColorBad + "]" + BallisticIcon + "[/color]"
		}
		blt.Print(MapSizeX, 1+y, ballisticStr)
		explosiveStr := ""
		if y < c.Explosive {
			explosiveStr =
				"[color=" + ExplosiveColorGood + "]" + ExplosiveIcon + "[/color]"
		} else {
			explosiveStr =
				"[color=" + ExplosiveColorBad + "]" + ExplosiveIcon + "[/color]"
		}
		blt.Print(MapSizeX+1, 1+y, explosiveStr)
		kineticStr := ""
		if y < c.Kinetic {
			kineticStr =
				"[color=" + KineticColorGood + "]" + KineticIcon + "[/color]"
		} else {
			kineticStr =
				"[color=" + KineticColorBad + "]" + KineticIcon + "[/color]"
		}
		blt.Print(MapSizeX, MapSizeY-2-y, kineticStr)
		electromagneticStr := ""
		if y < c.Electromagnetic {
			electromagneticStr =
				"[color=" + ElectromagneticColorGood + "]" +
					ElectromagneticIcon + "[/color]"
		} else {
			electromagneticStr =
				"[color=" + ElectromagneticColorBad + "]" +
					ElectromagneticIcon + "[/color]"
		}
		blt.Print(MapSizeX+1, MapSizeY-2-y, electromagneticStr)
	}
	var numbersTemp = []string{"1", "2", "3", "4"}
	var numbers = []string{}
	for i, v := range numbersTemp {
		if i == c.Active {
			numbers = append(numbers, "[color=white]"+v+"[/color]")
		} else {
			numbers = append(numbers, "[color=gray]"+v+"[/color]")
		}
	}
	blt.Print(MapSizeX, 0, numbers[0]+numbers[1])
	blt.Print(MapSizeX, MapSizeY-1, numbers[2]+numbers[3])
}

func RenderAll(b Board, c Creatures) {
	/* Function RenderAll prints every tile and character on game screen.
	   Takes board slice (ie level map), slice of objects, and slice of creatures
	   as arguments.
	   At first, it clears whole terminal window, then uses arguments:
	   CastRays (for raycasting FOV) of first object (assuming that it is player),
	   then calls functions for printing map, objects and creatures.
	   At the end, RenderAll calls blt.Refresh() that makes
	   changes to the game window visible. */
	blt.Clear()
	PrintBoard(b, c)
	PrintCreatures(b, c)
	PrintUI((c)[0])
	blt.Refresh()
}

func WinScreen() {
	blt.Clear()
	txt := "You have won!"
	blt.Print((WindowSizeX-utf8.RuneCountInString(txt))/2,
		WindowSizeY/2, "You have won!")
	blt.Refresh()
	DeleteSaves()
	blt.Read()
	blt.Close()
}
