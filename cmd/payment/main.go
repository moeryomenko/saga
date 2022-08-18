package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moeryomenko/healing"
	"github.com/moeryomenko/squad"

	"github.com/moeryomenko/saga/internal/payment/config"
	"github.com/moeryomenko/saga/internal/payment/infrastructure/eventhandler"
	"github.com/moeryomenko/saga/internal/payment/infrastructure/repository"
	"github.com/moeryomenko/saga/internal/payment/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, `load env vars: %s`, err)
		os.Exit(1)
	}

	group, err := squad.New(
		squad.WithSignalHandler(squad.WithGracefulPeriod(cfg.Health.GracePeriod)),
		squad.WithBootstrap(repository.Init(cfg), eventhandler.Init(cfg)),
		squad.WithCloses(repository.Close, eventhandler.Close),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, `bootstrap service: %s`, err)
		os.Exit(1)
	}

	health := healing.New(
		cfg.Health.Port,
		healing.WithCheckPeriod(cfg.Health.Period),
		healing.WithHealthzEndpoint(cfg.Health.LiveEndpoint),
		healing.WithReadyEndpoint(cfg.Health.ReadyEndpoint),
	)

	group.Run(eventhandler.HandleEvents(service.HandlePayments))
	group.RunGracefully(health.Heartbeat, health.Stop)

	errs := group.Wait()
	for _, err := range errs {
		log.Println(err)
	}
}
