package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/transport/dto"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/PureTeamLead/go-test-assessment-developstoday/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CatService interface {
	CreateCat(ctx context.Context, req cat.CreateCatSvc) (uuid.UUID, error)
	DeleteCat(ctx context.Context, id uuid.UUID) error
	ListCats(ctx context.Context) ([]*cat.Cat, error)
	GetCatByID(ctx context.Context, id uuid.UUID) (*cat.Cat, error)
	GetCatByName(ctx context.Context, name string) (*cat.Cat, error)
	UpdateCat(ctx context.Context, params cat.UpdateCatParams) (*cat.Cat, error)
}

const idParam = "id"

func (h *Handler) GetCats(c *gin.Context) {
	const op = "handler.GetCats"

	cats, err := h.CatService.ListCats(h.Ctx)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, ErrorObj("internal server error"))
		return
	}

	c.JSON(http.StatusOK, cats)
	logger.GetLoggerFromCtx(h.Ctx).Info("success on list cats")
}

func (h *Handler) CreateCat(c *gin.Context) {
	const op = "handler.CreateCat"

	var req dto.CreateCatInput
	if err := c.BindJSON(&req); err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(fmt.Sprintf("%s: failed to map input on post request", op), err)
		c.JSON(http.StatusBadRequest, BadRequestObj())
		return
	}

	id, err := h.CatService.CreateCat(h.Ctx, cat.CreateCatSvc{
		Name:     req.Name,
		Breed:    req.Breed,
		YearsExp: req.ExperienceInYears,
		Salary:   req.SalaryCents,
	})
	switch {
	case errors.Is(err, utils.ErrConflictingData):
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj("duplicate of unique data"))
		return
	case errors.Is(err, utils.ErrInvalidBreed):
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj("invalid breed"))
		return
	case errors.Is(err, utils.ErrValidatingCat):
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj("failed to pass validation on cat object"))
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

func (h *Handler) GetCat(c *gin.Context) {
	const op = "handler.GetCat"

	id := c.Param(idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	fetchedCat, err := h.CatService.GetCatByID(h.Ctx, parsedID)
	if err != nil {
		if errors.Is(err, utils.ErrCatNotFound) {
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrCatNotFound.Error()))
			return
		}

		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusInternalServerError, InternalErrorObj())
		return
	}

	c.JSON(http.StatusOK, fetchedCat)
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) UpdateCat(c *gin.Context) {
	const op = "handler.GetCat"

	id := c.Param(idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	var req dto.UpdateCatInput
	if err = c.BindJSON(&req); err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(fmt.Sprintf("%s: failed to map input on put request", op), err)
		c.JSON(http.StatusBadRequest, BadRequestObj())
		return
	}

	updatedCat, err := h.CatService.UpdateCat(h.Ctx, cat.UpdateCatParams{
		ID:          parsedID,
		Name:        req.Name,
		Breed:       req.Breed,
		YearsXP:     req.ExperienceInYears,
		SalaryCents: req.SalaryCents,
	})
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidBreed):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj("invalid breed"))
			return
		case errors.Is(err, utils.ErrValidatingCat):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj("failed to pass validation on cat object"))
			return
		case errors.Is(err, utils.ErrCatNotFound):
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusBadRequest, ErrorObj("cat not found by that ID"))
			return
		default:
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, InternalErrorObj())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCat)
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}

func (h *Handler) DeleteCat(c *gin.Context) {
	const op = "handler.GetCat"

	id := c.Param(idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
		c.JSON(http.StatusBadRequest, ErrorObj(utils.ErrInvalidID.Error()))
		return
	}

	if err = h.CatService.DeleteCat(h.Ctx, parsedID); err != nil {
		if errors.Is(err, utils.ErrCatNotFound) {
			logger.GetLoggerFromCtx(h.Ctx).Error(op, err)
			c.JSON(http.StatusInternalServerError, ErrorObj("cat not found by that ID"))
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success on cat deletion operation"})
	logger.GetLoggerFromCtx(h.Ctx).Info(fmt.Sprintf("%s: success", op))
}
