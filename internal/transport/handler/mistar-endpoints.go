package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/mission"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/service"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/dto"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const (
	missionIDParam = "id"
	targetIDParam  = "target-id"
)

type MisTargetService interface {
	CreateMission(ctx context.Context, rawTargets []service.CreateUpdateTargetSvc) (uuid.UUID, error)
	DeleteMission(ctx context.Context, id uuid.UUID) error
	UpdateMissionState(ctx context.Context, id uuid.UUID) error
	SetMissionTargetState(ctx context.Context, missionID uuid.UUID, targetID uuid.UUID) error
	UpdateMissionTargetNotes(ctx context.Context, missionID uuid.UUID, targetID uuid.UUID, notes string) error
	DeleteTargetFromMission(ctx context.Context, id uuid.UUID) error
	AddTargetToMission(ctx context.Context, missionID uuid.UUID, tarReq service.CreateUpdateTargetSvc) error
	AssignCatToMission(ctx context.Context, missionID uuid.UUID, catID uuid.UUID) error
	ListMissions(ctx context.Context) ([]*service.FullMission, error)
	GetMission(ctx context.Context, id uuid.UUID) (*service.FullMission, error)
}

func (h *Handler) CreateMission(c *gin.Context) {
	const op = "handler.CreateMission"

	var req dto.CreateMissionReq
	if err := c.BindJSON(&req); err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(fmt.Sprintf("%s: failed to map input on post request", op), err)
		c.JSON(http.StatusBadRequest, BadRequestObj())
		return
	}

	if len(req.Targets) > mission.MissionSize {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, utils.ErrTargetOverflow)
		c.JSON(http.StatusBadRequest, ErrorObj("too much targets specified"))
		return
	}

	id, err := h.MisTargetService.CreateMission(h.Ctx, dto.MapTargetsToRaw(req.Targets))
	switch {
	case errors.Is(err, utils.ErrNoTargets):
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj("no targets were specified"))
		return
	case errors.Is(err, utils.ErrConflictingData):
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj("duplicate of unique data"))
		return
	case errors.Is(err, utils.ErrValidatingTargets):
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj("failed to pass validation on target object"))
		return
	case err == nil:
		logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
		c.JSON(http.StatusCreated, map[string]interface{}{"obj_id": id})
		return
	default:
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, InternalErrorObj())
		return
	}
}

func (h *Handler) DeleteMission(c *gin.Context) {
	const op = "handler.DeleteMission"

	id := c.Param(idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	if err = h.MisTargetService.DeleteMission(h.Ctx, parsedID); err != nil {
		switch {
		case errors.Is(err, utils.ErrMissionNotFound):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, ErrorObj("mission object not found by that ID"))
			return
		case errors.Is(err, utils.ErrMissionCompleted):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, ErrorObj("failed to delete mission: it is already completed"))
			return
		case errors.Is(err, utils.ErrCatAssigned):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, ErrorObj("failed to delete mission: it is already assigned to the cat"))
			return
		default:
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, InternalErrorObj())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success on mission deletion operation"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) UpdateMissionState(c *gin.Context) {
	const op = "handler.UpdateMissionState"

	id := c.Param(idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	err = h.MisTargetService.UpdateMissionState(h.Ctx, parsedID)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMissionNotFound):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj("invalid breed"))
			return
		default:
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, InternalErrorObj())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "mission state updated"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) UpdateMissionTarget(c *gin.Context) {
	const op = "handler.UpdateMissionTarget"

	targetID := c.Param(targetIDParam)
	parsedTargetID, err := uuid.Parse(targetID)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	missionID := c.Param(missionIDParam)
	parsedMissionID, err := uuid.Parse(missionID)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	var req dto.UpdateMissionReq
	if err = c.BindJSON(&req); err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(fmt.Sprintf("%s: failed to map input on post request", op), err)
		c.JSON(http.StatusBadRequest, BadRequestObj())
		return
	}

	err = h.MisTargetService.UpdateMissionTargetNotes(h.Ctx, parsedMissionID, parsedTargetID, req.Notes)
	if err != nil {
		switch true {
		case errors.Is(err, utils.ErrMissionCompleted):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj("mission is already completed"))
			return
		case errors.Is(err, utils.ErrTargetCompleted):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj("target is already completed"))
			return
		case errors.Is(err, utils.ErrMissionNotFound):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrMissionNotFound.Error()))
			return
		default:
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, InternalErrorObj())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success on mission's target update"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) DeleteMissionTarget(c *gin.Context) {
	const op = "handler.DeleteMissionTarget"

	targetID := c.Param(targetIDParam)
	parsedTargetID, err := uuid.Parse(targetID)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	err = h.MisTargetService.DeleteTargetFromMission(h.Ctx, parsedTargetID)
	if err != nil {
		if errors.Is(err, utils.ErrTargetNotFound) {
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrTargetNotFound.Error()))
			return
		}

		if errors.Is(err, utils.ErrTargetCompleted) {
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrTargetCompleted.Error()))
			return
		}

		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, InternalErrorObj())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success on mission's target deletion"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) AddMissionTarget(c *gin.Context) {
	const op = "handler.AddMissionTarget"

	missionID := c.Param(missionIDParam)
	parsedMissionID, err := uuid.Parse(missionID)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	var req dto.CreateTargetReq
	if err = c.BindJSON(&req); err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(fmt.Sprintf("%s: failed to map input on post request", op), err)
		c.JSON(http.StatusBadRequest, BadRequestObj())
		return
	}

	err = h.MisTargetService.AddTargetToMission(h.Ctx, parsedMissionID, dto.MapTargetToRaw(req))
	if err != nil {
		switch true {
		case errors.Is(err, utils.ErrMissionCompleted):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrMissionCompleted.Error()))
			return
		case errors.Is(err, utils.ErrMissionNotFound):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrMissionNotFound.Error()))
			return
		case errors.Is(err, utils.ErrTargetNotFound):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrTargetNotFound.Error()))
			return
		case errors.Is(err, utils.ErrTargetOverflow):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrTargetOverflow.Error()))
			return
		default:
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, InternalErrorObj())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success on adding target to the mission"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) AssignMission(c *gin.Context) {
	const op = "handler.AssignMission"

	rawMissionID := c.Param(idParam)
	missionID, err := uuid.Parse(rawMissionID)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	var req dto.AssignToMissionReq
	if err = c.BindJSON(&req); err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(fmt.Sprintf("%s: failed to map input on post request", op), err)
		c.JSON(http.StatusBadRequest, BadRequestObj())
		return
	}

	catID, err := uuid.Parse(req.CatID)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	err = h.MisTargetService.AssignCatToMission(h.Ctx, missionID, catID)
	if err != nil {
		if errors.Is(err, utils.ErrMissionNotFound) {
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, utils.ErrMissionNotFound.Error())
			return
		}

		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, InternalErrorObj())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success on assigning cat to the mission"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) ListMissions(c *gin.Context) {
	const op = "handler.ListMissions"

	missions, err := h.MisTargetService.ListMissions(h.Ctx)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, ErrorObj("internal server error"))
		return
	}

	c.JSON(http.StatusOK, missions)
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) GetMission(c *gin.Context) {
	const op = "handler.GetMission"

	id := c.Param(idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	fetchedMission, err := h.MisTargetService.GetMission(h.Ctx, parsedID)
	if err != nil {
		if errors.Is(err, utils.ErrMissionNotFound) {
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrMissionNotFound.Error()))
			return
		}

		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, InternalErrorObj())
		return
	}

	c.JSON(http.StatusOK, fetchedMission)
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}
