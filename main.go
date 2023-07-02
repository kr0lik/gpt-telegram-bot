package main

import (
	"context"
	"gpt-telegran-bot/internal/di"
	"gpt-telegran-bot/internal/di/config"
	"gpt-telegran-bot/internal/domain/usecase"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	RetryTimeout       = 60
	MaxInitializeCount = 10
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("cancel program on signal: %v", sig)
		cancel()
		os.Exit(0)
	}()

	run(ctx)
}

func run(ctx context.Context) {
	useCase, err, failCount := getUseCase()
	if err != nil {
		log.Fatalf("stop after %d fails", failCount)
	}

	for {
		if err := useCase.Start(ctx); err != nil {
			log.Printf("error on start: %v", err)
			log.Printf("will try restart after %d secconds", RetryTimeout)
			time.Sleep(RetryTimeout * time.Second)
		}
	}
}

func getUseCase() (*usecase.Messaging, error, int) {
	failCount := 0

	for {
		useCase, err := di.InitialiseMessaging()
		if err != nil {
			log.Printf("failed to initialize messaging: %v", err)

			if failCount >= MaxInitializeCount {
				return nil, err, failCount
			}

			failCount++

			log.Printf("will try restart initializing after %d secconds", RetryTimeout)
			time.Sleep(RetryTimeout * time.Second)
		}

		return useCase, err, failCount
	}
}

func init() {
	configPath := "config.yaml"

	if len(os.Args) > 1 && strings.Contains(os.Args[1], ".yaml") {
		configPath = os.Args[1]
	}

	if err := config.ReadConfig(configPath); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
}
