package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type Register struct {
	CommandsList map[string]Command
}

type Command struct {
	name        string
	Callback    Callback
	description string
}

type Callback func(c *Config, args ...string) error

func GetRegister() *Register {
	return &Register{
		map[string]Command{
			"exit": {
				name:        "exit",
				description: "Exits the CLI",
				Callback:    cmdExit,
			},
			"add": {
				name:        "add",
				description: "Adds new task to the storage",
				Callback:    cmdAdd,
			},
			"update": {
				name:        "update",
				description: "Updates description of the task with the given ID",
				Callback:    cmdUpdate,
			},
			"delete": {
				name:        "delete",
				description: "Soft deletes task from the storage",
				Callback:    cmdSoftDelete,
			},
			"hard-delete": {
				name:        "hard-delete",
				description: "Hard deletes task from the storage",
				Callback:    cmdHardDelete,
			},
			"mark-in-progress": {
				name:        "mark-in-progress",
				description: "Changes status to 'in progress'",
				Callback:    cmdMarkInProgress,
			},
			"mark-done": {
				name:        "mark-in-progress",
				description: "Changes status to 'in progress'",
				Callback:    cmdMarkDone,
			},
			"list": {
				name:        "list",
				description: "Prints all registered tasks",
				Callback:    cmdList,
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

func cmdExit(c *Config, args ...string) error {
	fmt.Println("Exiting tasks tool... Stay productive!")
	os.Exit(0)
	return nil
}

func cmdAdd(c *Config, args ...string) error {
	_, err := c.Repo.CreateTask(args[0])
	if err != nil {
		return err
	}
	return nil
}

func cmdUpdate(c *Config, args ...string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("could't convert id arg: %w", err)
	}

	if err = c.Repo.UpdateTask(id, args[1]); err != nil {
		return err
	}

	return nil
}

func cmdSoftDelete(c *Config, args ...string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("couldn't convert id arg: %w", err)
	}

	if err = c.Repo.SoftDeleteTask(id); err != nil {
		return err
	}

	return nil
}

func cmdHardDelete(c *Config, args ...string) error {
	id, err := strconv.Atoi(args[0])

	if err != nil {
		return fmt.Errorf("couldn't convert arg to int: %w", err)
	}

	_, err = c.Repo.HardDeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func cmdMarkInProgress(c *Config, args ...string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("could't convert id arg to int: %w", err)
	}

	err = c.Repo.MarkInProgress(id)

	if err != nil {
		return err
	}

	return nil
}

func cmdMarkDone(c *Config, args ...string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("could't convert id arg to int: %w", err)
	}

	err = c.Repo.MarkDone(id)

	if err != nil {
		return err
	}

	return nil
}

func cmdList(c *Config, args ...string) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 0, 0, '\t', 0)
	var m string
	if len(args) > 0 {
		m = args[0]
	}
	tasks, err := c.Repo.List(m)
	if err != nil {
		return err
	}

	fmt.Println("ID\tDescription\tStatus\tCreatedAt")
	for _, t := range tasks {
		fmt.Printf("%d\t%s\t%s\t%v\n", t.ID, t.Description, t.Status, t.CreatedAt)
	}
	w.Flush()

	return nil
}
