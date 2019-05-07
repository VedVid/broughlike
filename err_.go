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
	"strconv"
	"unicode/utf8"
)

func LayerError(layer int) string {
	/* Function LayerError is helper function that returns string
	   to error; it takes layer integer as argument and returns string. */
	return "\n    <layer:  " + strconv.Itoa(layer) + ">"
}

func LayerWarning(layerMonster, layerDefault int) string {
	/* Function LayerWarning is helper function that returns string to error.
	   It is called when, during creation, monster layer is not equal
	   to CreaturesLayer constant defined at the top of render.go file. */
	txt := "\n    <monster layer: " + strconv.Itoa(layerMonster) +
		";    default layer: " + strconv.Itoa(layerDefault) + ">"
	return txt
}

func CoordsError(x, y int) string {
	/* Function CoordsError is helper function that returns string
	   to error; it takes coords x, y as arguments and returns string,
	   with use global MapSizeX and MapSizeY constants. */
	txt := "\n    <x: " + strconv.Itoa(x) + "; y: " + strconv.Itoa(y) +
		"; map width: " + strconv.Itoa(MapSizeX) + "; map height: " +
		strconv.Itoa(MapSizeY) + ">"
	return txt
}

func CharacterLengthError(character string) string {
	/* Function CharacterLengthError is helper function that returns string
	   to error; it takes character string as argument and returns string.
	   Character (as something's representation on map) is supposed to be
	   one-letter long. */
	txt := "\n    <length: " + strconv.Itoa(utf8.RuneCountInString(character)) +
		"; character: " + character + ">"
	return txt
}

func PlayerAIError(ai int) string {
	/* Function PlayerAIError is helper function that returns string to error;
	   it takes ai code (integer) as argument and returns string.
	   Player AI is supposed to be PlayerAI (defined in ai.go).
	   It's supposed to be warning, not error. */
	txt := "\n    <player ai code: " + strconv.Itoa(ai) + ">"
	return txt
}

func InitialHPError(hp int) string {
	/* Function InitialHPError is helper function that returns string to error;
	   it takes creature's HPMax as argument and returns string.
	   It will be warning instead of error sometimes - negative hp for newly created
	   creatures is unusual, but it is not bug per se. */
	txt := "\n    <fighter hp: " + strconv.Itoa(hp) + ">"
	return txt
}

func InitialAttackError(attack int) string {
	/* Function InitialAttackError is helper function that returns string
	   to error; it takes creature's attack value as argument and returns string.
	   Attack value should not be negative. */
	txt := "\n    <fighter attack: " + strconv.Itoa(attack) + ">"
	return txt
}

func InitialDefenseError(defense int) string {
	/* Function InitialDefenseError is helper function that returns string
	   to error; it takes creature's defense value as argument.
	   Defense value should not be negative. */
	txt := "\n    <fighter defense: " + strconv.Itoa(defense) + ">"
	return txt
}

func VectorCoordinatesOutOfMapBounds(startX, startY, targetX, targetY int) string {
	/* Function VectorCoordinatesOutOfMapBounds is helper function that returns
	   string to error; it takes vector source and vector target coords as arguments.
	   It is called if source or target is out of map bounds. */
	sx, sy := strconv.Itoa(startX), strconv.Itoa(startY)
	tx, ty := strconv.Itoa(targetX), strconv.Itoa(targetY)
	txt := "\n    <MapSizeX: 0.." + strconv.Itoa(MapSizeX-1) + "; MapSizeY: 0.." +
		strconv.Itoa(MapSizeY-1) + ";" +
		"\n    VectorStartPoint:  " + sx + ", " + sy + "; " +
		"\n    VectorTargetPoint: " + tx + ", " + ty + ">"
	return txt
}

func CorruptedSaveError(errBoard, errCreatures error) string {
	/* Function CorruptedSaveError is helper function that returns string to error.
	   It takes three specific errors as arguments (only one of them has to be != nil).
	   It is called when game can not find all two save files in directory. */
	errorBoard, errorCreatures := "", ""
	if errBoard != nil {
		errorBoard = "map.gob "
	}
	if errCreatures != nil {
		errorCreatures = "monsters.gob "
	}
	txt := "\n    <Following files are missing: " + errorBoard + errorCreatures + ">"
	return txt
}
