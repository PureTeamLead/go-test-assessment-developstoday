package utils

import "errors"

var (
	ErrValidatingCat     = errors.New("invalid cat input structure")
	ErrCatNotFound       = errors.New("cat not found")
	ErrConflictingData   = errors.New("conflict of data occurred")
	ErrMissionNotFound   = errors.New("mission not found")
	ErrTargetNotFound    = errors.New("target not found")
	ErrInvalidBreed      = errors.New("invalid breed")
	ErrApiServerError    = errors.New("api server error")
	ErrNoTargets         = errors.New("empty targets")
	ErrValidatingTargets = errors.New("failed to validate and create targets")
	ErrCatAssigned       = errors.New("cat is already assigned to the mission, operation is impossible")
	ErrMissionCompleted  = errors.New("mission is already completed, operation is impossible")
	ErrTargetCompleted   = errors.New("target is already completed, operation is impossible")
	ErrInvalidID         = errors.New("invalid ID format")
)
