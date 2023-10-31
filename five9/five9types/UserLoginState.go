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

type UserRole string

type SkillID string

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
	UI     []Server `json:"uiUrls"`
	API    []Server `json:"apiUrls"`
	Login  []Server `json:"loginUrls"`
	Active bool     `json:"active"`
}
type Server struct {
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

type WebSocketIncrementalAgentStateData struct {
	DataSource DataSource   `json:"dataSource"`
	Added      []AgentState `json:"added"`
	Updated    []AgentState `json:"updated"`
	Removed    []AgentState `json:"removed"`
}

type WebSocketIncrementalAgentStatisticsData struct {
	DataSource DataSource        `json:"dataSource"`
	Added      []AgentStatistics `json:"added"`
	Updated    []AgentStatistics `json:"updated"`
	Removed    []AgentStatistics `json:"removed"`
}

type WebSocketIncrementalACDStateData struct {
	DataSource DataSource `json:"dataSource"`
	Added      []ACDState `json:"added"`
	Updated    []ACDState `json:"updated"`
	Removed    []ACDState `json:"removed"`
}

type WebSocketStatisticsAgentStateData struct {
	ID                         UserID                   `json:"id"`
	CallType                   any                      `json:"callType"`
	CampaignID                 *CampaignID              `json:"campaignId"`
	Customer                   any                      `json:"customer"`
	MediaAvailability          string                   `json:"mediaAvailability"`
	ParkedCallsCount           uint64                   `json:"parkedCallsCount"`
	ReasonCodeID               ReasonCodeID             `json:"reasonCodeId"`
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
	NotReadyReasonCode         NotReadyReasonCode       `json:"notReadyReasonCode"`
	ChannelAvailability        map[Channel]ChannelState `json:"channelAvailability"`
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

type ACDState struct {
	ID                      QueueID `json:"id"`
	CallsInQueue            uint64  `json:"callsInQueue"`
	CallbacksInQueue        uint64  `json:"callbacksInQueue"`
	VoicemailsInQueue       uint64  `json:"voicemailsInQueue"`
	VoicemailsInProgress    uint64  `json:"voicemailsInProgress"`
	VoicemailsTotal         uint64  `json:"voicemailsTotal"`
	AgentsInVoicemailQueue  uint64  `json:"agentsInVoicemailQueue"`
	AgentsActive            uint64  `json:"agentsActive"`
	AgentsLoggedIn          uint64  `json:"agentsLoggedIn"`
	AgentsInQueue           uint64  `json:"agentsInQueue"`
	AgentsOnCall            uint64  `json:"agentsOnCall"`
	AgentsNotReadyForCalls  uint64  `json:"agentsNotReadyForCalls"`
	LongestQueueTime        uint64  `json:"longestQueueTime"`
	CurrentLongestQueueTime uint64  `json:"currentLongestQueueTime"`
	VivrCallsInQueue        uint64  `json:"vivrCallsInQueue"`
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

type AgentStatistics struct {
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

type ChannelState struct {
	Current uint64 `json:"current"`
	Max     uint64 `json:"max"`
	Status  string `json:"status"`
}

type ChannelAvailability struct {
	Video     ChannelState `json:"Video"`
	Total     ChannelState `json:"Total"`
	Chat      ChannelState `json:"Chat"`
	Voicemail ChannelState `json:"Voicemail"`
	Voice     ChannelState `json:"Voice"`
}

type Presence struct {
	OnVoice                  bool   `json:"onVoice"`
	OnSCC                    bool   `json:"onSCC"`
	ChangeTimestamp          uint64 `json:"changeTimestamp"`
	NextStateChangeTimestamp uint64 `json:"nextStateChangeTimestamp"`
	GracefulModeOn           bool   `json:"gracefulModeOn"`
	CurrentState             State  `json:"currentState"`
	PendingState             State  `json:"pendingState"`
}

type State struct {
	ReadyChannels      []string           `json:"readyChannels"`
	NotReadyReasonCode NotReadyReasonCode `json:"notReadyReasonCode"`
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
	NotReadyReasonCode    int
	ReasonCodeID          string
	AuthenticationTokenID string
)

type UserState string

const (
	UserStateAfterCallWork UserState = "ACW"
	UserStateLoggedOut     UserState = "LOGGED_OUT"
	UserStateNotReady      UserState = "NOT_READY"
	UserStateReady         UserState = "READY"
	UserStateOnCall        UserState = "ON_CALL"
	UserStateRinging       UserState = "RINGING"
)

type ReasonCodeInfo struct {
	ID         ReasonCodeID `json:"id"`
	Name       string       `json:"name"`
	Selectable bool         `json:"selectable"`
}

type MaintenanceNoticeInfo struct {
	Accepted   bool                `json:"accepted"`
	Annotation string              `json:"annotation"`
	ID         MaintenanceNoticeID `json:"id"`
	Text       string              `json:"text"`
}

type WebsocketSupervisorStateData struct {
	Data []AgentState `json:"data"`
}

type WebsocketSupervisorStatisticsData struct {
	Data []AgentStatistics `json:"data"`
}

type WebsocketSupervisorACDData struct {
	Data []ACDState `json:"data"`
}

type AgentState struct {
	ID                         UserID                   `json:"id"`
	CallType                   any                      `json:"callType"`
	CampaignID                 *CampaignID              `json:"campaignId"`
	Customer                   any                      `json:"customer"`
	MediaAvailability          string                   `json:"mediaAvailability"`
	ParkedCallsCount           uint64                   `json:"parkedCallsCount"`
	ReasonCodeID               ReasonCodeID             `json:"reasonCodeId"`
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
	Presence                   Presence                 `json:"presence"`
	ChannelAvailability        map[Channel]ChannelState `json:"channelAvailability"`
}

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
