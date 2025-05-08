package model

import (
	"time"
)

type StartAttackRequestBody struct {
	Name         string        `json:"name" example:"string" validate:"required"`
	WaitTimeSec  float64       `json:"wait_time_sec" example:"1" validate:"min=0.1,max=30"`
	DurationSec  *int64        `json:"duration_sec" example:"1" validate:"omitempty,min=1,max=2592000"`
	ConstConfig  *ConstConfig  `json:"const_config"`
	LinearConfig *LinearConfig `json:"linear_config"`
}

type ConstConfig struct {
	Scenarios map[string]int64 `json:"scenarios" validate:"required"`
}

type LinearConfig struct {
	WarmUpSec       *int64   `json:"warm_up_sec,omitempty" example:"1" validate:"omitempty,min=1"`
	StartCounter    int64    `json:"start_counter" example:"1" validate:"min=1"`
	EndCounter      int64    `json:"end_counter" example:"1" validate:"min=1"`
	CounterStep     *int64   `json:"counter_step,omitempty" example:"1" validate:"omitempty,min=1"`
	StepIntervalSec *int64   `json:"step_interval_sec,omitempty" example:"1" validate:"omitempty,min=1"`
	Scenarios       []string `json:"scenarios" validate:"required"`
}

type StartIncrementRequestBody struct {
	Scenarios map[string]int64 `json:"scenarios"`
}

type ScenarioInfo struct {
	Name        string `json:"name" example:"string"`
	Description string `json:"description" example:"string"`
}

type IncrementInfo struct {
	ID        int64             `json:"id" example:"1"`
	Scenarios []ScenarioCounter `json:"scenarios"`
}

type ScenarioCounter struct {
	Scenario string `json:"scenario" example:"string"`
	Counter  int64  `json:"counter" example:"1"`
}

type AttackInfo struct {
	ID           int64           `json:"id" example:"1"`
	Name         string          `json:"name" example:"string"`
	WaitTimeSec  float64         `json:"wait_time_sec" example:"1"`
	CreatedAt    time.Time       `json:"created_at" example:"2024-09-02T13:54:00Z"`
	DurationSec  *int64          `json:"duration_sec,omitempty" example:"1"`
	ConstConfig  *ConstConfig    `json:"const_config"`
	LinearConfig *LinearConfig   `json:"linear_config"`
	Increments   []IncrementInfo `json:"increments"`
}

type StartAttackResponse struct {
	Status string     `json:"status" example:"OK"`
	Attack AttackInfo `json:"data"`
}

type StartIncrementResponse struct {
	Status string        `json:"status" example:"OK"`
	Attack IncrementInfo `json:"data"`
}

type GetScenariosResponse struct {
	Status    string         `json:"status" example:"OK"`
	Scenarios []ScenarioInfo `json:"data"`
}

type GetAttacksResponse struct {
	Status  string       `json:"status" example:"OK"`
	Attacks []AttackInfo `json:"data"`
}

type NodeInfo struct {
	Name      string       `json:"name" example:"string"`
	Scenarios []string     `json:"scenarios"`
	Attacks   []AttackInfo `json:"attacks"`
	IsActive  bool         `json:"is_active" example:"true"`
}

type GetNodesResponse struct {
	Status string     `json:"status" example:"OK"`
	Nodes  []NodeInfo `json:"data"`
}

type StopAttackResponse struct {
	Status string `json:"status" example:"OK"`
}

type StopIncrementResponse struct {
	Status string `json:"status" example:"OK"`
}

type NoContentResponse struct{}
