package models

type ModuleRunQueue struct {
	Base

	ModuleID string

	Items []ModuleRunQueueItem
}

type ModuleQueuePriority uint

// Plans are run before applys. This is to handle the edge case where a plan and apply are queued
// at approximately the same time (ex. a forced merge on Github).
const (
	ModuleQueuePriorityPlan    ModuleQueuePriority = 3
	ModuleQueuePriorityDestroy ModuleQueuePriority = 2
	ModuleQueuePriorityApply   ModuleQueuePriority = 1
)

type LockPriority uint

const (
	NoLockID  LockPriority = 2
	HasLockID LockPriority = 1
)

type ModuleRunQueueItem struct {
	Base

	ModuleRunQueueID string
	ModuleRunID      string

	ModuleRunKind ModuleRunKind

	LockPriority LockPriority
	LockID       string
	LockKind     ModuleLockKind

	Priority ModuleQueuePriority
}
