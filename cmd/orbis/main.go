package main

import (
	"fmt"
	"os"

	"git.front.kjuulh.io/kjuulh/orbis/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	app := app.NewApp()

	if err := newRoot(app).Execute(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
