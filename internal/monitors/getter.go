package monitors

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

func GetAllMonitorsForModuleRun(repo repository.Repository, teamID string, moduleRun *models.ModuleRun) ([]*models.ModuleMonitor, error) {
	if moduleRun.Kind == models.ModuleRunKindMonitor {
		monitors, _, err := repo.ModuleMonitor().ListModuleMonitorsByTeamID(teamID)

		if err != nil {
			return nil, err
		}

		res := make([]*models.ModuleMonitor, 0)

		for _, monitor := range monitors {
			if monitor.ShouldRunForModule(moduleRun.ModuleID) {
				cMonitor := *monitor

				res = append(res, &cMonitor)
			}
		}

		return res, nil
	}

	before, after, err := GetInlineMonitorsForModuleRun(repo, teamID, moduleRun)

	if err != nil {
		return nil, err
	}

	return append(before, after...), nil
}

func GetInlineMonitorsForModuleRun(repo repository.Repository, teamID string, moduleRun *models.ModuleRun) ([]*models.ModuleMonitor, []*models.ModuleMonitor, error) {
	monitors, _, err := repo.ModuleMonitor().ListModuleMonitorsByTeamID(teamID)

	if err != nil {
		return nil, nil, err
	}

	beforeMonitors := make([]*models.ModuleMonitor, 0)
	afterMonitors := make([]*models.ModuleMonitor, 0)

	for _, monitor := range monitors {
		if monitor.ShouldRunForModule(moduleRun.ModuleID) {
			cMonitor := *monitor

			k := monitor.Kind

			matchedBefore := (k == models.MonitorKindBeforeApply && moduleRun.Kind == models.ModuleRunKindApply) ||
				(k == models.MonitorKindBeforePlan && moduleRun.Kind == models.ModuleRunKindPlan) ||
				(k == models.MonitorKindBeforeDestroy && moduleRun.Kind == models.ModuleRunKindDestroy)

			matchedAfter := (k == models.MonitorKindAfterApply && moduleRun.Kind == models.ModuleRunKindApply) ||
				(k == models.MonitorKindAfterPlan && moduleRun.Kind == models.ModuleRunKindPlan) ||
				(k == models.MonitorKindAfterDestroy && moduleRun.Kind == models.ModuleRunKindDestroy)

			if matchedBefore {
				beforeMonitors = append(beforeMonitors, &cMonitor)
			}

			if matchedAfter {
				afterMonitors = append(afterMonitors, &cMonitor)
			}
		}
	}

	return beforeMonitors, afterMonitors, nil
}
