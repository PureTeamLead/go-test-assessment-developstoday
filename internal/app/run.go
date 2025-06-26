package app

import (
	"context"
	"flag"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/config"
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

	db, err := database.NewPostgres(ctx, cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to set up database: %s", err.Error())
	}

	_ = db
}
