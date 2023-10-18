package five9

import (
	"context"
	"fmt"
	"net/http"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type AgentService struct {
	authState           *authenticationState
	websocketHandler    webSocketHandler
	webSocketCache      *agentWebsocketCache
	domainMetadataCache *domainMetadataCache
}

type agentWebsocketCache struct {
	// lastPong *time.Time
}

func (s AgentService) GetAllMaintenanceNoticesForSelf(ctx context.Context) ([]five9types.MaintenanceNoticeInfo, error) {
	var target []five9types.MaintenanceNoticeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/appsvcs/rs/svc/agents/:userID/maintenance_notices",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s AgentService) AcceptMaintenanceNoticeForSelf(
	ctx context.Context,
	maintenanceNoticeID five9types.MaintenanceNoticeID,
) (five9types.MaintenanceNoticeInfo, error) {
	var target five9types.MaintenanceNoticeInfo
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf(
			"/appsvcs/rs/svc/agents/:userID/maintenance_notices/%s/accept",
			maintenanceNoticeID,
		),
		http.NoBody,
	)
	if err != nil {
		return five9types.MaintenanceNoticeInfo{}, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return five9types.MaintenanceNoticeInfo{}, err
	}

	return target, nil
}

func (s *AgentService) GetAllReasonCodes(ctx context.Context) ([]five9types.ReasonCodeInfo, error) {
	logoutCodes, err := s.getAllLogoutReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	notReadyCodes, err := s.getAllNotReadyReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	return append(logoutCodes, notReadyCodes...), nil
}

func (s *AgentService) getAllLogoutReasonCodes(ctx context.Context) ([]five9types.ReasonCodeInfo, error) {
	var target []five9types.ReasonCodeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/appsvcs/rs/svc/orgs/:organizationID/logout_reason_codes",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *AgentService) getAllNotReadyReasonCodes(ctx context.Context) ([]five9types.ReasonCodeInfo, error) {
	var target []five9types.ReasonCodeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/appsvcs/rs/svc/orgs/:organizationID/not_ready_reason_codes",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}
