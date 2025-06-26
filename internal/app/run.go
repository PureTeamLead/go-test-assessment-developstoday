package app

import (
	"context"
	"flag"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/config"
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

}
