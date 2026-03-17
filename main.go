package main

import (
	"log"

	"0ximg.sh/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
