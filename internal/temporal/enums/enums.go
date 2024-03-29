package enums

const (
	BackgroundQueueName string = "background"
	ModuleRunQueueName  string = "runner"
)

const (
	BackgroundLogFlushID string = "log_flusher"
	BackgroundNotifierID string = "notifier"
	// TODO: module run scheduler should in background queue?
	BackgroundQueueCheckerID       string = "queue_checker"
	BackgroundModuleRunSchedulerID string = "module_run_scheduler"
)

const (
	WorkflowTypeNameLogFlush string = "FlushLogs"
	WorkflowTypeNameNotifier string = "NotifyWorkflow"
	// TODO: consolidate language here
	WorkflowTypeNameCheckModuleQueue string = "ScheduleFromQueue"
	WorkflowTypeNameCheckAllQueues   string = "CheckQueues"
	WorkflowTypeNameProvision        string = "Provision"
	WorkflowTypeNameDispatchMonitors string = "DispatchMonitors"
	WorkflowTypeNameRunMonitor       string = "RunMonitor"
)
