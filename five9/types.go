package five9

type LoginPayload struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
	AppKey              string              `json:"appKey"`
	Policy              Policy              `json:"policy"`
}

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Policy string

const (
	PolicyAttachExisting Policy = "AttachExisting"
	PolicyForceIn        Policy = "ForceIn"
)

type SupervisorMetadataResponse struct {
	TokenID   AuthenticationTokenID `json:"tokenID"`
	SessionID SessionID             `json:"sessionId"`
	OrgID     OrganizationID        `json:"orgID"`
	UserID    UserID                `json:"userID"`
	Context   SessionContext        `json:"context"`
	Metadata  SessionMetadata       `json:"metadata"`
}

type LoginAgentResponse struct {
	TokenID  AuthenticationTokenID `json:"tokenID"`
	OrgID    OrganizationID        `json:"orgID"`
	UserID   UserID                `json:"userID"`
	Context  SessionContext        `json:"context"`
	Metadata SessionMetadata       `json:"metadata"`
	Roles    []UserRole            `json:"roles"`
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
	Name      string   `json:"name"`
	UIURLs    []Server `json:"uiUrls"`
	APIURLs   []Server `json:"apiUrls"`
	LoginURLs []Server `json:"loginUrls"`
	Active    bool     `json:"active"`
}

type Server struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	RouteKey string `json:"routeKey"`
	Version  string `json:"version"`
}

type (
	AuthenticationTokenID string
	CampaignID            string
	CorrelationID         string
	FarmID                string
	MessageID             string
	OrganizationID        string
	ProfileID             string
	SessionID             string
	StationID             string
	TenantID              string
	UserID                string
	UserName              string
)

type UserState string

const (
	UserAfterCallWork UserState = "ACW"
	UserLoggedOut     UserState = "LOGGED_OUT"
	UserNotReady      UserState = "NOT_READY"
	UserReady         UserState = "READY"
	UserOnCall        UserState = "ON_CALL"
)

type UserRole string

const (
	RoleDomainAdmin      UserRole = "DomainAdmin"
	RoleDomainSupervisor UserRole = "DomainSupervisor"
	RoleAgent            UserRole = "Agent"
	RoleReporting        UserRole = "Reporting"
)

type EventReason string

const (
	WebSocketEventConnectionSuccessful EventReason = "Successful WebSocket Connection"
	WebSocketEventUpdated              EventReason = "UPDATED"
)

type EventID string

const (
	WebSocketServerConnected                    EventID = "1010"
	WebSocketDuplicateConnection                EventID = "1020"
	WebSocketSupervisorStats                    EventID = "5000" // Statistics data has been received.
	WebSocketDispositionsInvalidated            EventID = "5002" // Disposition has been removed or	created, or disposition name has been changed.
	WebSocketSkillsInvalidated                  EventID = "5003" // Skill has been removed or created, or	queue name has been changed.
	WebSocketAgentGroupsInvalidated             EventID = "5004"
	WebSocketCampaignsInvalidated               EventID = "5005"
	WebSocketUsersInvalidated                   EventID = "5006"
	WebSocketReasonCodesInvalidated             EventID = "5007"
	WebSocketCampaignProfilesInvalidated        EventID = "5008"
	WebSocketCampaignOutOfNumbers               EventID = "5009"
	WebSocketListsInvalidated                   EventID = "5010"
	WebSocketCampaignListsChanged               EventID = "5011"
	WebSocketIncremtalStatsUpdate               EventID = "5012"
	WebSocketIncrementalUserProfilesUpdate      EventID = "5013"
	WebSocketFilterSettingsUpdated              EventID = "6001"
	WebSocketAgentsInvalidated                  EventID = "6002"
	WebSocketPermissionsUpdated                 EventID = "6003"
	WebSocketResetCampaignDispositionsCompleted EventID = "6004"
	WebSocketMonitoringStateUpdated             EventID = "6005"
	WebSocketRandomMonitoringStarted            EventID = "6006"
	WebSocketFDSRealTime                        EventID = "6007"
	WebSocketIncrementalInteractions            EventID = "6008"
)

type LoginSupervisorResponse struct {
	TokenID   AuthenticationTokenID `json:"tokenID"`
	SessionID SessionID             `json:"sessionId"`
	OrgID     OrganizationID        `json:"orgID"`
	UserID    UserID                `json:"userID"`
	Context   SessionContext        `json:"context"`
	Metadata  SessionMetadata       `json:"metadata"`
}

type WebSocketEvent struct {
	Context WebSocketContext `json:"context"`
	PayLoad any              `json:"payLoad"`
}

type WebSocketContext struct {
	EventID       EventID        `json:"eventId"`
	EventReason   *EventReason   `json:"eventReason"`
	MessageID     *MessageID     `json:"messageId"`
	UserID        *UserID        `json:"userId"`
	CorrelationID *CorrelationID `json:"correlationId"`
	UserName      *UserName      `json:"userName"`
	Timestamp     uint64         `json:"timeStamp"`
	TenantID      TenantID       `json:"tenantId"`
	BroadCast     bool           `json:"broadCast"`
}

type WebSocketResponse interface {
	GetContext() *WebSocketContext
	GetPayload() any
}

