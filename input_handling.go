package main

import (
	"fmt"
	"slices"
)

func (r *replState) keyInput(n int, buf []byte) {
	if n == 1 && buf[0] >= 32 && buf[0] <= 126 {
		if len(r.currentInput) == 0 {
			r.currentInput = append(r.currentInput, buf[0])
		} else {
			r.currentInput = slices.Insert(r.currentInput, r.currentInputIndex+1, buf[0])
		}
		fmt.Print(ClearLine + Prompt + string(r.currentInput))
		r.currentInputIndex++
		if r.currentInputIndex != len(r.currentInput)-1 {
			fmt.Printf("\x1b[%dD", (len(r.currentInput)-1)-r.currentInputIndex)
		}
	}
}

func (r *replState) backInput() {
	if len(r.currentInput) > 0 && r.currentInputIndex != -1 {

		r.currentInput = slices.Delete(r.currentInput, r.currentInputIndex, r.currentInputIndex+1)

		fmt.Print(ClearLine + Prompt + string(r.currentInput))
		r.currentInputIndex--
		if r.currentInputIndex != len(r.currentInput)-1 {
			fmt.Printf("\x1b[%dD", (len(r.currentInput)-1)-r.currentInputIndex)
		}
	}
}

func (r *replState) upInput() {
	if len(r.history) > 0 {
		if r.historyIndex == -1 {
			r.historyIndex = len(r.history) - 1
		} else if r.historyIndex != 0 {
			r.historyIndex--
		}
		lastCmd := r.history[r.historyIndex]
		r.currentInput = []byte(lastCmd)
		r.currentInputIndex = len(r.currentInput) - 1

		fmt.Print(ClearLine + Prompt + lastCmd)
	}
}

func (r *replState) downInput() {
	if r.historyIndex != -1 {
		if r.historyIndex != len(r.history)-1 {
			r.historyIndex++
		} else {
			r.historyIndex = -1
		}

		var nextCmd string
		if r.historyIndex != -1 {
			nextCmd = r.history[r.historyIndex]
		} else {
			nextCmd = ""
		}
		r.currentInput = []byte(nextCmd)
		r.currentInputIndex = len(r.currentInput) - 1

		fmt.Print(ClearLine + Prompt + nextCmd)
	}
}

func (r *replState) leftInput() {
	if r.currentInputIndex >= 0 {
		fmt.Print(LeftArrow)
		r.currentInputIndex--
	}
}

func (r *replState) rightInput() {
	if r.currentInputIndex < len(r.currentInput)-1 {
		fmt.Print(RightArrow)
		r.currentInputIndex++
	}
}

/*
	func (r *replState) enterInput(cfg *config) bool {
		dirtyInput := string(r.currentInput)
		if dirtyInput != "" {
			fmt.Print("\r\n")
			input := cleanInput(dirtyInput)
			if len(input) == 0 {
				fmt.Print(Prompt)
				return false
			}

			command := input[0]
			var args []string
			if len(input) > 1 {
				args = input[1:]
			}

			if commandStruct, ok := commandRegistry[command]; ok {
				err := commandStruct.callback(cfg, args...)
				if err != nil {
					fmt.Printf("%v\r\n", err)
				}
			} else {
				fmt.Print("Unknown command\r\n")
			}

			if len(r.history) == 0 || r.history[len(r.history)-1] != dirtyInput {
				r.history = append(r.history, dirtyInput)
			}

			if command == "exit" {
				return true
			}

			r.currentInput = nil
			r.currentInputIndex = -1
			r.historyIndex = -1
			fmt.Print(Prompt)
		}
		return false
	}
*/
func (r *replState) ctrlU() []byte {
	deleted := r.currentInput[:r.currentInputIndex+1]
	r.currentInput = r.currentInput[r.currentInputIndex+1:]
	r.currentInputIndex = len(r.currentInput) - 1
	fmt.Print(ClearLine + Prompt + string(r.currentInput))
	for r.currentInputIndex != -1 {
		r.leftInput()
	}
	return deleted
}

func (r *replState) ctrlK() []byte {
	deleted := r.currentInput[r.currentInputIndex+1:]
	r.currentInput = r.currentInput[:r.currentInputIndex+1]
	fmt.Print(ClearLine + Prompt + string(r.currentInput))
	return deleted
}

func (r *replState) ctrlY() {
	moveLeftAmount := (len(r.currentInput) - 1) - r.currentInputIndex

	r.currentInput = slices.Insert(r.currentInput, r.currentInputIndex+1, r.yanked...)
	r.currentInputIndex = len(r.currentInput) - 1
	fmt.Print(ClearLine + Prompt + string(r.currentInput))

	for range moveLeftAmount {
		r.leftInput()
	}
}
