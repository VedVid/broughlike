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

const (
	// Special characters.
	CorpseChar = "%"
)

type Creature struct {
	/* Creatures are living objects that
	   moves, attacks, dies, etc. */
	BasicProperties
	VisibilityProperties
	CollisionProperties
	FighterProperties
}

// Creatures holds every creature on map.
type Creatures []*Creature

func NewCreature(x, y int, monsterFile string) (*Creature, error) {
	/* NewCreature is function that returns new Creature from
	   json file passed as argument. It replaced old code that
	   was encouraging hardcoding data in go files.
	   Errors returned by json package are not very helpful, and
	   hard to work with, so there is lazy panic for them. */
	var monster = &Creature{}
	err := CreatureFromJson(CreaturesPathJson+monsterFile, monster)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	monster.X, monster.Y = x, y
	var err2 error
	if monster.Layer < 0 {
		txt := LayerError(monster.Layer)
		err2 = errors.New("Creature layer is smaller than 0." + txt)
	}
	if monster.Layer != CreaturesLayer {
		txt := LayerWarning(monster.Layer, CreaturesLayer)
		err2 = errors.New("Creature layer is not equal to CreaturesLayer constant." + txt)
	}
	if monster.X < 0 || monster.X >= MapSizeX || monster.Y < 0 || monster.Y >= MapSizeY {
		txt := CoordsError(monster.X, monster.Y)
		err2 = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(monster.Char) != 1 {
		txt := CharacterLengthError(monster.Char)
		err2 = errors.New("Creature character string length is not equal to 1." + txt)
	}
	if monster.HPMax < 0 {
		txt := InitialHPError(monster.HPMax)
		err2 = errors.New("Creature HPMax is smaller than 0." + txt)
	}
	if monster.Attack < 0 {
		txt := InitialAttackError(monster.Attack)
		err2 = errors.New("Creature attack value is smaller than 0." + txt)
	}
	if monster.Defense < 0 {
		txt := InitialDefenseError(monster.Defense)
		err2 = errors.New("Creature defense value is smaller than 0." + txt)
	}
	var monsterColors = []string{BallisticColorGood, KineticColorGood,
		ElectromagneticColorGood, ExplosiveColorGood}
	monster.Color = monsterColors[rand.Intn(len(monsterColors))]
	monster.ColorDark = monster.Color
	switch monster.Color {
	case BallisticColorGood:
		monster.Ballistic = 1
	case KineticColorGood:
		monster.Kinetic = 1
	case ElectromagneticColorGood:
		monster.Electromagnetic = 1
	case ExplosiveColorGood:
		monster.Explosive = 1
	}
	return monster, err2
}

func (c *Creature) Move(tx, ty int, b Board) bool {
	/* Move is method of Creature; it takes target x, y as arguments;
	   check if next move won't put Creature off the screen, then updates
	   Creature coords. */
	turnSpent := false
	newX, newY := c.X+tx, c.Y+ty
	if newX >= 0 &&
		newX <= MapSizeX-1 &&
		newY >= 0 &&
		newY <= MapSizeY-1 {
		if b[newX][newY].Blocked == false {
			c.X = newX
			c.Y = newY
			if b[newX][newY].Stairs == false {
				turnSpent = true
			} else {
				if c.AIType == PlayerAI {
					if CurrentLevel < len(LevelMaps) {
						CurrentLevel++
					} else {
						GameWon = true
					}
				}
			}
		}
	}
	return turnSpent
}

func (c *Creature) PickUp(b Board) bool {
	/* PickUp is method that has *Creature as receiver.
	   It will use *Tile as argument.
	   The idea is to check, if tile has deposits of mana first,
	   then allow player to "charge" energy from this deposit. */
	turnSpent := false
	t := b[c.X][c.Y]
	if t.Drained == true {
		return turnSpent
	}
	if (t.Resources == BallisticResource && c.Ballistic < AmmoMax) ||
		(t.Resources == ExplosiveResource && c.Explosive < AmmoMax) ||
		(t.Resources == KineticResource && c.Kinetic < AmmoMax) ||
		(t.Resources == ElectromagneticResource &&
			c.Electromagnetic < AmmoMax) {
		c.AddAmmo(t.Resources)
	} else {
		return turnSpent
	}
	t.Drained = true
	t.Color = ResourcesColors[t.Resources][1]
	turnSpent = true
	return turnSpent
}

func (c *Creature) AddAmmo(resource int) {
	switch resource {
	case BallisticResource:
		c.Ballistic += RandRange(1, 3)
		if c.Ballistic > AmmoMax {
			c.Ballistic = AmmoMax
		}
	case ExplosiveResource:
		c.Explosive += RandRange(1, 3)
		if c.Explosive > AmmoMax {
			c.Explosive = AmmoMax
		}
	case KineticResource:
		c.Kinetic += RandRange(1, 3)
		if c.Kinetic > AmmoMax {
			c.Kinetic = AmmoMax
		}
	case ElectromagneticResource:
		c.Electromagnetic += RandRange(1, 3)
		if c.Electromagnetic > AmmoMax {
			c.Electromagnetic = AmmoMax
		}
	}
}

func (c *Creature) Die() {
	/* Method Die is called when Creature's HP drops below zero.
	   Die() has *Creature as receiver.
	   Receiver properties changes to fit better to corpse. */
	c.Layer = DeadLayer
	c.Name = "corpse of " + c.Name
	c.Blocked = false
	c.BlocksSight = false
	c.AIType = NoAI
}

func FindMonsterByXY(x, y int, c Creatures) *Creature {
	/* Function FindMonsterByXY takes desired coords and list
	   of all available creatures. It iterates through this list,
	   and returns nil or creature that occupies specified coords. */
	var monster *Creature
	for i := 0; i < len(c); i++ {
		if x == c[i].X && y == c[i].Y {
			monster = c[i]
			break
		}
	}
	return monster
}
