package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calendar/internal/app"
	"calendar/internal/logger"
	internalhttp "calendar/internal/server/http"
	memorystorage "calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig("config.json")
	if err != nil {
		if errors.Is(err, ErrInvalidConfig) {
			fmt.Println(fmt.Errorf("error validating confing %w", err))
			config = NewDefaultConfig()
		} else {
			fmt.Printf("Error loading confing %s", err.Error())
			return
		}
	}

	logg := logger.New(config.Logger)
	defer logg.Close()

	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(config.Server, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
