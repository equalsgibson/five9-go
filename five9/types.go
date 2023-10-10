package five9

import "fmt"

type (
	UserName              string
	CampaignID            string
	UserID                string
	UserState             string
	authenticationTokenID string
	farmID                string
	organizationID        string
	policy                string
	sessionID             string
	eventID               string
	userLoginState        string
	dataSource            string
)

const (
	eventIDServerConnected                    eventID = "1010"
	eventIDDuplicateConnection                eventID = "1020"
	eventIDPongReceived                       eventID = "1202" // Pong response to ping request
	eventIDSupervisorStats                    eventID = "5000" // Statistics data has been received.
	eventIDDispositionsInvalidated            eventID = "5002" // Disposition has been removed or	created, or disposition name has been changed.
	eventIDSkillsInvalidated                  eventID = "5003" // Skill has been removed or created, or	queue name has been changed.
	eventIDAgentGroupsInvalidated             eventID = "5004"
	eventIDCampaignsInvalidated               eventID = "5005"
	eventIDUsersInvalidated                   eventID = "5006"
	eventIDReasonCodesInvalidated             eventID = "5007"
	eventIDCampaignProfilesInvalidated        eventID = "5008"
	eventIDCampaignOutOfNumbers               eventID = "5009"
	eventIDListsInvalidated                   eventID = "5010"
	eventIDCampaignListsChanged               eventID = "5011"
	eventIDIncrementalStatsUpdate             eventID = "5012"
	eventIDIncrementalUserProfilesUpdate      eventID = "5013"
	eventIDFilterSettingsUpdated              eventID = "6001"
	eventIDAgentsInvalidated                  eventID = "6002"
	eventIDPermissionsUpdated                 eventID = "6003"
	eventIDResetCampaignDispositionsCompleted eventID = "6004"
	eventIDMonitoringStateUpdated             eventID = "6005"
	eventIDRandomMonitoringStarted            eventID = "6006"
	eventIDFdsRealTime                        eventID = "6007"
	eventIDIncrementalInteractions            eventID = "6008"
)

const (
	PolicyAttachExisting policy = "AttachExisting"
	PolicyForceIn        policy = "ForceIn"
)

const (
	dataSourceACDStatus                  dataSource = "ACD_STATUS"
	dataSourceAgentState                 dataSource = "AGENT_STATE"
	dataSourceAgentStatistic             dataSource = "AGENT_STATISTIC"
	dataSourceCampaignState              dataSource = "CAMPAIGN_STATE"
	dataSourceInboundCampaignStatistics  dataSource = "INBOUND_CAMPAIGN_STATISTICS"
	dataSourceStations                   dataSource = "STATIONS"
	dataSourceOutboundCampaignStatistics dataSource = "OUTBOUND_CAMPAIGN_STATISTICS"
	dataSourceOutboundCampaignManager    dataSource = "OUTBOUND_CAMPAIGN_MANAGER"
	dataSourceUserSession                dataSource = "USER_SESSION"
)

type webSocketIncrementalStatsUpdateData struct {
	DataSource dataSource   `json:"dataSource"`
	Added      []AgentState `json:"added"`
	Updated    []AgentState `json:"updated"`
	Removed    []AgentState `json:"removed"`
}

type websocketSupervisorStatsData struct {
	Data []AgentState `json:"data"`
}

type AgentState struct {
	ID                         UserID      `json:"id"`
	CallType                   any         `json:"callType"`
	CampaignID                 *CampaignID `json:"campaignId"`
	Customer                   any         `json:"customer"`
	MediaAvailability          string      `json:"mediaAvailability"`
	ParkedCallsCount           uint64      `json:"parkedCallsCount"`
	ReasonCodeID               string      `json:"reasonCodeId"`
	State                      UserState   `json:"state"`
	StateSince                 uint64      `json:"stateSince"`
	StateDuration              uint64      `json:"stateDuration"`
	OnHoldStateSince           uint64      `json:"onHoldStateSince"`
	OnHoldStateDuration        uint64      `json:"onHoldStateDuration"`
	OnParkStateSince           uint64      `json:"onParkStateSince"`
	OnParkStateDuration        uint64      `json:"onParkStateDuration"`
	ReasonCodeSince            uint64      `json:"reasonCodeSince"`
	ReasonCodeDuration         uint64      `json:"reasonCodeDuration"`
	AfterCallWorkStateSince    uint64      `json:"afterCallWorkStateSince"`
	AfterCallWorkStateDuration uint64      `json:"afterCallWorkStateDuration"`
	LoggedOutStateSince        uint64      `json:"loggedOutStateSince"`
	LoggedOutStateDuration     uint64      `json:"loggedOutStateDuration"`
	NotReadyStateSince         uint64      `json:"notReadyStateSince"`
	NotReadyStateDuration      uint64      `json:"notReadyStateDuration"`
	OnCallStateSince           uint64      `json:"onCallStateSince"`
	OnCallStateDuration        uint64      `json:"onCallStateDuration"`
	ReadyStateSince            uint64      `json:"readyStateSince"`
	ReadyStateDuration         uint64      `json:"readyStateDuration"`
	PermanentRecording         bool        `json:"permanentRecording"`
	SessionRecording           bool        `json:"sessionRecording"`
	ReadyChannels              string      `json:"readyChannels"`
	NotReadyReasonCode         uint64      `json:"notReadyReasonCode"`
	// ChannelAvailability        map[Channel]channelState `json:"channelAvailability"`
}

type StationInfo struct {
	StationID   string `json:"stationId"`
	StationType string `json:"stationType"`
}

type AgentInfo struct {
	ID       UserID   `json:"id"`
	UserName UserName `json:"userName"`
}

type loginResponse struct {
	TokenID   authenticationTokenID `json:"tokenID"`
	SessionID sessionID             `json:"sessionId"`
	OrgID     organizationID        `json:"orgID"`
	UserID    UserID                `json:"userID"`
	Context   sessionContext        `json:"context"`
	Metadata  sessionMetadata       `json:"metadata"`
}

func (v loginResponse) GetAPIHost() string {
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

type loginPayload struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
	AppKey              string              `json:"appKey"`
	Policy              policy              `json:"policy"`
}

type sessionContext struct {
	CloudClientURL string `json:"cloudClientUrl"`
	CloudTokenURL  string `json:"cloudTokenUrl"`
	FarmID         farmID `json:"farmId"`
}

type sessionMetadata struct {
	FreedomURL  string       `json:"freedomUrl"`
	DataCenters []dataCenter `json:"dataCenters"`
}

type dataCenter struct {
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
