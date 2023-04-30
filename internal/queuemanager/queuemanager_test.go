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

		err = queueManager.Enqueue(testutils.ModuleModels[0], moduleRun, &queuemanager.LockOpts{})

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

		err = queueManager.Enqueue(testutils.ModuleModels[0], moduleRun1, &queuemanager.LockOpts{})

		assert.Nil(t, err)

		err = queueManager.Enqueue(testutils.ModuleModels[0], moduleRun2, &queuemanager.LockOpts{})

		assert.Nil(t, err)

		queueLen, err := queueManager.Len(testutils.ModuleModels[0])

		assert.Nil(t, err)
		assert.Equal(t, 2, queueLen, "queue length should be 2")

		// peek should return the plan
		gotModuleRun, err := queueManager.Peek(testutils.ModuleModels[0])

		assert.Nil(t, err)

		assert.Equal(t, moduleRun2.ID, gotModuleRun.ID, "module run ids are equal")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams, testutils.InitModules)
}

func TestWithLockIDs(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		queueManager := queuemanager.NewDefaultModuleRunQueueManager(conf.Repository)

		// create 5 module runs:
		// 1. Has lock-id-1, is a plan
		// 2. Has lock-id-2, is a plan
		// 3. Has lock-id-2, is an apply
		// 4. Has lock-id-1, is an apply
		// 5. Has no lock id, is an apply.
		//
		// Should be removed from queue in the following order: 5-1-4-(remove module lock)-2-3.
		//
		// This is because runs without a lock id are processed first.
		// Next, runs that were queued first are processed, with plans processed before apply's.
		basePlan := models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindPlan,
		}

		baseApply := models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindApply,
		}

		moduleRun1 := basePlan
		moduleRun2 := basePlan
		moduleRun3 := baseApply
		moduleRun4 := baseApply
		moduleRun5 := baseApply

		_moduleRun1, err := conf.Repository.Module().CreateModuleRun(&moduleRun1)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun1, &queuemanager.LockOpts{
			LockID:   "lock-id-1",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		_moduleRun2, err := conf.Repository.Module().CreateModuleRun(&moduleRun2)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun2, &queuemanager.LockOpts{
			LockID:   "lock-id-2",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		_moduleRun3, err := conf.Repository.Module().CreateModuleRun(&moduleRun3)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun3, &queuemanager.LockOpts{
			LockID:   "lock-id-2",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		_moduleRun4, err := conf.Repository.Module().CreateModuleRun(&moduleRun4)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun4, &queuemanager.LockOpts{
			LockID:   "lock-id-1",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		_moduleRun5, err := conf.Repository.Module().CreateModuleRun(&moduleRun5)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun5, &queuemanager.LockOpts{
			LockID: "",
		})

		assert.Nil(t, err)

		queueLen, err := queueManager.Len(testutils.ModuleModels[0])

		assert.Nil(t, err)
		assert.Equal(t, 5, queueLen, "queue length should be 5")

		gotRun, err := queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun5.ID, gotRun.ID, "first peek returns module run 5")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		gotRun, err = queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun1.ID, gotRun.ID, "second peek returns module run 1")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		gotRun, err = queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun4.ID, gotRun.ID, "third peek returns module run 4")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		gotRun, err = queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun2.ID, gotRun.ID, "fourth peek returns module run 4")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		gotRun, err = queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun3.ID, gotRun.ID, "fifth peek returns module run 3")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams, testutils.InitModules)
}

func TestWithMultiKind(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		queueManager := queuemanager.NewDefaultModuleRunQueueManager(conf.Repository)

		// create 3 module runs:
		// 1. Has lock-id-1, is a plan
		// 2. Has lock-id-1, is a plan
		// 3. Has lock-id-1, is an apply
		//
		// 1 should be removed from queue, and then 2 and 3 should be popped from the queue.
		basePlan := models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindPlan,
		}

		baseApply := models.ModuleRun{
			ModuleID: testutils.ModuleModels[0].ID,
			Status:   models.ModuleRunStatusQueued,
			Kind:     models.ModuleRunKindApply,
		}

		moduleRun1 := basePlan
		moduleRun2 := basePlan
		moduleRun3 := baseApply

		_moduleRun1, err := conf.Repository.Module().CreateModuleRun(&moduleRun1)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun1, &queuemanager.LockOpts{
			LockID:   "lock-id-1",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		_moduleRun2, err := conf.Repository.Module().CreateModuleRun(&moduleRun2)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun2, &queuemanager.LockOpts{
			LockID:   "lock-id-1",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		_moduleRun3, err := conf.Repository.Module().CreateModuleRun(&moduleRun3)

		err = queueManager.Enqueue(testutils.ModuleModels[0], _moduleRun3, &queuemanager.LockOpts{
			LockID:   "lock-id-1",
			LockKind: models.ModuleLockKindVCSBranch,
		})

		assert.Nil(t, err)

		queueLen, err := queueManager.Len(testutils.ModuleModels[0])

		assert.Nil(t, err)
		assert.Equal(t, 2, queueLen, "queue length should be 2")

		gotRun, err := queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun2.ID, gotRun.ID, "first peek returns module run 2")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		gotRun, err = queueManager.Peek(testutils.ModuleModels[0])
		assert.Equal(t, _moduleRun3.ID, gotRun.ID, "second peek returns module run 3")
		queueManager.Remove(testutils.ModuleModels[0], gotRun)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams, testutils.InitModules)
}
