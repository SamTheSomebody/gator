package main

import (
	"internal/config"
	"fmt"
)

func main() {
	cfg := config.Read()
	config.SetUser("Sam", cfg)
	cfg = config.Read()
	fmt.Printf("DB URL: %v\n", cfg.Db_url)
	fmt.Printf("Current User Name: %v\n", cfg.Current_user_name)
}
