package main

import (
	"fmt"
	"os"

	"github.com/SebastianGeroli/gator-boot-dev/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	cfg.SetUser("sebastian")
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", cfg)
}
