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
	"fmt"
)

func (c *Creature) AttackTarget(t *Creature) {
	/* Receiver "c" is attacker, argument "t" is target. */
	t.TakeDamage(c.Attack - t.Defense)
}

func (c *Creature) Shoot(dx, dy int, b Board, cs Creatures) bool {
	turnSpent := false
	tx, ty := c.X, c.Y
	if dx == (-1) {
		tx = 0
	} else if dx == 1 {
		tx = MapSizeX - 1
	} else if dy == (-1) {
		ty = 0
	} else if dy == 1 {
		ty = MapSizeY - 1
	}
	vec, err := NewVector(c.X, c.Y, tx, ty)
	if err != nil {
		fmt.Println(err)
	}
	_ = ComputeVector(vec)
	_, _, target := ValidateVector(vec, b, cs)
	var attacks = []string{
		"ballistic", "explosive", "kinetic", "electromagnetic"}
	activeAttack := attacks[c.Active]
	switch activeAttack {
	case "ballistic":
		if c.Ballistic <= 0 {
			return turnSpent
		} else {
			c.Ballistic--
		}
	case "kinetic":
		if c.Kinetic <= 0 {
			return turnSpent
		} else {
			c.Kinetic--
		}
	case "electromagnetic":
		if c.Electromagnetic <= 0 {
			return turnSpent
		} else {
			c.Electromagnetic--
		}
	case "explosive":
		if c.Explosive <= 0 {
			return turnSpent
		} else {
			c.Explosive--
		}
	}
	turnSpent = true
	if target != nil {
		if (activeAttack == "ballistic" && target.Ballistic > 0) ||
			(activeAttack == "kinetic" && target.Kinetic > 0) ||
			(activeAttack == "electromagnetic" && target.Electromagnetic > 0) ||
			(activeAttack == "explosive" && target.Explosive > 0) {
			target.TakeDamage((c.Attack - target.Defense) * 2)
		}
	}
	return turnSpent
}

func (c *Creature) TakeDamage(dmg int) {
	/* Method TakeDamage has *Creature as receiver and takes damage integer
	   as argument. dmg value is deducted from Creature current HP.
	   If HPCurrent is below zero after taking damage, Creature dies. */
	c.HPCurrent -= dmg
	if c.HPCurrent <= 0 {
		c.Die()
	}
}
