package five9types

import "fmt"

type CampaignMode string

const (
	CampaignModeBasic    CampaignMode = "BASIC"
	CampaignModeAdvanced CampaignMode = "ADVANCED"
)

type CampaignStateLabel string

const (
	CampaignStateLabelRunning    CampaignStateLabel = "RUNNING"
	CampaignStateLabelNotRunning CampaignStateLabel = "NOT_RUNNING"
)

type UserLoginState string

const (
	UserLoginStateWorking       UserLoginState = "WORKING"
	UserLoginStateSelectStation UserLoginState = "SELECT_STATION"
	UserLoginStateAcceptNotice  UserLoginState = "ACCEPT_NOTICE"
)

type Channel string

const (
	ChannelVideo     Channel = "Video"
	ChannelTotal     Channel = "Total"
	ChannelChat      Channel = "Chat"
	ChannelVoicemail Channel = "Voicemail"
	ChannelVoice     Channel = "Voice"
)

type PresenceStateCode string

const (
	PresenceStateCodePendingState PresenceStateCode = "pendingState"
	PresenceStateCodeCurrentState PresenceStateCode = "currentState"
)

type DataSource string

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

type MaintenanceNoticeID string

type UserState string

type UserRole string

const (
	UserRoleDomainAdmin      UserRole = "DomainAdmin"
	UserRoleDomainSupervisor UserRole = "DomainSupervisor"
	UserRoleAgent            UserRole = "Agent"
	UserRoleReporting        UserRole = "Reporting"
)

type EventReason string

const (
	EventReasonConnectionSuccessful EventReason = "Successful WebSocket Connection"
	EventReasonUpdated              EventReason = "UPDATED"
)

type EventID string

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

type Policy string

const (
	PolicyAttachExisting Policy = "AttachExisting"
	PolicyForceIn        Policy = "ForceIn"
)

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionContext struct {
	CloudClientURL string `json:"cloudClientUrl"`
	CloudTokenURL  string `json:"cloudTokenUrl"`
	FarmID         FarmID `json:"farmId"`
}

type SessionMetadata struct {
	FreedomURL  string       `json:"freedomUrl"`
	DataCenters []DataCenter `json:"dataCenters"`
}

type DataCenter struct {
	Name   string   `json:"name"`
	UI     []server `json:"uiUrls"`
	API    []server `json:"apiUrls"`
	Login  []server `json:"loginUrls"`
	Active bool     `json:"active"`
}
type server struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	RouteKey string `json:"routeKey"`
	Version  string `json:"version"`
}

type (
	CampaignID     string
	CorrelationID  string
	FarmID         string
	MessageID      string
	OrganizationID string
	ProfileID      string
	SessionID      string
	StationID      string
	TenantID       string
	UserID         string
	UserName       string
)

type WebSocketIncrementalStatsUpdateData struct {
	DataSource DataSource   `json:"dataSource"`
	Added      []AgentState `json:"added"`
	Updated    []AgentState `json:"updated"`
	Removed    []AgentState `json:"removed"`
}

type WebSocketStatisticsAgentStateData struct {
	ID                         UserID                   `json:"id"`
	CallType                   any                      `json:"callType"`
	CampaignID                 *CampaignID              `json:"campaignId"`
	Customer                   any                      `json:"customer"`
	MediaAvailability          string                   `json:"mediaAvailability"`
	ParkedCallsCount           uint64                   `json:"parkedCallsCount"`
	ReasonCodeID               string                   `json:"reasonCodeId"`
	State                      UserState                `json:"state"`
	StateSince                 uint64                   `json:"stateSince"`
	StateDuration              uint64                   `json:"stateDuration"`
	OnHoldStateSince           uint64                   `json:"onHoldStateSince"`
	OnHoldStateDuration        uint64                   `json:"onHoldStateDuration"`
	OnParkStateSince           uint64                   `json:"onParkStateSince"`
	OnParkStateDuration        uint64                   `json:"onParkStateDuration"`
	ReasonCodeSince            uint64                   `json:"reasonCodeSince"`
	ReasonCodeDuration         uint64                   `json:"reasonCodeDuration"`
	AfterCallWorkStateSince    uint64                   `json:"afterCallWorkStateSince"`
	AfterCallWorkStateDuration uint64                   `json:"afterCallWorkStateDuration"`
	LoggedOutStateSince        uint64                   `json:"loggedOutStateSince"`
	LoggedOutStateDuration     uint64                   `json:"loggedOutStateDuration"`
	NotReadyStateSince         uint64                   `json:"notReadyStateSince"`
	NotReadyStateDuration      uint64                   `json:"notReadyStateDuration"`
	OnCallStateSince           uint64                   `json:"onCallStateSince"`
	OnCallStateDuration        uint64                   `json:"onCallStateDuration"`
	ReadyStateSince            uint64                   `json:"readyStateSince"`
	ReadyStateDuration         uint64                   `json:"readyStateDuration"`
	PermanentRecording         bool                     `json:"permanentRecording"`
	SessionRecording           bool                     `json:"sessionRecording"`
	ReadyChannels              string                   `json:"readyChannels"`
	NotReadyReasonCode         uint64                   `json:"notReadyReasonCode"`
	ChannelAvailability        map[Channel]channelState `json:"channelAvailability"`
}

