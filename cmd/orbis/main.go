package main

import (
	"context"
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

	ctx := context.Background()
	if err := newRoot(app).ExecuteContext(ctx); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
