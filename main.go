package main

import (
	"hse-cat/internal/di"
	"os"
)

func main() {
	c := di.NewContainer(os.Stdin, os.Stdout)
	c.App().Run()
}
