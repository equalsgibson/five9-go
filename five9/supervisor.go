package five9

import (
	"context"
	"net/http"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type SupervisorService struct {
	authState           *authenticationState
	webSocketHandler    webSocketHandler
	webSocketCache      *supervisorWebSocketCache
	domainMetadataCache *domainMetadataCache
}

func (s *SupervisorService) getDomainUserInfoMap(ctx context.Context) (map[five9types.UserID]five9types.AgentInfo, error) {
	if s.domainMetadataCache.agentInfoState.GetLastUpdated() != nil {
		if time.Since(*s.domainMetadataCache.agentInfoState.GetLastUpdated()) < time.Hour {
			return s.domainMetadataCache.agentInfoState.GetAll().Items, nil
		}
	}

	domainUserSlice, err := s.GetAllDomainUsers(ctx)
	if err != nil {
		return nil, err
	}

	freshData := map[five9types.UserID]five9types.AgentInfo{}

	for _, domainUser := range domainUserSlice {
		freshData[domainUser.ID] = domainUser
	}

	s.domainMetadataCache.agentInfoState.Replace(freshData)
	return freshData, nil
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

// func (s *SupervisorService) DomainReasonCodes(ctx context.Context) map[five9types.ReasonCodeID]five9types.ReasonCodeInfo {
// 	return s.domainMetadataCache.reasonCodes
// }
