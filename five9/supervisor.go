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
	websocketReady      chan bool
}

func (s *SupervisorService) GetAllDomainUsers(ctx context.Context) (map[five9types.UserID]five9types.AgentInfo, error) {
	// Check to see if we already have the data
	if s.domainMetadataCache.agentInfoState.state != nil {
		// Check to see when data was last fetched - older than 1 hour is considered stale.
		if time.Since(*s.domainMetadataCache.agentInfoState.state) < time.Hour {
			return s.domainMetadataCache.agentInfoState.agentInfo, nil
		}
	}

	s.domainMetadataCache.agentInfoState.mutex.Lock()
	defer s.domainMetadataCache.agentInfoState.mutex.Unlock()

	// Check to make sure another func hasn't managed to set the data
	if s.domainMetadataCache.agentInfoState.state != nil {
		// Check to see when data was last fetched - older than 1 hour is considered stale.
		if time.Since(*s.domainMetadataCache.agentInfoState.state) < time.Hour {
			return s.domainMetadataCache.agentInfoState.agentInfo, nil
		}
	}

	// Assume data is stale or nil at this point, so reach out to API to get fresh data
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

	//TODO: find elegant solution to removing stale users (deleted, suspended, no longer in result set)
	for _, agentInfo := range target {
		s.domainMetadataCache.agentInfoState.agentInfo[agentInfo.ID] = agentInfo
	}
	completedTime := time.Now()

	s.domainMetadataCache.agentInfoState.state = &completedTime
	return s.domainMetadataCache.agentInfoState.agentInfo, nil
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
