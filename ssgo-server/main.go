package main

import (
	"ssgo-server/cmd"

	_ "modernc.org/sqlite"
)

func main() {
	cmd.Execute()
}
