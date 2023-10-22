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
	timeLastUpdated := s.domainMetadataCache.agentInfoState.GetCacheAge()

	if timeLastUpdated != nil {
		if *timeLastUpdated < time.Hour {
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

func (s *SupervisorService) GetStatisticsFilterSettings(ctx context.Context) ([]five9types.AgentInfo, error) {
	var target []five9types.AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/supervisors/:userID/stats_filter_settings",
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

func (s *SupervisorService) SetStatisticsFilterSettings(ctx context.Context, payload any) ([]five9types.AgentInfo, error) {
	var target []five9types.AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/supsvcs/rs/svc/supervisors/:userID/stats_filter_settings",
		structToReaderCloser(payload),
	)
	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
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

func (s *SupervisorService) getQueueInfoMap(ctx context.Context) (map[five9types.QueueID]five9types.SkillInfo, error) {
	timeLastUpdated := s.domainMetadataCache.queueInfoState.GetCacheAge()

	if timeLastUpdated != nil {
		if *timeLastUpdated < time.Hour {
			return s.domainMetadataCache.queueInfoState.GetAll().Items, nil
		}
	}

	queues, err := s.GetAllQueues(ctx)
	if err != nil {
		return nil, err
	}

	freshData := map[five9types.QueueID]five9types.SkillInfo{}

	for _, queue := range queues {
		freshData[queue.ID] = queue
	}

	s.domainMetadataCache.queueInfoState.Replace(freshData)
	return freshData, nil
}

func (s *SupervisorService) GetAllQueues(ctx context.Context) ([]five9types.SkillInfo, error) {
	var target []five9types.SkillInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/skills",
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

func (s *SupervisorService) getReasonCodeInfoMap(ctx context.Context) (map[five9types.QueueID]five9types.SkillInfo, error) {
	timeLastUpdated := s.domainMetadataCache.queueInfoState.GetCacheAge()

	if timeLastUpdated != nil {
		if *timeLastUpdated < time.Hour {
			return s.domainMetadataCache.queueInfoState.GetAll().Items, nil
		}
	}

	queues, err := s.GetAllQueues(ctx)
	if err != nil {
		return nil, err
	}

	freshData := map[five9types.QueueID]five9types.SkillInfo{}

	for _, queue := range queues {
		freshData[queue.ID] = queue
	}

	s.domainMetadataCache.queueInfoState.Replace(freshData)
	return freshData, nil
}

func (s *SupervisorService) GetAllReasonCodes(ctx context.Context) ([]five9types.SkillInfo, error) {
	var target []five9types.SkillInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/skills",
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
