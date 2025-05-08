package attack

import (
	"load-generation-system/api/rest/manager/handlers/model"
	"load-generation-system/internal/core"
	"slices"
	"sort"
)

type StartAttackPresenter model.StartAttackRequestBody

type StartIncrementPresenter model.StartIncrementRequestBody

func PresentScenario(scenario core.ScenarioDetails) model.ScenarioInfo {
	return model.ScenarioInfo{
		Name:        scenario.Name,
		Description: scenario.Description,
	}
}

func PresentIncrement(increment core.IncrementDetails) model.IncrementInfo {
	scenarioCounters := make([]model.ScenarioCounter, 0, len(increment.Scenarios))
	for scenario, counter := range increment.Scenarios {
		scenarioCounters = append(scenarioCounters, model.ScenarioCounter{
			Scenario: scenario,
			Counter:  counter,
		})
	}
	sort.Slice(scenarioCounters, func(i, j int) bool {
		return scenarioCounters[i].Scenario < scenarioCounters[j].Scenario
	})

	return model.IncrementInfo{
		ID:        increment.ID,
		Scenarios: scenarioCounters,
	}
}

func PresentAttack(attack core.AttackDetails) model.AttackInfo {
	incrementInfos := make([]model.IncrementInfo, 0, len(attack.Increments))
	for _, increment := range attack.Increments {
		incrementInfos = append(incrementInfos, PresentIncrement(increment))
	}
	sort.Slice(incrementInfos, func(i, j int) bool {
		return incrementInfos[i].ID < incrementInfos[j].ID
	})

	var constConfig *model.ConstConfig
	if attack.ConstConfig != nil {
		scenarios := make(map[string]int64)
		for _, scenario := range attack.ConstConfig.Scenarios {
			scenarios[scenario.Name] = *scenario.Counter
		}
		constConfig = &model.ConstConfig{
			Scenarios: scenarios,
		}
	}

	var linearConfig *model.LinearConfig
	if attack.LinearConfig != nil {
		scenarios := make([]string, 0, len(attack.LinearConfig.Scenarios))
		for _, scenario := range attack.LinearConfig.Scenarios {
			scenarios = append(scenarios, scenario.Name)
		}

		linearConfig = &model.LinearConfig{
			WarmUpSec:       attack.LinearConfig.WarmUpSec,
			StartCounter:    attack.LinearConfig.StartCounter,
			EndCounter:      attack.LinearConfig.EndCounter,
			CounterStep:     attack.LinearConfig.CounterStep,
			StepIntervalSec: attack.LinearConfig.StepIntervalSec,
			Scenarios:       scenarios,
		}
	}

	return model.AttackInfo{
		Name:         attack.Name,
		ID:           attack.ID,
		WaitTimeSec:  attack.WaitTimeSec,
		CreatedAt:    attack.CreatedAt,
		DurationSec:  attack.DurationSec,
		ConstConfig:  constConfig,
		LinearConfig: linearConfig,
		Increments:   incrementInfos,
	}
}

func PresentNode(node core.NodeDetails) model.NodeInfo {
	attackPresenters := make([]model.AttackInfo, 0, len(node.Attacks))
	for _, attack := range node.Attacks {
		attackPresenters = append(attackPresenters, PresentAttack(attack))
	}
	sort.Slice(attackPresenters, func(i, j int) bool {
		return attackPresenters[i].ID < attackPresenters[j].ID
	})
	scenarios := make([]string, 0, len(node.Scenarios))
	for _, scenario := range node.Scenarios {
		scenarios = append(scenarios, scenario.Name)
	}
	sort.Strings(scenarios)

	return model.NodeInfo{
		Name:      node.Name,
		Scenarios: scenarios,
		Attacks:   attackPresenters,
		IsActive:  node.IsActive,
	}
}

func PresentNodeList(nodes []core.NodeDetails) []model.NodeInfo {
	pres := make([]model.NodeInfo, 0, len(nodes))
	for _, node := range nodes {
		pres = append(pres, PresentNode(node))
	}
	sort.Slice(pres, func(i, j int) bool {
		return pres[i].Name < pres[j].Name
	})

	return pres
}

