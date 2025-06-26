package handler

import (
	"context"
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/server"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	CatService       CatService
	MisTargetService MisTargetService
	Router           *gin.Engine
	Server           *http.Server
}

const (
	catsPath    = "/cats"
	missionPath = "/missions"
	targetsPath = "/:mission_id/targets"
)

// TODO: Gin default?
func New(cfg server.Config, catService CatService, misTarService MisTargetService) *Handler {
	router := gin.New()
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Handler{CatService: catService, MisTargetService: misTarService, Router: router, Server: srv}
}

func (h *Handler) InitRoutes() {
	catsGroup := h.Router.Group(catsPath)
	{
		catsGroup.GET("", h.GetCats)
		catsGroup.GET("/:id", h.GetCat)
		catsGroup.POST("", h.CreateCat)
		catsGroup.DELETE("/:id", h.DeleteCat)
		catsGroup.PUT("/:id", h.UpdateCat)
	}

	missionsGroup := h.Router.Group(missionPath)
	{
		missionsGroup.GET("", h.ListMissions)
		missionsGroup.POST("", h.CreateMission)
		missionsGroup.GET("/:id", h.GetMission)
		missionsGroup.DELETE("/:id", h.DeleteMission)
		missionsGroup.PUT("/:id", h.UpdateMission)
		missionsGroup.PUT("/:id/assign", h.AssignMission)

		targetsGroup := missionsGroup.Group(targetsPath)
		targetsGroup.PUT("/:target_id", h.UpdateMissionTarget)
		targetsGroup.POST("/", h.AddMissionTarget)
		targetsGroup.DELETE("/:target_id", h.DeleteMissionTarget)
	}
}

func (h *Handler) Run() error {
	err := h.Server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Stop(ctx context.Context) error {
	const op = "HTTP server shutdown"
	if err := h.Server.Shutdown(ctx); err != nil {
		logger.GetLoggerFromCtx(ctx).Error("Fail on HTTP server shutdown", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.GetLoggerFromCtx(ctx).Info("Successful graceful shutdown of HTTP server")
	return nil
}
