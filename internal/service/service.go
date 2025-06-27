package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/cat"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/mission"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/domain/target"
	"github.com/PureTeamLead/go-test-assessment-developstoday/internal/utils"
	"github.com/google/uuid"
)

type MissionRepository interface {
	AddMission(ctx context.Context) (uuid.UUID, error)
	DeleteMission(ctx context.Context, id uuid.UUID) error
	SetMissionCompleted(ctx context.Context, id uuid.UUID) error
	AddCatID(ctx context.Context, missionID uuid.UUID, catID uuid.UUID) error
	GetMissions(ctx context.Context) ([]*mission.Mission, error)
	GetMissionByCatID(ctx context.Context, catID uuid.UUID) (*mission.Mission, error)
	GetMissionByID(ctx context.Context, id uuid.UUID) (*mission.Mission, error)
	GetAssignedCat(ctx context.Context, missionID uuid.UUID) (*cat.Cat, error)
}

type TargetRepository interface {
	GetTargetsByMissionID(ctx context.Context, missionID uuid.UUID) ([]*target.Target, error)
	GetTargetByID(ctx context.Context, id uuid.UUID) (*target.Target, error)
	UpdateTargetNotes(ctx context.Context, id uuid.UUID, notes string) error
	SetTargetCompleted(ctx context.Context, id uuid.UUID) error
	AddTarget(ctx context.Context, missionID uuid.UUID, target *target.Target) (uuid.UUID, error)
	DeleteTarget(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	mr MissionRepository
	tr TargetRepository
}

const completedState = "completed"

func New(mr MissionRepository, tr TargetRepository) *Service {
	return &Service{mr: mr, tr: tr}
}

func (s *Service) CreateMission(ctx context.Context, rawTargets []CreateUpdateTargetSvc) (uuid.UUID, error) {
	const op = "service.CreateMission"

	if len(rawTargets) == 0 {
		return uuid.Nil, utils.ErrNoTargets
	}

	var validatedTargets []*target.Target
	for _, rawTarget := range rawTargets {
		notes := ""
		if rawTarget.Notes != nil {
			notes = *(rawTarget.Notes)
		}

		newTarget := target.NewEntity(rawTarget.Name, rawTarget.Country, notes)
		if err := newTarget.Validate(); err != nil {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}

		validatedTargets = append(validatedTargets, newTarget)
	}

	if len(validatedTargets) != len(rawTargets) {
		return uuid.Nil, utils.ErrValidatingTargets
	}

	id, err := s.mr.AddMission(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, validTarget := range validatedTargets {
		if _, err = s.tr.AddTarget(ctx, id, validTarget); err != nil {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}

	}

	return id, nil
}

func (s *Service) DeleteMission(ctx context.Context, id uuid.UUID) error {
	const op = "service.DeleteMission"

	targets, err := s.tr.GetTargetsByMissionID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, tar := range targets {
		if err = s.tr.DeleteTarget(ctx, tar.ID); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = s.mr.DeleteMission(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateMissionState(ctx context.Context, id uuid.UUID) error {
	const op = "service.UpdateMissionState"

	err := s.mr.SetMissionCompleted(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) SetMissionTargetState(ctx context.Context, missionID uuid.UUID, targetID uuid.UUID) error {
	const op = "service.SetMissionTargetState"

	mis, err := s.mr.GetMissionByID(ctx, missionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if mis.State == completedState {
		return utils.ErrMissionCompleted
	}

	if err = s.tr.SetTargetCompleted(ctx, targetID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateMissionTargetNotes(ctx context.Context, missionID uuid.UUID, targetID uuid.UUID, notes string) error {
	const op = "service.SetMissionTargetState"

	mis, err := s.mr.GetMissionByID(ctx, missionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if mis.State == completedState {
		return utils.ErrMissionCompleted
	}

	tar, err := s.tr.GetTargetByID(ctx, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tar.State == completedState {
		return utils.ErrTargetCompleted
	}

	if err = s.tr.UpdateTargetNotes(ctx, targetID, notes); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteTargetFromMission(ctx context.Context, id uuid.UUID) error {
	const op = "service.DeleteTargetsFromMission"

	tar, err := s.tr.GetTargetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tar.State == completedState {
		return utils.ErrTargetCompleted
	}

	if err = s.tr.DeleteTarget(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AddTargetToMission(ctx context.Context, missionID uuid.UUID, tarReq CreateUpdateTargetSvc) error {
	const op = "service.AddTargetToMission"

	mis, err := s.mr.GetMissionByID(ctx, missionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	targets, err := s.tr.GetTargetsByMissionID(ctx, missionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(targets) == mission.MissionSize {
		return utils.ErrTargetOverflow
	}

	if mis.State == completedState {
		return utils.ErrMissionCompleted
	}

	tar := MapTargetSvcToEntity(tarReq)

	if _, err = s.tr.AddTarget(ctx, missionID, tar); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AssignCatToMission(ctx context.Context, missionID uuid.UUID, catID uuid.UUID) error {
	const op = "service.AssignCatToMission"

	if err := s.mr.AddCatID(ctx, missionID, catID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ListMissions(ctx context.Context) ([]*FullMission, error) {
	const op = "service.ListMission"

	missions, err := s.mr.GetMissions(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fullMissions := make([]*FullMission, 0, len(missions))

	for _, mis := range missions {
		var fullMis *FullMission
		if fullMis, err = s.GetMission(ctx, mis.ID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		fullMissions = append(fullMissions, fullMis)
	}

	return fullMissions, nil
}

func (s *Service) GetMission(ctx context.Context, id uuid.UUID) (*FullMission, error) {
	const op = "service.GetMission"
	var fullMis FullMission

	mis, err := s.mr.GetMissionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	targets, err := s.tr.GetTargetsByMissionID(ctx, mis.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fullMis.ID = mis.ID
	fullMis.State = mis.State
	fullMis.CreatedAt = mis.CreatedAt
	fullMis.UpdatedAt = mis.UpdatedAt
	fullMis.Targets = targets

	assignedCat, err := s.mr.GetAssignedCat(ctx, id)
	if err != nil {
		if !errors.Is(err, utils.ErrCatNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return &fullMis, nil
	}

	fullMis.Cat = assignedCat
	return &fullMis, nil
}
