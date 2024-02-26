package five9types

type (
	AuthenticationTokenID       string
	CampaignID                  string
	CampaignMode                string
	CampaignStateLabel          string
	Channel                     string
	CorrelationID               string
	DataSource                  string
	EventID                     string
	EventReason                 string
	FarmID                      string
	FilterSettingsSelectionType string
	MaintenanceNoticeID         string
	MessageID                   string
	NotReadyReasonCode          int
	OrganizationID              string
	Policy                      string
	PresenceStateCode           string
	ProfileID                   string
	QueueID                     string
	ReasonCodeID                string
	SessionID                   string
	SkillID                     string
	StationID                   string
	StatisticsRange             string // Enumeration of time periods to use as statistics filters.
	StatisticsRollingPeriod     string // Enumeration with the time period for the list and campaign statistics.
	TenantID                    string
	UserID                      string
	UserLoginState              string
	UserName                    string
	UserRole                    string
	UserState                   string
)

const (
	CampaignModeBasic    CampaignMode = "BASIC"
	CampaignModeAdvanced CampaignMode = "ADVANCED"
)

const (
	CampaignStateLabelRunning    CampaignStateLabel = "RUNNING"
	CampaignStateLabelNotRunning CampaignStateLabel = "NOT_RUNNING"
)

const (
	UserLoginStateWorking       UserLoginState = "WORKING"
	UserLoginStateSelectStation UserLoginState = "SELECT_STATION"
	UserLoginStateAcceptNotice  UserLoginState = "ACCEPT_NOTICE"
	UserLoginStateRelogin       UserLoginState = "RELOGIN"
)

const (
	ChannelVideo     Channel = "Video"
	ChannelTotal     Channel = "Total"
	ChannelChat      Channel = "Chat"
	ChannelVoicemail Channel = "Voicemail"
	ChannelVoice     Channel = "Voice"
)

const (
	PresenceStateCodePendingState PresenceStateCode = "pendingState"
	PresenceStateCodeCurrentState PresenceStateCode = "currentState"
)

const (
	DataSourceACDStatus                  DataSource = "ACD_STATUS"
	DataSourceAgentState                 DataSource = "AGENT_STATE"
	DataSourceAgentStatistic             DataSource = "AGENT_STATISTIC"
	DataSourceCampaignState              DataSource = "CAMPAIGN_STATE"
	DataSourceInboundCampaignStatistics  DataSource = "INBOUND_CAMPAIGN_STATISTICS"
	DataSourceStations                   DataSource = "STATIONS"
	DataSourceOutboundCampaignStatistics DataSource = "OUTBOUND_CAMPAIGN_STATISTICS"
	DataSourceOutboundCampaignManager    DataSource = "OUTBOUND_CAMPAIGN_MANAGER"
	DataSourceUserSession                DataSource = "USER_SESSION"
)

const (
	UserRoleDomainAdmin      UserRole = "DomainAdmin"
	UserRoleDomainSupervisor UserRole = "DomainSupervisor"
	UserRoleAgent            UserRole = "Agent"
	UserRoleReporting        UserRole = "Reporting"
)

const (
	EventReasonConnectionSuccessful EventReason = "Successful WebSocket Connection"
	EventReasonUpdated              EventReason = "UPDATED"
)

const (
	EventIDServerConnected                    EventID = "1010"
	EventIDDuplicateConnection                EventID = "1020"
	EventIDPongReceived                       EventID = "1202" // Pong response to ping request
	EventIDSupervisorStats                    EventID = "5000" // Statistics data has been received.
	EventIDDispositionsInvalidated            EventID = "5002" // Disposition has been removed or	created, or disposition name has been changed.
	EventIDSkillsInvalidated                  EventID = "5003" // Skill has been removed or created, or	queue name has been changed.
	EventIDAgentGroupsInvalidated             EventID = "5004"
	EventIDCampaignsInvalidated               EventID = "5005"
	EventIDUsersInvalidated                   EventID = "5006"
	EventIDReasonCodesInvalidated             EventID = "5007"
	EventIDCampaignProfilesInvalidated        EventID = "5008"
	EventIDCampaignOutOfNumbers               EventID = "5009"
	EventIDListsInvalidated                   EventID = "5010"
	EventIDCampaignListsChanged               EventID = "5011"
	EventIDIncrementalStatsUpdate             EventID = "5012"
	EventIDIncrementalUserProfilesUpdate      EventID = "5013"
	EventIDFilterSettingsUpdated              EventID = "6001"
	EventIDAgentsInvalidated                  EventID = "6002"
	EventIDPermissionsUpdated                 EventID = "6003"
	EventIDResetCampaignDispositionsCompleted EventID = "6004"
	EventIDMonitoringStateUpdated             EventID = "6005"
	EventIDRandomMonitoringStarted            EventID = "6006"
	EventIDFdsRealTime                        EventID = "6007"
	EventIDIncrementalInteractions            EventID = "6008"
)

const (
	PolicyAttachExisting Policy = "AttachExisting"
	PolicyForceIn        Policy = "ForceIn"
)

const (
	SelectAll     FilterSettingsSelectionType = "ALL"     // Show statistics for all the queues or agent groups
	SelectMy      FilterSettingsSelectionType = "MY"      // Show statistics for the user's queues or agent groups
	SelectSpecify FilterSettingsSelectionType = "SPECIFY" // Show statistics for the specified queue or agent groups.

)

const (
	RangeCurrentDay   StatisticsRange = "CURRENT_DAY"   // Current day
	RangeCurrentMonth StatisticsRange = "CURRENT_MONTH" // Current month
	RangeCurrentShift StatisticsRange = "CURRENT_SHIFT" // Current day shift by the number of hours specified in StatsFilterSettingsInfo
	RangeCurrentWeek  StatisticsRange = "CURRENT_WEEK"  // Current week
	RangeLifetime     StatisticsRange = "LIFETIME"      // All time
	RangeRollingHour  StatisticsRange = "ROLLING_HOUR"  // Last hour divided into five minute intervals
)

const (
	FiveMinutes    StatisticsRollingPeriod = "MINUTES5"
	TenMinutes     StatisticsRollingPeriod = "MINUTES10"
	FifteenMinutes StatisticsRollingPeriod = "MINUTES15"
	ThirtyMinutes  StatisticsRollingPeriod = "MINUTES30"
	OneHour        StatisticsRollingPeriod = "HOUR1"
	TwoHours       StatisticsRollingPeriod = "HOUR2"
	ThreeHours     StatisticsRollingPeriod = "HOUR3"
	Today          StatisticsRollingPeriod = "TODAY"
)

const (
	UserStateAfterCallWork UserState = "ACW"
	UserStateLoggedOut     UserState = "LOGGED_OUT"
	UserStateNotReady      UserState = "NOT_READY"
	UserStateReady         UserState = "READY"
	UserStateOnCall        UserState = "ON_CALL"
	UserStateRinging       UserState = "RINGING"
)
