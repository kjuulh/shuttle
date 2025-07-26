package main

import (
	"os"

	"github.com/kjuulh/shuttle/cmd"
)

func main() {
	cmd.Execute(os.Stdout, os.Stderr)
}
