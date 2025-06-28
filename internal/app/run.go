package app

import (
	"context"
	"flag"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/config"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/mission"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/service"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/handler"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/logger"
	database "github.com/PureTeamLead/go-test-assessment-developstoday/pkg/storage/postgres"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeoutDuration = 10 * time.Second
	portCtx         = "port"
)

func Run() {
	ctx := context.Background()

	configPath := flag.String("config", "./configs/prod-config.yaml", "Specifying the path of the config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Config is not loaded: %s", err.Error())
	}

	log.Println(cfg.DBConfig.Username, cfg.DBConfig.Password)

	ctx = logger.New(ctx, cfg.Env)
	ctx = context.WithValue(ctx, portCtx, cfg.HTTPSrvConfig.Port)

	db, err := database.NewPostgres(ctx, cfg.DBConfig)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal("Failed to set up database: " + err.Error())
	}
	defer db.Close()

	catRepo := cat.NewRepository(db)
	missionRepo := mission.NewRepository(db)
	targetRepo := target.NewRepository(db)

	catSvc := cat.NewService(catRepo)
	misTarSvc := service.New(missionRepo, targetRepo)

	logger.GetLoggerFromCtx(ctx).WithPort(ctx, portCtx)
	transport := handler.New(ctx, cfg.HTTPSrvConfig, catSvc, misTarSvc)
	transport.InitRoutes()

	go func() {
		if err = transport.Run(); err != nil {
			logger.GetLoggerFromCtx(ctx).Error("HTTP server stopped", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	sigStr := <-sig
	logger.GetLoggerFromCtx(ctx).Info("Interrupted by signal", zap.String("type", sigStr.String()))

	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	if err = transport.Stop(ctx); err != nil {
		logger.GetLoggerFromCtx(ctx).Error("failed server shutdown", err)
	}
}
