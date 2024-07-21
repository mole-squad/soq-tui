package main

import (
	"fmt"
	"os"

	"github.com/mole-squad/soq-tui/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	os.Exit(0)
}
