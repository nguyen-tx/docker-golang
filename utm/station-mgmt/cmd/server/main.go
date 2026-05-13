package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/utm/station-mgmt/internal/client"
	"github.com/utm/station-mgmt/internal/config"
	"github.com/utm/station-mgmt/internal/gapi"
	"github.com/utm/station-mgmt/internal/hapi"
	"github.com/utm/station-mgmt/internal/service"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	// khởi tạo gRPC clients đến 3 service con
	nsClient, err := client.NewNetworkStationClient(cfg.NetworkStation.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("connect network station")
	}
	defer nsClient.Close()

	rtkClient, err := client.NewRTKStationClient(cfg.RTKStation.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("connect rtk station")
	}
	defer rtkClient.Close()

	miClient, err := client.NewManageInfoClient(cfg.ManageInfo.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("connect manage info")
	}
	defer miClient.Close()

	overviewSvc := service.NewOverviewService(nsClient, rtkClient, miClient)

	// HTTP server
	httpRouter := hapi.NewRouter(cfg, overviewSvc)
	httpSrv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: httpRouter,
	}

	// gRPC server
	grpcSrv := gapi.NewServer(cfg.App.GRPCPort, overviewSvc)

	go func() {
		log.Info().Str("port", cfg.App.Port).Msg("HTTP server starting")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server")
		}
	}()

	go func() {
		if err := grpcSrv.Run(); err != nil {
			log.Fatal().Err(err).Msg("grpc server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
	grpcSrv.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
}
