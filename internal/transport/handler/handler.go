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
	Ctx              context.Context
}

const (
	catsPath    = "/cats"
	missionPath = "/missions"
	targetsPath = "/:id/targets"
)

func New(ctx context.Context, cfg server.Config, catService CatService, misTarService MisTargetService) *Handler {
	router := gin.New()
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Handler{Ctx: ctx, CatService: catService, MisTargetService: misTarService, Router: router, Server: srv}
}

func (h *Handler) InitRoutes() {
	defer h.assignRouter()

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
		targetsGroup.PUT("/:target-id", h.UpdateMissionTarget)
		targetsGroup.POST("/", h.AddMissionTarget)
		targetsGroup.DELETE("/:target-id", h.DeleteMissionTarget)
	}
}

func (h *Handler) assignRouter() {
	h.Server.Handler = h.Router
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
