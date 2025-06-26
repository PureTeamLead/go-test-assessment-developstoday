package app

import (
	"context"
	"flag"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/config"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/mission"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/service"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/logger"
	database "github.com/PureTeamLead/go-test-assessment-developstoday/pkg/storage/postgres"
	"log"
)

func Run() {
	ctx := context.Background()

	configPath := flag.String("config", "./configs/prod-config.yaml", "Specifying the path of the config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Config is not loaded: %s", err.Error())
	}

	ctx = logger.New(ctx, cfg.Env)

	db, err := database.NewPostgres(ctx, cfg.DBConfig)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal("Failed to set up database: " + err.Error())
	}

	catRepo := cat.NewRepository(db)
	missionRepo := mission.NewRepository(db)
	targetRepo := target.NewRepository(db)

	catSvc := cat.NewService(catRepo)
	misTarSvc := service.New(missionRepo, targetRepo)

	_ = catSvc
	_ = misTarSvc
}
