package main

import (
	"fmt"
	"os"
)

const (
	UpArrow    = "\x1b[A"
	DownArrow  = "\x1b[B"
	RightArrow = "\x1b[C"
	LeftArrow  = "\x1b[D"
	ClearLine  = "\x1b[2K\r"
	Prompt     = "gator "
)

type replState struct {
	history           []string
	historyIndex      int
	currentInput      []byte
	currentInputIndex int
	yanked            []byte
}

func startRepl() {
	state := replState{
		historyIndex:      -1,
		currentInputIndex: -1,
	}

	buf := make([]byte, 512)

	fmt.Print(Prompt)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("error reading from Stdin: %v\r\n", err)
			return
		}

		inputStr := string(buf[:n])
		switch inputStr {
		case UpArrow:
			state.upInput()
		case DownArrow:
			state.downInput()
		case LeftArrow:
			state.leftInput()
		case RightArrow:
			state.rightInput()
		case "\f":
			fmt.Printf("\033[H\033[2J%s", Prompt+string(state.currentInput))
		case "\r", "\n":
			/*
				if state.enterInput(cfg) {
					return
				}
			*/
		case "\x7f", "\x08":
			state.backInput()
		case "\x01":
			for state.currentInputIndex != -1 {
				state.leftInput()
			}
		case "\x03":
			fmt.Print("^C\r\n")
			return
		case "\x05":
			for state.currentInputIndex < len(state.currentInput)-1 {
				state.rightInput()
			}
		case "\x15":
			state.yanked = state.ctrlU()
		case "\x0B":
			state.yanked = state.ctrlK()
		case "\x19":
			state.ctrlY()
		default:
			state.keyInput(n, buf)
		}
	}
}
