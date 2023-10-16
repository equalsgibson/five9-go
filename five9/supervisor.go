package five9

import (
	"context"
	"net/http"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type SupervisorService struct {
	authState           *authenticationState
	websocketHandler    websocketHandler
	webSocketCache      *supervisorWebsocketCache
	domainMetadataCache *domainMetadata
}

func (s *SupervisorService) getDomainUserInfoMap(ctx context.Context) (map[five9types.UserID]five9types.AgentInfo, error) {
	// Check to see if we already have the data
	if s.domainMetadataCache.agentInfoState.lastUpdated != nil {
		// Check to see when data was last fetched - older than 1 hour is considered stale.
		if time.Since(*s.domainMetadataCache.agentInfoState.lastUpdated) < time.Hour {
			return s.domainMetadataCache.agentInfoState.agentInfo, nil
		}
	}

	s.domainMetadataCache.agentInfoState.mutex.Lock()
	defer s.domainMetadataCache.agentInfoState.mutex.Unlock()

	// Check to see if we already have the data
	if s.domainMetadataCache.agentInfoState.lastUpdated != nil {
		// Check to see when data was last fetched - older than 1 hour is considered stale.
		if time.Since(*s.domainMetadataCache.agentInfoState.lastUpdated) < time.Hour {
			return s.domainMetadataCache.agentInfoState.agentInfo, nil
		}
	}

	domainUserInfo, err := s.GetAllDomainUsers(ctx)
	if err != nil {
		return nil, err
	}

	s.domainMetadataCache.agentInfoState.agentInfo = map[five9types.UserID]five9types.AgentInfo{}
	for _, agentInfo := range domainUserInfo {
		s.domainMetadataCache.agentInfoState.agentInfo[agentInfo.ID] = agentInfo
	}

	completedTime := time.Now()

	s.domainMetadataCache.agentInfoState.lastUpdated = &completedTime

	return s.domainMetadataCache.agentInfoState.agentInfo, nil
}

func (s *SupervisorService) GetAllDomainUsers(ctx context.Context) ([]five9types.AgentInfo, error) {
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