type WebSocketStatisticsUserSessionData struct {
	ID           SessionID `json:"id"`
	UserID       UserID    `json:"userId"`
	UserName     UserName  `json:"userName"`
	FullName     string    `json:"fullName"`
	Role         UserRole  `json:"role"`
	SessionStart uint64    `json:"sessionStart"`
	Station      StationID `json:"station"`
}

type WebSocketStatisticsACDData struct {
	ID                      string `json:"id"`
	CallsInQueue            uint64 `json:"callsInQueue"`
	CallbacksInQueue        uint64 `json:"callbacksInQueue"`
	VoicemailsInQueue       uint64 `json:"voicemailsInQueue"`
	VoicemailsInProgress    uint64 `json:"voicemailsInProgress"`
	VoicemailsTotal         uint64 `json:"voicemailsTotal"`
	AgentsInVoicemailQueue  uint64 `json:"agentsInVoicemailQueue"`
	AgentsActive            uint64 `json:"agentsActive"`
	AgentsLoggedIn          uint64 `json:"agentsLoggedIn"`
	AgentsInQueue           uint64 `json:"agentsInQueue"`
	AgentsOnCall            uint64 `json:"agentsOnCall"`
	AgentsNotReadyForCalls  uint64 `json:"agentsNotReadyForCalls"`
	LongestQueueTime        uint64 `json:"longestQueueTime"`
	CurrentLongestQueueTime uint64 `json:"currentLongestQueueTime"`
	VivrCallsInQueue        uint64 `json:"vivrCallsInQueue"`
}

type WebSocketStatisticsInboundCampaignStatisticsData struct {
	ID                               CampaignID        `json:"id"`
	AssociatedWithAgentsCallsCount   uint64            `json:"associatedWithAgentsCallsCount"`
	ConnectedPlusAbandonedCallsCount uint64            `json:"connectedPlusAbandonedCallsCount"`
	AbandonCallRate                  float64           `json:"abandonCallRate"`
	TotalCallsCount                  uint64            `json:"totalCallsCount"`
	AverageAvailabilityTime          uint64            `json:"averageAvailabilityTime"`
	AverageCallTime                  uint64            `json:"averageCallTime"`
	HandledCallsCount                uint64            `json:"handledCallsCount"`
	AverageHandleTime                uint64            `json:"averageHandleTime"`
	AverageSpeedOfAnswer             uint64            `json:"averageSpeedOfAnswer"`
	AverageWrapTime                  uint64            `json:"averageWrapTime"`
	CallCharges                      float64           `json:"callCharges"`
	AbandonedCallsCount              uint64            `json:"abandonedCallsCount"`
	ConnectedCallsCount              uint64            `json:"connectedCallsCount"`
	FinishedInIVRErrorCallsCount     uint64            `json:"finishedInIVRErrorCallsCount"`
	FinishedInIVRSuccessCallsCount   uint64            `json:"finishedInIVRSuccessCallsCount"`
	RejectedCallsCount               uint64            `json:"rejectedCallsCount"`
	Dispositions                     map[string]uint64 `json:"dispositions"`
	FirstCallResolution              uint64            `json:"firstCallResolution"`
	LongestHoldTime                  uint64            `json:"longestHoldTime"`
	LongestQueueTime                 uint64            `json:"longestQueueTime"`
	ServiceLevelQueue                float64           `json:"serviceLevelQueue"`
	ServiceLevelTalk                 float64           `json:"serviceLevelTalk"`
	VivrSessionsCounty               uint64            `json:"vivrSessionsCounty"`
}

