package five9

// type CampaignMode string

// const (
// 	CampaignModeBasic    CampaignMode = "BASIC"
// 	CampaignModeAdvanced CampaignMode = "ADVANCED"
// )

// type CampaignStateLabel string

// const (
// 	CampaignStateLabelRunning    CampaignStateLabel = "RUNNING"
// 	CampaignStateLabelNotRunning CampaignStateLabel = "NOT_RUNNING"
// )

// type UserLoginState string

// const (
// 	UserLoginStateWorking UserLoginState = "WORKING"
// )

// type Channel string

// const (
// 	ChannelVideo     Channel = "Video"
// 	ChannelTotal     Channel = "Total"
// 	ChannelChat      Channel = "Chat"
// 	ChannelVoicemail Channel = "Voicemail"
// 	ChannelVoice     Channel = "Voice"
// )

// type PresenceStateCode string

// const (
// 	PresenceStateCodePendingState PresenceStateCode = "pendingState"
// 	PresenceStateCodeCurrentState PresenceStateCode = "currentState"
// )

// type DataSource string

// const (
// 	DataSourceACDStatus                  DataSource = "ACD_STATUS"
// 	DataSourceAgentState                 DataSource = "AGENT_STATE"
// 	DataSourceAgentStatistic             DataSource = "AGENT_STATISTIC"
// 	DataSourceCampaignState              DataSource = "CAMPAIGN_STATE"
// 	DataSourceInboundCampaignStatistics  DataSource = "INBOUND_CAMPAIGN_STATISTICS"
// 	DataSourceStations                   DataSource = "STATIONS"
// 	DataSourceOutboundCampaignStatistics DataSource = "OUTBOUND_CAMPAIGN_STATISTICS"
// 	DataSourceOutboundCampaignManager    DataSource = "OUTBOUND_CAMPAIGN_MANAGER"
// 	DataSourceUserSession                DataSource = "USER_SESSION"
// )

// type UserState string

// const (
// 	UserStateAfterCallWork UserState = "ACW"
// 	UserStateLoggedOut     UserState = "LOGGED_OUT"
// 	UserStateNotReady      UserState = "NOT_READY"
// 	UserStateReady         UserState = "READY"
// 	UserStateOnCall        UserState = "ON_CALL"
// )

// type UserRole string

// const (
// 	UserRoleDomainAdmin      UserRole = "DomainAdmin"
// 	UserRoleDomainSupervisor UserRole = "DomainSupervisor"
// 	UserRoleAgent            UserRole = "Agent"
// 	UserRoleReporting        UserRole = "Reporting"
// )

// type EventReason string

// const (
// 	EventReasonConnectionSuccessful EventReason = "Successful WebSocket Connection"
// 	EventReasonUpdated              EventReason = "UPDATED"
// )

// type EventID string

// const (
// 	EventIDServerConnected                    EventID = "1010"
// 	EventIDDuplicateConnection                EventID = "1020"
// 	EventIDPongReceived                       EventID = "1202" // Pong response to ping request
// 	EventIDSupervisorStats                    EventID = "5000" // Statistics data has been received.
// 	EventIDDispositionsInvalidated            EventID = "5002" // Disposition has been removed or	created, or disposition name has been changed.
// 	EventIDSkillsInvalidated                  EventID = "5003" // Skill has been removed or created, or	queue name has been changed.
// 	EventIDAgentGroupsInvalidated             EventID = "5004"
// 	EventIDCampaignsInvalidated               EventID = "5005"
// 	EventIDUsersInvalidated                   EventID = "5006"
// 	EventIDReasonCodesInvalidated             EventID = "5007"
// 	EventIDCampaignProfilesInvalidated        EventID = "5008"
// 	EventIDCampaignOutOfNumbers               EventID = "5009"
// 	EventIDListsInvalidated                   EventID = "5010"
// 	EventIDCampaignListsChanged               EventID = "5011"
// 	EventIDIncrementalStatsUpdate             EventID = "5012"
// 	EventIDIncrementalUserProfilesUpdate      EventID = "5013"
// 	EventIDFilterSettingsUpdated              EventID = "6001"
// 	EventIDAgentsInvalidated                  EventID = "6002"
// 	EventIDPermissionsUpdated                 EventID = "6003"
// 	EventIDResetCampaignDispositionsCompleted EventID = "6004"
// 	EventIDMonitoringStateUpdated             EventID = "6005"
// 	EventIDRandomMonitoringStarted            EventID = "6006"
// 	EventIDFdsRealTime                        EventID = "6007"
// 	EventIDIncrementalInteractions            EventID = "6008"
// )

// type Policy string

// const (
// 	PolicyAttachExisting Policy = "AttachExisting"
// 	PolicyForceIn        Policy = "ForceIn"
// )
