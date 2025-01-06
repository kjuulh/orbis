package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop

		app.Logger().Info("stop signal received: shutting down orbis")
		cancel()

		// Start timer for hard stop
		time.Sleep(time.Second * 10)
		fmt.Println("orbis failed to stop in time, forced to hard cancel")
		os.Exit(1)
	}()

	if err := newRoot(app).ExecuteContext(ctx); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
