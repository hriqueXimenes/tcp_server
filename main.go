package main

import (
	"os"

	"github.com/hriqueXimenes/sumo_logic_server/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
