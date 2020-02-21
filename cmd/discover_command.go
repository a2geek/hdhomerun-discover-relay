package cmd

import (
	"fmt"
)

type DiscoverCommand struct {
}

func (cmd DiscoverCommand) Execute(args []string) error {
	fmt.Println("TBD!")
	return fmt.Errorf("Not implemented")
}
