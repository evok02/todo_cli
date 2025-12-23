package input

import (
	"fmt"
	"github.com/google/shlex"
	"strings"
)

type Args struct {
	Arguments []string
}

func Normalize(input string) ([]string, error) {
	out, err := splitInput(input)
	if err != nil {
		return nil, err
	}
	return cleanInput(out), nil
}

func cleanInput(args []string) []string {
	var cleanArgs []string
	for _, v := range args {
		cleanArgs = append(cleanArgs, strings.ToLower(v))
	}
	return cleanArgs
}

func splitInput(input string) ([]string, error) {
	out, err := shlex.Split(input)
	if err != nil {
		return nil, fmt.Errorf("couldn't split users input: %w", err)
	}
	return out, nil
}