type WebSocketStatsResponse struct {
	Context WebSocketContext     `json:"context"`
	PayLoad []WebSocketStatsData `json:"payLoad"`
}

func (w WebSocketStatsResponse) GetContext() WebSocketContext {
	return w.Context
}

func (w WebSocketStatsResponse) GetPayload() []WebSocketStatsData {
	return w.PayLoad
}

type WebSocketStatsData struct {
	DataSource DataSource `json:"dataSource"`
	Data       []any      `json:"data"`
}

type WebSocketIncrementalStatsUpdateResponse struct {
	Context WebSocketContext                       `json:"context"`
	PayLoad []WebSocketIncremtalStatsUpdatePayload `json:"payLoad"`
}

func (w WebSocketIncrementalStatsUpdateResponse) GetContext() WebSocketContext {
	return w.Context
}

func (w WebSocketIncrementalStatsUpdateResponse) GetPayload() []WebSocketIncremtalStatsUpdatePayload {
	return w.PayLoad
}

type WebSocketIncremtalStatsUpdatePayload struct {
	WebSocketIncrementalStatsUpdateData
}

type DataSource string

const (
	ACDStatus                  DataSource = "ACD_STATUS"
	AgentState                 DataSource = "AGENT_STATE"
	AgentStatistic             DataSource = "AGENT_STATISTIC"
	CampaignState              DataSource = "CAMPAIGN_STATE"
	InboundCampaignStatistics  DataSource = "INBOUND_CAMPAIGN_STATISTICS"
	Stations                   DataSource = "STATIONS"
	OutboundCampaignStatistics DataSource = "OUTBOUND_CAMPAIGN_STATISTICS"
	OutboundCampaignManager    DataSource = "OUTBOUND_CAMPAIGN_MANAGER"
	UserSession                DataSource = "USER_SESSION"
)

type WebSocketIncrementalStatsUpdateData struct {
	DataSource DataSource `json:"dataSource"`
	Added      []any      `json:"added"`
	Updated    []any      `json:"updated"`
	Removed    []any      `json:"removed"`
}

type WebSocketStatisticsAgentStateData struct {
	ID                         UserID                              `json:"id"`
	CallType                   any                                 `json:"callType"`
	CampaignID                 *CampaignID                         `json:"campaignId"`
	Customer                   any                                 `json:"customer"`
	MediaAvailability          string                              `json:"mediaAvailability"`
	ParkedCallsCount           uint64                              `json:"parkedCallsCount"`
	ReasonCodeID               string                              `json:"reasonCodeId"`
	State                      UserState                           `json:"state"`
	StateSince                 uint64                              `json:"stateSince"`
	StateDuration              uint64                              `json:"stateDuration"`
	OnHoldStateSince           uint64                              `json:"onHoldStateSince"`
	OnHoldStateDuration        uint64                              `json:"onHoldStateDuration"`
	OnParkStateSince           uint64                              `json:"onParkStateSince"`
	OnParkStateDuration        uint64                              `json:"onParkStateDuration"`
	ReasonCodeSince            uint64                              `json:"reasonCodeSince"`
	ReasonCodeDuration         uint64                              `json:"reasonCodeDuration"`
	AfterCallWorkStateSince    uint64                              `json:"afterCallWorkStateSince"`
	AfterCallWorkStateDuration uint64                              `json:"afterCallWorkStateDuration"`
	LoggedOutStateSince        uint64                              `json:"loggedOutStateSince"`
	LoggedOutStateDuration     uint64                              `json:"loggedOutStateDuration"`
	NotReadyStateSince         uint64                              `json:"notReadyStateSince"`
	NotReadyStateDuration      uint64                              `json:"notReadyStateDuration"`
	OnCallStateSince           uint64                              `json:"onCallStateSince"`
	OnCallStateDuration        uint64                              `json:"onCallStateDuration"`
	ReadyStateSince            uint64                              `json:"readyStateSince"`
	ReadyStateDuration         uint64                              `json:"readyStateDuration"`
	PermanentRecording         bool                                `json:"permanentRecording"`
	SessionRecording           bool                                `json:"sessionRecording"`
	Presence                   map[PresenceStateCode]PresenceState `json:"presence"`
	ChannelAvailability        map[Channel]ChannelState            `json:"channelAvailability"`
}

type WebSocketStatisticsUserSessionData struct {
	ID           SessionID `json:"id"`
	UserId       UserID    `json:"userId"`
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

type CampaignMode string

const (
	Basic    CampaignMode = "BASIC"
	Advanced CampaignMode = "ADVANCED"
)

type CampaignStateLabel string

const (
	Running    CampaignStateLabel = "RUNNING"
	NotRunning CampaignStateLabel = "NOT_RUNNING"
)

type Channel string

const (
	Video     Channel = "Video"
	Total     Channel = "Total"
	Chat      Channel = "Chat"
	Voicemail Channel = "Voicemail"
	Voice     Channel = "Voice"
)

type ChannelState struct {
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

type PresenceStateCode string

const (
	PendingState PresenceStateCode = "pendingState"
	CurrentState PresenceStateCode = "currentState"
)

type PresenceState struct {
	ReadyChannels      []string `json:"readyChannels"`
	NotReadyReasonCode uint64   `json:"notReadyReasonCode"`
}