type WebSocketStatisticsAgentStatisticsData struct {
	ID                              UserID            `json:"id"`
	TotalCallsCount                 uint64            `json:"totalCallsCount"`
	AgentCallsCount                 uint64            `json:"agentCallsCount"`
	TotalCallsWithoutInternalsCount uint64            `json:"totalCallsWithoutInternalsCount"`
	BreaksCount                     uint64            `json:"breaksCount"`
	AverageBreakTime                uint64            `json:"averageBreakTime"`
	AverageCallTime                 uint64            `json:"averageCallTime"`
	AverageHoldTime                 uint64            `json:"averageHoldTime"`
	AverageIdleTime                 uint64            `json:"averageIdleTime"`
	InternalCallsCount              uint64            `json:"internalCallsCount"`
	AverageInternalCallTime         uint64            `json:"averageInternalCallTime"`
	PreviewCallsCount               uint64            `json:"previewCallsCount"`
	PreviewTime                     uint64            `json:"previewTime"`
	AveragePreviewTime              uint64            `json:"averagePreviewTime"`
	AverageHandleTime               uint64            `json:"averageHandleTime"`
	ProcessedVoicemailCount         uint64            `json:"processedVoicemailCount"`
	AverageVoicemailProcessingTime  uint64            `json:"averageVoicemailProcessingTime"`
	AverageVoicemailReadyTime       uint64            `json:"averageVoicemailReadyTime"`
	AverageWrapTime                 uint64            `json:"averageWrapTime"`
	CallCharges                     float64           `json:"callCharges"`
	SkippedInPreviewCallsCount      uint64            `json:"skippedInPreviewCallsCount"`
	Dispositions                    map[string]uint64 `json:"dispositions"`
	FirstCallResolution             uint64            `json:"firstCallResolution"`
	InboundCallsCount               uint64            `json:"inboundCallsCount"`
	SuccessfulInternalCallsCount    uint64            `json:"successfulInternalCallsCount"`
	LoginTime                       uint64            `json:"loginTime"`
	Occupancy                       float64           `json:"occupancy"`
	OutboundCallsCount              uint64            `json:"outboundCallsCount"`
	OffBreakTime                    uint64            `json:"offBreakTime"`
	Utilization                     float64           `json:"utilization"`
}

type WebSocketStatisticsOutboundCampaignManagerData struct {
	ID                                CampaignID         `json:"id"`
	ReadyForCallAgentsCount           uint64             `json:"readyForCallAgentsCount"`
	DispositionedRecordsCount         uint64             `json:"dispositionedRecordsCount"`
	DialingAttemptsCount              uint64             `json:"dialingAttemptsCount"`
	ContactedCallsCount               uint64             `json:"contactedCallsCount"`
	SkippedInPreviewCallsCount        uint64             `json:"skippedInPreviewCallsCount"`
	CallsToAgentRatio                 float64            `json:"callsToAgentRatio"`
	CallsToAgentTargetRatio           float64            `json:"callsToAgentTargetRatio"`
	TotalRecordsCount                 uint64             `json:"totalRecordsCount"`
	AvailableRecordsCount             uint64             `json:"availableRecordsCount"`
	RedialedWithTimerRecordsCount     uint64             `json:"redialedWithTimerRecordsCount"`
	DialedWithoutTimerRecordsCount    uint64             `json:"dialedWithoutTimerRecordsCount"`
	DialedWithASAPRequestRecordsCount uint64             `json:"dialedWithASAPRequestRecordsCount"`
	NoPartyContactSystemCallsCount    uint64             `json:"noPartyContactSystemCallsCount"`
	AbandonedCallsCount               uint64             `json:"abandonedCallsCount"`
	UnreachableRecordsCount           uint64             `json:"unreachableRecordsCount"`
	CampaignState                     CampaignStateLabel `json:"campaignState"`
}

type WebSocketStatisticsCampaignStateData struct {
	ID            CampaignID         `json:"id"`
	CampaignState CampaignStateLabel `json:"campaignState"`
	Priority      *uint64            `json:"priority"`
	Ratio         *uint64            `json:"ratio"`
	CurrentAction string             `json:"currentAction"`
	StateSince    uint64             `json:"stateSince"`
	Mode          *CampaignMode      `json:"mode"`
	ProfileID     *ProfileID         `json:"profileId"`
}

type channelState struct {
	Current uint64 `json:"current"`
	Max     uint64 `json:"max"`
	Status  string `json:"status"`
}

