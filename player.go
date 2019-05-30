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
	"strconv"
	"unicode/utf8"
)

const AmmoMax = 5

func NewPlayer(x, y int) (*Creature, error) {
	/* NewPlayer is function that returns new Creature
	   (that is supposed to be player) from json file passed as argument.
	   It replaced old code that was encouraging hardcoding data in go files.
	   Errors returned by json package are not very helpful, and
	   hard to work with, so there is lazy panic for them. */
	const playerPath = "./data/player/player.json"
	var player = &Creature{}
	err := CreatureFromJson(playerPath, player)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	player.X, player.Y = x, y
	var err2 error
	if player.Layer < 0 {
		txt := LayerError(player.Layer)
		err2 = errors.New("Creature layer is smaller than 0." + txt)
	}
	if player.Layer != PlayerLayer {
		txt := LayerWarning(player.Layer, PlayerLayer)
		err2 = errors.New("Creature layer is not equal to CreaturesLayer constant." + txt)
	}
	if player.X < 0 || player.X >= MapSizeX || player.Y < 0 || player.Y >= MapSizeY {
		txt := CoordsError(player.X, player.Y)
		err2 = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(player.Char) != 1 {
		txt := CharacterLengthError(player.Char)
		err2 = errors.New("Creature character string length is not equal to 1." + txt)
	}
	if player.AIType != PlayerAI {
		txt := PlayerAIError(player.AIType)
		err2 = errors.New("Warning: Player AI is supposed to be " +
			strconv.Itoa(PlayerAI) + "." + txt)
	}
	if player.HPMax < 0 {
		txt := InitialHPError(player.HPMax)
		err2 = errors.New("Creature HPMax is smaller than 0." + txt)
	}
	if player.Attack < 0 {
		txt := InitialAttackError(player.Attack)
		err2 = errors.New("Creature attack value is smaller than 0." + txt)
	}
	if player.Defense < 0 {
		txt := InitialDefenseError(player.Defense)
		err2 = errors.New("Creature defense value is smaller than 0." + txt)
	}
	return player, err2
}

func (c *Creature) MoveOrAttack(tx, ty int, b Board, all Creatures) bool {
	/* Method MoveOrAttack decides if Creature will move or attack other Creature;
	   It has *Creature receiver, and takes tx, ty (coords) integers as arguments,
	   and map of current level, and list of all Creatures.
	   Starts by target that is nil, then iterates through Creatures. If there is
	   Creature on targeted tile, that Creature becomes new target for attack.
	   Otherwise, Creature moves to specified Tile.
	   It's supposed to take player as receiver (attack / moving enemies is
	   handled differently - check ai.go and combat.go). */
	var target *Creature
	turnSpent := false
	for i, _ := range all {
		if all[i].X == c.X+tx && all[i].Y == c.Y+ty {
			if all[i].HPCurrent > 0 {
				target = all[i]
				break
			}
		}
	}
	if target != nil {
		c.AttackTarget(target)
		turnSpent = true
	} else {
		turnSpent = c.Move(tx, ty, b)
	}
	return turnSpent
}

func (c *Creature) SetWeapon(i int) bool {
	i--
	if i < 0 || i > 3 {
		return false
	}
	if i == c.Active {
		return false
	} else {
		c.Active = i
		return true
	}
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
