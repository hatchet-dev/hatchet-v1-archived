package models

type ModuleRunQueue struct {
	Base

	ModuleID string

	Items []ModuleRunQueueItem
}

type ModuleQueuePriority uint

const (
	ModuleQueuePriorityPlan    ModuleQueuePriority = 1
	ModuleQueuePriorityApply   ModuleQueuePriority = 2
	ModuleQueuePriorityDestroy ModuleQueuePriority = 3
)

type ModuleRunQueueItem struct {
	Base

	ModuleRunQueueID string
	ModuleRunID      string

	Priority ModuleQueuePriority
}