type Presence struct {
	OnVoice                  bool   `json:"onVoice"`
	OnSCC                    bool   `json:"onSCC"`
	ChangeTimestamp          uint64 `json:"changeTimestamp"`
	NextStateChangeTimestamp uint64 `json:"nextStateChangeTimestamp"`
	GracefulModeOn           bool   `json:"gracefulModeOn"`
}

type StationInfo struct {
	StationID   string `json:"stationId"`
	StationType string `json:"stationType"`
}

type SupervisorUserInfo struct {
	Email    string   `json:"email"`
	ID       UserID   `json:"id"`
	UserName UserName `json:"userName"`
}

type (
	ReasonCodeID          string
	AuthenticationTokenID string
)

// type (
//
//	UserName              string
//	CampaignID            string
//	UserID                string
//	UserState             string
//	ReasonCodeID          string
//	NotReadyReasonCode    int64

// 	farmID                string
// 	organizationID        string
// 	policy                string
// 	sessionID             string
// 	eventID               string
// 	userLoginState        string
// 	dataSource            string
// )

const (
	UserStateAfterCallWork UserState = "ACW"
	UserStateLoggedOut     UserState = "LOGGED_OUT"
	UserStateNotReady      UserState = "NOT_READY"
	UserStateReady         UserState = "READY"
	UserStateOnCall        UserState = "ON_CALL"
	UserStateRinging       UserState = "RINGING"
)

// const (
// 	eventIDServerConnected                    eventID = "1010"
// 	eventIDDuplicateConnection                eventID = "1020"
// 	eventIDPongReceived                       eventID = "1202" // Pong response to ping request
// 	eventIDSupervisorStats                    eventID = "5000" // Statistics data has been received.
// 	eventIDDispositionsInvalidated            eventID = "5002" // Disposition has been removed or	created, or disposition name has been changed.
// 	eventIDSkillsInvalidated                  eventID = "5003" // Skill has been removed or created, or	queue name has been changed.
// 	eventIDAgentGroupsInvalidated             eventID = "5004"
// 	eventIDCampaignsInvalidated               eventID = "5005"
// 	eventIDUsersInvalidated                   eventID = "5006"
// 	eventIDReasonCodesInvalidated             eventID = "5007"
// 	eventIDCampaignProfilesInvalidated        eventID = "5008"
// 	eventIDCampaignOutOfNumbers               eventID = "5009"
// 	eventIDListsInvalidated                   eventID = "5010"
// 	eventIDCampaignListsChanged               eventID = "5011"
// 	eventIDIncrementalStatsUpdate             eventID = "5012"
// 	eventIDIncrementalUserProfilesUpdate      eventID = "5013"
// 	eventIDFilterSettingsUpdated              eventID = "6001"
// 	eventIDAgentsInvalidated                  eventID = "6002"
// 	eventIDPermissionsUpdated                 eventID = "6003"
// 	eventIDResetCampaignDispositionsCompleted eventID = "6004"
// 	eventIDMonitoringStateUpdated             eventID = "6005"
// 	eventIDRandomMonitoringStarted            eventID = "6006"
// 	eventIDFdsRealTime                        eventID = "6007"
// 	eventIDIncrementalInteractions            eventID = "6008"
// )

// const (
// 	PolicyAttachExisting policy = "AttachExisting"
// 	PolicyForceIn        policy = "ForceIn"
// )

// const (
// 	dataSourceACDStatus                  dataSource = "ACD_STATUS"
// 	dataSourceAgentState                 dataSource = "AGENT_STATE"
// 	dataSourceAgentStatistic             dataSource = "AGENT_STATISTIC"
// 	dataSourceCampaignState              dataSource = "CAMPAIGN_STATE"
// 	dataSourceInboundCampaignStatistics  dataSource = "INBOUND_CAMPAIGN_STATISTICS"
// 	dataSourceStations                   dataSource = "STATIONS"
// 	dataSourceOutboundCampaignStatistics dataSource = "OUTBOUND_CAMPAIGN_STATISTICS"
// 	dataSourceOutboundCampaignManager    dataSource = "OUTBOUND_CAMPAIGN_MANAGER"
// 	dataSourceUserSession                dataSource = "USER_SESSION"
// )

// type webSocketIncrementalStatsUpdateData struct {
// 	DataSource dataSource   `json:"dataSource"`
// 	Added      []AgentState `json:"added"`
// 	Updated    []AgentState `json:"updated"`
// 	Removed    []AgentState `json:"removed"`
// }

