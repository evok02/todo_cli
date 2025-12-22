package input

import "strings"

type Args struct {
	Arguments []string
}

func Normalize(input string) []string {
	return cleanInput(splitInput(input))
}

func cleanInput(args []string) []string {
	var cleanArgs []string
	for _, v := range args {
		cleanArgs = append(cleanArgs, strings.ToLower(v))
	}
	return cleanArgs
}

func splitInput(input string) []string {
	return strings.Fields(input)
}
