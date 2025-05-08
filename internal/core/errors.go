package core

import "errors"

var (
	ErrNoActiveNodes     = errors.New("there are no active nodes to perform the attack")
	ErrAttackNotFound    = errors.New("attack not found")
	ErrIncrementNotFound = errors.New("increment not found")
	ErrScenarioNotFound  = errors.New("scenario not found")
	ErrEmptyAttack       = errors.New("empty attack configuration")
	ErrBadConfig         = errors.New("bad attack configuration")
	ErrNodeAlreadyExists = errors.New("node already exists")

	ErrJobNotFound     = errors.New("job not found")
	ErrBrokenScheduler = errors.New("cannot schedule report job")

	ErrUnacceptableCode = errors.New("unacceptable status code")

	ErrScenarioExecutionViolation = errors.New("scenario execution violation")
)