type ReasonCodeInfo struct {
	ID         ReasonCodeID `json:"id"`
	Name       string       `json:"name"`
	Selectable bool         `json:"selectable"`
}

// type domainMetadata struct {
// 	reasonCodes map[ReasonCodeID]ReasonCodeInfo
// 	agentInfo   map[UserID]AgentInfo
// }

type MaintenanceNoticeInfo struct {
	Accepted   bool                `json:"accepted"`
	Annotation string              `json:"annotation"`
	ID         MaintenanceNoticeID `json:"id"`
	Text       string              `json:"text"`
}

type WebsocketSupervisorStatsData struct {
	Data []AgentState `json:"data"`
}

type AgentState struct {
	ID                         UserID       `json:"id"`
	CallType                   any          `json:"callType"`
	CampaignID                 *CampaignID  `json:"campaignId"`
	Customer                   any          `json:"customer"`
	MediaAvailability          string       `json:"mediaAvailability"`
	ParkedCallsCount           uint64       `json:"parkedCallsCount"`
	ReasonCodeID               ReasonCodeID `json:"reasonCodeId"`
	State                      UserState    `json:"state"`
	StateSince                 uint64       `json:"stateSince"`
	StateDuration              uint64       `json:"stateDuration"`
	OnHoldStateSince           uint64       `json:"onHoldStateSince"`
	OnHoldStateDuration        uint64       `json:"onHoldStateDuration"`
	OnParkStateSince           uint64       `json:"onParkStateSince"`
	OnParkStateDuration        uint64       `json:"onParkStateDuration"`
	ReasonCodeSince            uint64       `json:"reasonCodeSince"`
	ReasonCodeDuration         uint64       `json:"reasonCodeDuration"`
	AfterCallWorkStateSince    uint64       `json:"afterCallWorkStateSince"`
	AfterCallWorkStateDuration uint64       `json:"afterCallWorkStateDuration"`
	LoggedOutStateSince        uint64       `json:"loggedOutStateSince"`
	LoggedOutStateDuration     uint64       `json:"loggedOutStateDuration"`
	NotReadyStateSince         uint64       `json:"notReadyStateSince"`
	NotReadyStateDuration      uint64       `json:"notReadyStateDuration"`
	OnCallStateSince           uint64       `json:"onCallStateSince"`
	OnCallStateDuration        uint64       `json:"onCallStateDuration"`
	ReadyStateSince            uint64       `json:"readyStateSince"`
	ReadyStateDuration         uint64       `json:"readyStateDuration"`
	PermanentRecording         bool         `json:"permanentRecording"`
	SessionRecording           bool         `json:"sessionRecording"`
	ReadyChannels              string       `json:"readyChannels"`
	// ChannelAvailability        map[Channel]channelState `json:"channelAvailability"`
}

// type StationInfo struct {
// 	StationID   string `json:"stationId"`
// 	StationType string `json:"stationType"`
// }

type AgentInfo struct {
	ID       UserID   `json:"id"`
	UserName UserName `json:"userName"`
}

type LoginResponse struct {
	TokenID   AuthenticationTokenID `json:"tokenID"`
	SessionID SessionID             `json:"sessionId"`
	OrgID     OrganizationID        `json:"orgID"`
	UserID    UserID                `json:"userID"`
	Context   SessionContext        `json:"context"`
	Metadata  SessionMetadata       `json:"metadata"`
}

func (v LoginResponse) GetAPIHost() string {
	for _, dataCenter := range v.Metadata.DataCenters {
		if !dataCenter.Active {
			continue
		}

		for _, server := range dataCenter.API {
			return fmt.Sprintf("%s:%s", server.Host, server.Port)
		}
	}

	return "app.five9.com:443"
}

type LoginPayload struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
	AppKey              string              `json:"appKey"`
	Policy              Policy              `json:"policy"`
}

// type sessionContext struct {
// 	CloudClientURL string `json:"cloudClientUrl"`
// 	CloudTokenURL  string `json:"cloudTokenUrl"`
// 	FarmID         farmID `json:"farmId"`
// }

// type sessionMetadata struct {
// 	FreedomURL  string       `json:"freedomUrl"`
// 	DataCenters []dataCenter `json:"dataCenters"`
// }

// type server struct {
// 	Host     string `json:"host"`
// 	Port     string `json:"port"`
// 	RouteKey string `json:"routeKey"`
// 	Version  string `json:"version"`
// }
