package main

import (
	"os"

	"github.com/MelonSmasher/browserbeat/cmd"

	_ "github.com/MelonSmasher/browserbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
