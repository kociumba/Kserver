package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Command struct {
	Name        string
	NameANSI    string
	Description string
	Callback    func()
}

var Commands = map[string]Command{
	"help": {
		Name:        "help",
		NameANSI:    "\033[33mhelp\033[0m",
		Description: "show this message",
		Callback:    nil,
	},
	"exit": {
		Name:        "exit",
		NameANSI:    "\033[33mexit\033[0m",
		Description: "exit the server",
		Callback:    func() { os.Exit(0) },
	},
	"config": {
		Name:        "config",
		NameANSI:    "\033[33mconfig\033[0m",
		Description: "show the config",
		Callback:    showConfig,
	},
}

var (
	title = "\n\033[1mKserver help:\033[0m\n\n"

	HelpMsg = title + createHelpMessage(Commands)
)

func calculatePadding(commands map[string]Command) int {
	maxLength := 0
	for _, command := range commands {
		length := len(command.Name)
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

func createHelpMessage(commands map[string]Command) string {
	var builder strings.Builder
	padding := calculatePadding(commands) + 3 // Add extra spaces for separation

	for _, command := range commands {
		// Calculate the number of spaces needed for padding
		spaces := strings.Repeat(" ", padding-len(command.Name))
		builder.WriteString(fmt.Sprintf("%s%s- %s\n", command.NameANSI, spaces, command.Description))
	}

	return builder.String()
}

func showConfig() {
	f, err := os.Open("kserver.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		panic(err)
	}
}
