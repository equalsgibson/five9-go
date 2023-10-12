package five9

import (
	"context"
	"net/http"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type SupervisorService struct {
	authState           *authenticationState
	websocketHandler    websocketHandler
	webSocketCache      *supervisorWebsocketCache
	domainMetadataCache *domainMetadata
	websocketReady      chan bool
}

func (s *SupervisorService) getAllDomainUsers(ctx context.Context) ([]five9types.AgentInfo, error) {
	var target []five9types.AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/users",
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

func (s *SupervisorService) requestWebSocketFullStatistics(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/supsvcs/rs/svc/supervisors/:userID/request_full_statistics",
		http.NoBody,
	)
	if err != nil {
		return err
	}

	if err := s.authState.requestWithAuthentication(request, nil); err != nil {
		return err
	}

	return nil
}

func (s *SupervisorService) DomainReasonCodes(ctx context.Context) map[five9types.ReasonCodeID]five9types.ReasonCodeInfo {
	return s.domainMetadataCache.reasonCodes
}
