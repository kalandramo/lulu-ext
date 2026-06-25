package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kalandramo/lulu"
	httpServer "github.com/kalandramo/lulu-ext/transport/http"
	ginDriver "github.com/kalandramo/lulu-ext/transport/http/driver/gin"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/cors"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/recovery"
)

func main() {
	// 1. HTTP server (port 8080)
	httpSrv := httpServer.NewServer(
		":8080",
		httpServer.WithDriver(ginDriver.NewDriver()),
	)

	httpSrv.Use(
		recovery.Middleware(),
		cors.Middleware(),
	)

	httpSrv.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		id := "req-123456"
		fmt.Fprintf(w, "Hello from gin-server app! (request-id: %s)\n", id)
	})

	// 2. Assemble the App — one place to manage all servers.
	app := lulu.New(
		lulu.WithID("gin-server-demo"),
		lulu.WithName("multi-server-app"),
		lulu.WithVersion("1.0.0"),
		lulu.WithServer(httpSrv),
		lulu.WithStopTimeout(15*time.Second),
		// BeforeStop: runs BEFORE any server stops.
		// Typical use: deregister from service registry, drain request queue.
		lulu.WithBeforeStop(func(_ context.Context) error {
			fmt.Println("[hook] beforeStop: preparing for graceful shutdown...")
			return nil
		}),
		// AfterStop: runs AFTER all servers have stopped.
		// Typical use: close database connections, flush log buffers.
		lulu.WithAfterStop(func(_ context.Context) error {
			fmt.Println("[hook] afterStop: cleanup complete, exiting.")
			return nil
		}),
	)

	fmt.Println("Starting gin-server application...")
	fmt.Printf("  HTTP:  %s\n", httpSrv.Endpoint())
	fmt.Println()

	// App.Run blocks until SIGINT/SIGTERM/SIGQUIT or context cancellation.
	// No manual signal.NotifyContext needed!
	if err := app.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "application exited with error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("application stopped gracefully")
}
