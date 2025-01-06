package main

import (
	"fmt"
	"os"

	"git.front.kjuulh.io/kjuulh/orbis/internal/app"
)

func main() {
	app := app.NewApp()

	if err := newRoot(app).Execute(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
