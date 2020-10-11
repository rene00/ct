package main

import (
	"ct/cmd"
)

const version = "0.0.4"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
