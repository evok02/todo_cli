package main

import (
	"bufio"
	"fmt"
	"github.com/evok02/todo_cli/pkg/commands"
	"github.com/evok02/todo_cli/pkg/input"
	"github.com/evok02/todo_cli/pkg/repo"
	"github.com/evok02/todo_cli/storage/sqlite"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	reg := commands.GetRegister()
	var a input.Args

	db, err := sqlite.Init()
	if err != nil {
		fmt.Print(err)
	}

	c := commands.Config{
		Repo: repo.NewRepo(db),
	}

	for {
		fmt.Print("\ntasks > ")
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			continue
		} else {
			a.Arguments, err = input.Normalize(scanner.Text())
			if err != nil {
				fmt.Println(err)
				continue
			}

			command, err := reg.GetCommand(a.Arguments[0])
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				err := command.Callback(&c, a.Arguments[1:]...)
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	}
}
