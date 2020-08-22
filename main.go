package main

import (
	"ct/cmd"
	"fmt"
)

const version = "0.0.1"

func Version() {
	fmt.Println(version)
}

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
