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
)

const (
	/* Actions' identifiers.
	   They have to be strings due to
	   testing these values with strings
	   from options_controls.cfg file. */
	StrMoveNorth   = "MOVE_NORTH"
	StrMoveWest    = "MOVE_WEST"
	StrMoveEast    = "MOVE_EAST"
	StrMoveSouth   = "MOVE_SOUTH"
	StrAttackNorth = "ATTACK_NORTH"
	StrAttackWest  = "ATTACK_WEST"
	StrAttackEast  = "ATTACK_EAST"
	StrAttackSouth = "ATTACK_SOUTH"
	StrPickup      = "PICKUP"
	StrSetWeapon1  = "CHOOSE_WEAPON_1"
	StrSetWeapon2  = "CHOOSE_WEAPON_2"
	StrSetWeapon3  = "CHOOSE_WEAPON_3"
	StrSetWeapon4  = "CHOOSE_WEAPON_4"
)

var Actions = []string{
	// List of all possible actions.
	StrMoveNorth,
	StrMoveWest,
	StrMoveEast,
	StrMoveSouth,
	StrAttackNorth,
	StrAttackWest,
	StrAttackEast,
	StrAttackSouth,
	StrPickup,
	StrSetWeapon1,
	StrSetWeapon2,
	StrSetWeapon3,
	StrSetWeapon4,
}

var CommandKeys = map[int]string{
	// Mapping keyboard scancodes to Action identifiers.
	blt.TK_W:     StrMoveNorth,
	blt.TK_D:     StrMoveEast,
	blt.TK_S:     StrMoveSouth,
	blt.TK_A:     StrMoveWest,
	blt.TK_UP:    StrAttackNorth,
	blt.TK_RIGHT: StrAttackEast,
	blt.TK_DOWN:  StrAttackSouth,
	blt.TK_LEFT:  StrAttackWest,
	blt.TK_SPACE: StrPickup,
	blt.TK_1:     StrSetWeapon1,
	blt.TK_2:     StrSetWeapon2,
	blt.TK_3:     StrSetWeapon3,
	blt.TK_4:     StrSetWeapon4,
}

/* Place to store customized controls scheme,
   in the same manner as CommandKeys. */
var CustomCommandKeys = map[int]string{}

func Command(com string, p *Creature, b *Board, c *Creatures) bool {
	/* Function Command handles input received from Controls.
	   Most important argument passed to Command is string "com" that
	   is action identifier (action identifiers are stored as constants
	   at the top of this file). It calls player methods regarding to
	   passed command.
	   Returns true if command is valid and takes turn.
	   Otherwise, return false. */
	turnSpent := false
	switch com {
	case StrMoveNorth:
		turnSpent = p.MoveOrAttack(0, -1, *b, *c)
	case StrMoveEast:
		turnSpent = p.MoveOrAttack(1, 0, *b, *c)
	case StrMoveSouth:
		turnSpent = p.MoveOrAttack(0, 1, *b, *c)
	case StrMoveWest:
		turnSpent = p.MoveOrAttack(-1, 0, *b, *c)
	case StrAttackNorth:
		turnSpent = p.Shoot(0, -1, *b, *c)
	case StrAttackEast:
		turnSpent = p.Shoot(1, 0, *b, *c)
	case StrAttackSouth:
		turnSpent = p.Shoot(0, 1, *b, *c)
	case StrAttackWest:
		turnSpent = p.Shoot(-1, 0, *b, *c)
	case StrPickup:
		turnSpent = p.PickUp(*b)
	case StrSetWeapon1:
		turnSpent = p.SetWeapon(1)
	case StrSetWeapon2:
		turnSpent = p.SetWeapon(2)
	case StrSetWeapon3:
		turnSpent = p.SetWeapon(3)
	case StrSetWeapon4:
		turnSpent = p.SetWeapon(4)
	}
	return turnSpent
}

func Controls(k int, p *Creature, b *Board, c *Creatures) bool {
	/* Function Controls takes integer 'k' (that is pressed key - blt uses
	   scancodes internally) and trying to find match key-command in
	   CommandKeys.
	   Value to return is determined in Command func. */
	turnSpent := false
	var command string
	if CustomControls == false {
		command = CommandKeys[k]
	} else {
		command = CustomCommandKeys[k]
	}
	turnSpent = Command(command, p, b, c)
	return turnSpent
}

func ReadInput() int {
	/* Function ReadInput is replacement of default blt's Read function that
	   returns QWERTY scancode. To provide (still experimental - I don't have
	   access to non-QWERTY keyboard physically) support for different
	   keyboard layouts, there are maps (in options.go) that matches
	   non-QWERTY input with QWERTY scancodes.
	   Some keys are hardcoded - like numpad, enter, etc. These hardcoded
	   keys are tested as first place as it's much cheaper operation than
	   checking map.
	   KeyMap content depends on chosen keyboard layout. */
	key := blt.Read()
	for _, v := range HardcodedKeys {
		if key == v {
			return v
		}
	}
	var r rune
	if blt.Check(blt.TK_WCHAR) != 0 {
		r = rune(blt.State(blt.TK_WCHAR))
	}
	return KeyMap[r]
}
