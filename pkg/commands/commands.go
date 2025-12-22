package commands

import (
	"errors"
	"fmt"
	"os"
)

type Register struct {
	CommandsList map[string]Command
}

type Command struct {
	name        string
	Callback    Callback
	description string
}
type Callback func() error

func GetRegister() *Register {
	return &Register{
		map[string]Command{
			"exit": {
				name:        "exit",
				description: "Exits the CLI",
				Callback:    cmdExit,
			},
		},
	}
}

func (r Register) GetCommand(name string) (Command, error) {
	if command, ok := r.CommandsList[name]; ok {
		return command, nil
	} else {
		return Command{}, errors.New("invalid command name...")
	}
}

func cmdExit() error {
	fmt.Println("Exiting tasks tool... Stay productive!")
	os.Exit(0)
	return nil
}