func PresentAttackList(attacks []core.AttackDetails) []model.AttackInfo {
	pres := make([]model.AttackInfo, 0, len(attacks))
	for _, attack := range attacks {
		pres = append(pres, PresentAttack(attack))
	}
	sort.Slice(pres, func(i, j int) bool {
		return pres[i].ID < pres[j].ID
	})

	return pres
}

func PresentScenarioList(scenarios []core.ScenarioDetails) []model.ScenarioInfo {
	pres := make([]model.ScenarioInfo, 0, len(scenarios))
	for _, scenario := range scenarios {
		pres = append(pres, PresentScenario(scenario))
	}
	sort.Slice(pres, func(i, j int) bool {
		return pres[i].Name < pres[j].Name
	})

	return pres
}

func (sa *StartAttackPresenter) ToCore() (core.StartAttack, error) {
	if sa.ConstConfig == nil && sa.LinearConfig == nil {
		return core.StartAttack{}, core.ErrBadConfig
	}

	var constConfig *core.ConstConfig
	if sa.ConstConfig != nil {
		scenarios := make([]core.Scenario, 0, len(sa.ConstConfig.Scenarios))
		for scenario, counter := range sa.ConstConfig.Scenarios {
			scenarios = append(scenarios, core.Scenario{
				Name:    scenario,
				Counter: &counter,
			})
		}
		constConfig = &core.ConstConfig{
			Scenarios: scenarios,
		}
	}

	var linearConfig *core.LinearConfig
	if sa.LinearConfig != nil {
		if sa.LinearConfig.WarmUpSec == nil && sa.LinearConfig.CounterStep == nil {
			return core.StartAttack{}, core.ErrBadConfig
		}

		if sa.LinearConfig.EndCounter <= sa.LinearConfig.StartCounter {
			return core.StartAttack{}, core.ErrBadConfig
		}

		if sa.LinearConfig.WarmUpSec != nil && sa.LinearConfig.CounterStep != nil {
			computedInterval := float64((*sa.LinearConfig.WarmUpSec * (*sa.LinearConfig.CounterStep)) / (sa.LinearConfig.EndCounter - sa.LinearConfig.StartCounter))
			if computedInterval < 1 {
				return core.StartAttack{}, core.ErrBadConfig
			}
		}

		if sa.LinearConfig.WarmUpSec != nil && sa.DurationSec != nil && *sa.LinearConfig.WarmUpSec >= *sa.DurationSec {
			return core.StartAttack{}, core.ErrBadConfig
		}

		if sa.LinearConfig.WarmUpSec != nil && sa.LinearConfig.StepIntervalSec != nil && *sa.LinearConfig.StepIntervalSec > *sa.LinearConfig.WarmUpSec {
			return core.StartAttack{}, core.ErrBadConfig
		}

		if sa.DurationSec != nil && sa.LinearConfig.StepIntervalSec != nil && *sa.LinearConfig.StepIntervalSec >= *sa.DurationSec {
			return core.StartAttack{}, core.ErrBadConfig
		}

		scenarios := make([]core.Scenario, 0, len(sa.LinearConfig.Scenarios))
		for _, scenario := range sa.LinearConfig.Scenarios {
			scenarios = append(scenarios, core.Scenario{
				Name: scenario,
			})
		}

		linearConfig = &core.LinearConfig{
			WarmUpSec:       sa.LinearConfig.WarmUpSec,
			StartCounter:    sa.LinearConfig.StartCounter,
			EndCounter:      sa.LinearConfig.EndCounter,
			CounterStep:     sa.LinearConfig.CounterStep,
			StepIntervalSec: sa.LinearConfig.StepIntervalSec,
			Scenarios:       scenarios,
		}
	}

	return core.StartAttack{
		Name:         sa.Name,
		WaitTimeSec:  sa.WaitTimeSec,
		DurationSec:  sa.DurationSec,
		ConstConfig:  constConfig,
		LinearConfig: linearConfig,
	}, nil
}

func (si *StartIncrementPresenter) ToCore(attackID int64) core.OperationStart {
	return core.OperationStart{
		AttackID:  attackID,
		Scenarios: si.Scenarios,
	}
}

func PresentStringList(list []string) []string {
	slices.Sort(list)

	return list
}
