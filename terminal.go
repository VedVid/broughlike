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
	"runtime"
	"strconv"

	blt "bearlibterminal"
)

const (
	// Setting BearLibTerminal window.
	WindowSizeX = 14
	WindowSizeY = 14
	MapSizeX    = 12
	MapSizeY    = 12
	UIPosX      = 0
	UIPosY      = MapSizeY
	UISizeX     = WindowSizeX
	UISizeY     = WindowSizeY - MapSizeY
	GameTitle   = "Broughlike"
	GameVersion = "0.1"
	FontName    = "Deferral-Square.ttf"
	FontSize    = 24
)

func constrainThreads() {
	/* Constraining processor and threads is necessary,
	   because BearLibTerminal often crashes otherwise. */
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func InitializeBLT() {
	/* Constraining threads and setting BearLibTerminal window. */
	constrainThreads()
	blt.Open()
	sizeX, sizeY := strconv.Itoa(WindowSizeX), strconv.Itoa(WindowSizeY)
	sizeFont := strconv.Itoa(FontSize)
	window := "window: size=" + sizeX + "x" + sizeY
	blt.Set(window + ", title=' " + GameTitle + " " + GameVersion +
		"'; font: " + FontName + ", size=" + sizeFont)
	blt.Clear()
	blt.Refresh()
}
