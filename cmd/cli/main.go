package main

import (
	"bufio"
	"fmt"
	"github.com/evok02/todo_cli/pkg/commands"
	"github.com/evok02/todo_cli/pkg/input"
	"github.com/evok02/todo_cli/storage/sqlite"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	reg := commands.GetRegister()
	var a input.Args
	sqlite.Init()

	for {
		fmt.Print("\ntasks > ")
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		} else {
			a.Arguments = input.Normalize(scanner.Text())
			command, err := reg.GetCommand(a.Arguments[0])

			if err != nil {
				fmt.Println(err)
			} else {
				command.Callback()
			}

		}
	}
}
