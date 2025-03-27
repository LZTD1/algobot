package main

import (
	"algobot/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Printf("%#v\n", cfg)

	// TODO : create application

	// TODO : start bot app
	// TODO : start message scheduler app

	// graceful shutdown

	// TODO : add graceful shutdown
}
