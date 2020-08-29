package main

import (
	"ct/cmd"
)

const version = "0.0.2"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
