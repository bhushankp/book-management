package main

import (
	"book-management/internal/config"
	"fmt"
)

func main() {

	cfg := config.LoadConfig()
	fmt.Printf("server Port : %d", cfg.Server.Port)

}
