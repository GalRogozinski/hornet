package shutdown

const (
	PriorityCloseDatabase = iota
	PriorityFlushToDatabase
	PriorityRequestsProcessor
	PriorityMilestoneSolidifier
	PriorityMilestoneProcessor
	PrioritySolidifierGossip
	PriorityReceiveTxWorker
	PriorityBroadcastQueue
	PriorityMessageProcessor
	PriorityPeerSendQueue
	PriorityPeeringTCPServer
	PriorityPeerReconnecter
	PriorityHeartbeats
	PriorityWarpSync
	PriorityLocalSnapshots
	PriorityMetricsUpdater
	PriorityDashboard
	PriorityAPI
	PriorityMetricsPublishers
	PrioritySpammer
	PriorityStatusReport
	PriorityAutopeering
	PriorityCoordinator
	PriorityUpdateCheck
	PriorityPrometheus
	PriorityProfiler
)
