package main

import (
	"context"
	"errors"
	httpadapter "github.com/example/asset-maintenance-service/internal/adapter/http"
	"github.com/example/asset-maintenance-service/internal/application"
	memory "github.com/example/asset-maintenance-service/internal/infrastructure/memory"
	"github.com/example/asset-maintenance-service/internal/platform/config"
	"github.com/example/asset-maintenance-service/internal/platform/httpx"
	"github.com/example/asset-maintenance-service/internal/platform/observability"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load("configs/config.yaml")
	log := observability.NewLogger()
	metrics := &observability.Metrics{}
	assets := memory.NewAssetRepository()
	strategies := memory.NewStrategyRepository()
	templates := memory.NewTemplateRepository()
	plans := memory.NewPlanRepository()
	tasks := memory.NewTaskRepository()
	findings := memory.NewFindingRepository()
	history := memory.NewHistoryRepository()
	pub := memory.NewPublisher()
	svc := application.NewService(assets, strategies, templates, plans, tasks, findings, history, pub)
	handler := httpadapter.NewHandler(svc, metrics)
	server := &http.Server{Addr: cfg.Address, Handler: httpx.Chain(handler, log, metrics, cfg.RequestTimeout, cfg.RateLimit), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		log.Info("maintenance service started", "address", cfg.Address)
		if e := server.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			log.Error("server failed", "error", e)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if e := server.Shutdown(shutdown); e != nil {
		log.Error("graceful shutdown failed", "error", e)
	} else {
		log.Info("maintenance service stopped")
	}
}
