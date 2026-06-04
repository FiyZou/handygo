package main

import (
	"fmt"
	"os"

	"github.com/FiyZou/handygo/tools"
)

func main() {
	if err := tools.NewCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
