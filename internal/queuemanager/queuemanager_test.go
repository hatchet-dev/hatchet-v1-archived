package queuemanager_test

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestEnqueueSuccessful(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		queueManager := queuemanager.NewDefaultModuleRunQueueManager(conf.Repository)

		// create a module run
		moduleRun := &models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindPlan,
		}

		moduleRun, err := conf.Repository.Module().CreateModuleRun(moduleRun)

		if err != nil {
			return err
		}

		err = queueManager.Enqueue(testutils.ModuleModels[0], moduleRun)

		assert.Nil(t, err)

		queueLen, err := queueManager.Len(testutils.ModuleModels[0])

		assert.Nil(t, err)
		assert.Equal(t, 1, queueLen, "queue length should be 1")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams, testutils.InitModules)
}

func TestPeek(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		queueManager := queuemanager.NewDefaultModuleRunQueueManager(conf.Repository)

		// create a module run
		moduleRun1 := &models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindApply,
		}

		moduleRun1, err := conf.Repository.Module().CreateModuleRun(moduleRun1)

		if err != nil {
			return err
		}

		moduleRun2 := &models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindPlan,
		}

		moduleRun2, err = conf.Repository.Module().CreateModuleRun(moduleRun2)

		if err != nil {
			return err
		}

		err = queueManager.Enqueue(testutils.ModuleModels[0], moduleRun1)

		assert.Nil(t, err)

		err = queueManager.Enqueue(testutils.ModuleModels[0], moduleRun2)

		assert.Nil(t, err)

		queueLen, err := queueManager.Len(testutils.ModuleModels[0])

		assert.Nil(t, err)
		assert.Equal(t, 2, queueLen, "queue length should be 2")

		// peek should return the apply
		gotModuleRun, err := queueManager.Peek(testutils.ModuleModels[0])

		assert.Nil(t, err)

		assert.Equal(t, moduleRun1.ID, gotModuleRun.ID, "module run ids are equal")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams, testutils.InitModules)
}
