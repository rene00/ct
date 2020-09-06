package main

import (
	"ct/cmd"
)

const version = "0.0.3"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
